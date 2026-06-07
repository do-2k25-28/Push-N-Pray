package ceph

import (
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
