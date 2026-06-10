package apps

import (
	"context"
	"os"
	"path"
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

func (app StaticWebApp) Prepare(ctx context.Context, docker *dockerw.Client, manifest manifest.Manifest) error {
	cwd := ctx.Value(WorkingDirectoryContextKey).(string)
	dockerfilePath := path.Join(cwd, app.imageName(manifest))

	userDockerfile := strings.Replace(dockerfile, userContentKey, app.Path, 1)

	if err := os.WriteFile(dockerfilePath, []byte(userDockerfile), 0644); err != nil {
		return err
	}

	if err := docker.BuildImage(ctx, app.imageName(manifest), dockerfilePath, cwd); err != nil {
		return err
	}

	return nil
}

func (app StaticWebApp) ContainerConfig(ctx context.Context, manifest manifest.Manifest) (dockerw.ContainerConfig, error) {
	healthcheck, err := healthcheckConfig(app.App)
	if err != nil {
		return dockerw.ContainerConfig{}, err
	}

	return dockerw.ContainerConfig{
		Image:       app.imageName(manifest),
		Healthcheck: healthcheck,
	}, nil
}

func NewStaticWebApp(manifest manifest.StaticWepApp) StaticWebApp {
	return StaticWebApp{manifest}
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
