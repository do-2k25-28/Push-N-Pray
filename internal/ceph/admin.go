package ceph

import (
	"os"
	"log"
	"context"
	"pushnpray/internal/dockerw"
)

const containerName = "ceph"

// CreateUser creates a Ceph RGW user via radosgw-admin inside the ceph container.
func CreateUser(ctx context.Context, docker *dockerw.Client, uid, accessKey, secretKey string) error {
	return docker.ExecInContainer(ctx, containerName, []string{
		"radosgw-admin", "user", "create",
		"--uid=" + uid,
		"--display-name=" + uid,
		"--access-key=" + accessKey,
		"--secret-key=" + secretKey,
	})
}

const defaultCephEndpoint = "http://10.200.0.2:8080"

func GetCephEndpoint() string {
	var ceph_endpoint = os.Getenv("CEPH_ENDPOINT")
	if ceph_endpoint == "" {
		log.Printf("No ceph endpoint defined ")
		ceph_endpoint = defaultCephEndpoint
	}
	log.Printf("Ceph endpoint: %s", ceph_endpoint)
	return ceph_endpoint
}

