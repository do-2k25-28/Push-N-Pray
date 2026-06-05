package utils

import (
	"encoding/binary"
	"fmt"
	"net"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func ResolvePath(base, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(base, path)
}

func FindAvailablePort(defaultPort string) string {
	port, err := strconv.Atoi(defaultPort)
	if err != nil {
		port = 4000
	}

	for port <= 65535 {
		addr := fmt.Sprintf(":%d", port)
		l, err := net.Listen("tcp", addr)
		if err == nil {
			_ = l.Close()
			return strconv.Itoa(port)
		}
		port++
	}

	l, err := net.Listen("tcp", ":0")
	if err == nil {
		defer func() {
			_ = l.Close()
		}()
		return strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
	}

	return defaultPort
}

func getUsedDockerSubnets() (map[string]bool, error) {
	usedSubnets := map[string]bool{}

	lsOut, err := exec.Command("docker", "network", "ls", "-q").Output()
	if err != nil {
		return nil, fmt.Errorf("error listing networks: %w", err)
	}

	networkIDs := strings.Fields(string(lsOut))
	if len(networkIDs) == 0 {
		return usedSubnets, nil
	}

	args := append([]string{"network", "inspect", "--format", "{{range .IPAM.Config}}{{.Subnet}}{{end}}"}, networkIDs...)
	inspectOut, err := exec.Command("docker", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("error inspecting networks: %w", err)
	}

	for _, subnet := range strings.Fields(string(inspectOut)) {
		usedSubnets[subnet] = true
	}

	return usedSubnets, nil
}

func findFreeSubnet(usedSubnets map[string]bool) (string, error) {
	for secondByte := 18; secondByte <= 255; secondByte++ {
		subnet := fmt.Sprintf("172.%d.0.0/16", secondByte)
		if !usedSubnets[subnet] {
			return subnet, nil
		}
	}
	return "", fmt.Errorf("no available subnet found in 172.18-255.0.0/16 range")
}

func createDockerNetwork(networkName, subnet string) error {
	_, err := exec.Command("docker", "network", "create", "--subnet", subnet, networkName).Output()
	if err != nil {
		return fmt.Errorf("error creating network %s with subnet %s: %w", networkName, subnet, err)
	}
	return nil
}

type CephNetwork struct {
	Name      string
	Subnet    string
	MonitorIP string
}

func deriveMonitorIP(cidr string) (string, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", fmt.Errorf("invalid subnet %q: %w", cidr, err)
	}

	ip := ipNet.IP.To4()
	if ip == nil {
		return "", fmt.Errorf("only IPv4 subnets are supported, got %q", cidr)
	}

	n := binary.BigEndian.Uint32(ip)
	n += 2
	result := make(net.IP, 4)
	binary.BigEndian.PutUint32(result, n)

	if !ipNet.Contains(result) {
		return "", fmt.Errorf("derived monitor IP %s is outside subnet %s", result, cidr)
	}

	return result.String(), nil
}

func getExistingNetworkSubnet(networkName string) (string, error) {
	out, err := exec.Command("docker", "network", "inspect", "--format",
		"{{range .IPAM.Config}}{{.Subnet}}{{end}}", networkName).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func InitCephNetwork(networkName string) (*CephNetwork, error) {
	if subnet, err := getExistingNetworkSubnet(networkName); err == nil && subnet != "" {
		monitorIP, err := deriveMonitorIP(subnet)
		if err != nil {
			return nil, err
		}
		return &CephNetwork{Name: networkName, Subnet: subnet, MonitorIP: monitorIP}, nil
	}

	usedSubnets, err := getUsedDockerSubnets()
	if err != nil {
		return nil, err
	}

	subnet, err := findFreeSubnet(usedSubnets)
	if err != nil {
		return nil, err
	}

	if err := createDockerNetwork(networkName, subnet); err != nil {
		return nil, err
	}

	monitorIP, err := deriveMonitorIP(subnet)
	if err != nil {
		return nil, err
	}

	return &CephNetwork{Name: networkName, Subnet: subnet, MonitorIP: monitorIP}, nil
}
