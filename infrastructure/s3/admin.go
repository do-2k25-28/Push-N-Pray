package s3

import (
	"context"
	"pushnpray/internal/dockerw"
)

const cephContainer = "ceph"

// CreateCephUser creates a Ceph RGW user via radosgw-admin inside the ceph container.
func CreateCephUser(ctx context.Context, docker *dockerw.Client, uid, accessKey, secretKey string) error {
	return docker.ExecInContainer(ctx, cephContainer, []string{
		"radosgw-admin", "user", "create",
		"--uid=" + uid,
		"--display-name=" + uid,
		"--access-key=" + accessKey,
		"--secret-key=" + secretKey,
	})
}
