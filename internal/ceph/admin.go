package ceph

import (
	"context"
	"log"
	"os"
	"strings"

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

// EnsureUser creates the user if it doesn't exist, or adds a new key if it does.
// This handles the case where the Ceph user exists but the DB record was lost.
func EnsureUser(ctx context.Context, docker *dockerw.Client, uid, accessKey, secretKey string) error {
	err := CreateUser(ctx, docker, uid, accessKey, secretKey)
	if err == nil {
		return nil
	}
	if !strings.Contains(err.Error(), "user already exists") {
		return err
	}
	return docker.ExecInContainer(ctx, containerName, []string{
		"radosgw-admin", "key", "create",
		"--uid=" + uid,
		"--key-type=s3",
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
