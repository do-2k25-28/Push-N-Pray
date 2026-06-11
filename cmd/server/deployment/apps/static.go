package apps

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"
	"strings"
)

type workingDirectoryCtxKeyType int

const WorkingDirectoryContextKey workingDirectoryCtxKeyType = 0

type StaticWebApp struct {
	manifest.StaticWepApp
}

func (app StaticWebApp) imageName(manifest manifest.Manifest) string {
	return "img-static-" + manifest.ProjectId + "-" + app.Name
}

func (app StaticWebApp) AppName() string {
	return app.Name
}

func (app StaticWebApp) LinkedApps() []string {
	return app.Links
}

func (app StaticWebApp) GetAllowOriginFrom() string {
	return app.AllowOriginFrom
}

func (app StaticWebApp) Prepare(ctx context.Context, docker *dockerw.Client, manifest manifest.Manifest) error {
	cwd := ctx.Value(WorkingDirectoryContextKey).(string)
	dockerfilePath := path.Join(cwd, app.imageName(manifest))

	// Validate app.Path before injecting it into the Dockerfile COPY instruction.
	if err := validateStaticPath(cwd, app.Path); err != nil {
		return err
	}

	userDockerfile := strings.Replace(dockerfile, userContentKey, app.Path, 1)

	if err := os.WriteFile(dockerfilePath, []byte(userDockerfile), 0644); err != nil {
		return err
	}

	if err := docker.BuildImage(ctx, app.imageName(manifest), dockerfilePath, cwd); err != nil {
		return err
	}

	return nil
}

func (app StaticWebApp) ContainerConfig(ctx context.Context, manifest manifest.Manifest) dockerw.ContainerConfig {
	return dockerw.ContainerConfig{
		Image: app.imageName(manifest),
	}
}

func NewStaticWebApp(manifest manifest.StaticWepApp) StaticWebApp {
	return StaticWebApp{manifest}
}

// validateStaticPath ensures that the user-supplied static path is a safe
// relative path within the workspace. It rejects absolute paths, traversals,
// and symlinks that point outside the workspace.
func validateStaticPath(workspace, p string) error {
	if filepath.IsAbs(p) {
		return fmt.Errorf("static path %q must be relative to the workspace", p)
	}
	if strings.HasPrefix(filepath.Clean(p)+string(filepath.Separator), ".."+string(filepath.Separator)) ||
		filepath.Clean(p) == ".." {
		return fmt.Errorf("static path %q escapes workspace root", p)
	}

	// Resolve symlinks and verify containment within the workspace.
	candidate := filepath.Join(workspace, p)
	resolved, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		return fmt.Errorf("cannot resolve static path %q: %w", p, err)
	}
	cleanBase := filepath.Clean(workspace)
	if !strings.HasPrefix(resolved+string(filepath.Separator), cleanBase+string(filepath.Separator)) {
		return fmt.Errorf("static path %q escapes workspace root", p)
	}

	return nil
}

const userContentKey = "$USER_CONTENT"

const dockerfile = `
FROM nginx:alpine

# Remove the default Nginx static assets
RUN rm -rf /usr/share/nginx/html/*

# Copy the user's static website folder into the Nginx HTML directory
COPY ` + userContentKey + ` /usr/share/nginx/html

# Expose port 80 to the outside world
EXPOSE 80

# Start Nginx and keep it running in the foreground
CMD ["nginx", "-g", "daemon off;"]
`
