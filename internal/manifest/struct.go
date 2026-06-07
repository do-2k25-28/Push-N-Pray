package manifest

type App struct {
	Name string `toml:"name"`
}

type DockerFileApp struct {
	App
	Dockerfile string `toml:"dockerfile"`
	Context    string `toml:"context"`
}

type DockerApp struct {
	App
	Image string `toml:"image"`
}

type Service struct {
	Name   string   `toml:"name"`
	UsedBy []string `toml:"used-by"`
}

type PostgresService struct {
	Service
	Version string `toml:"version"`
}

type RedisService struct {
	Service
	Version string `toml:"version"`
}

type S3Service struct {
	Service
}

type Manifest struct {
	ProjectId string `toml:"project-id"`
	Server    string `toml:"server,omitempty"`

	Apps struct {
		Dockerfile []DockerFileApp `toml:"dockerfile,omitempty"`
		Docker     []DockerApp     `toml:"docker,omitempty"`
	} `toml:"apps,omitempty"`

	Services struct {
		Postgres []PostgresService `toml:"postgres,omitempty"`
		Redis    []RedisService    `toml:"redis,omitempty"`
		S3       []S3Service       `toml:"s3,omitempty"`
	} `toml:"services,omitempty"`
}

func (m *Manifest) GetApplicationCount() int {
	return len(m.Apps.Docker) + len(m.Apps.Dockerfile)
}

func (m *Manifest) GetApps() []App {
	apps := make([]App, 0, m.GetApplicationCount())
	for _, app := range m.Apps.Docker {
		apps = append(apps, app.App)
	}
	for _, app := range m.Apps.Dockerfile {
		apps = append(apps, app.App)
	}
	return apps
}

func (m *Manifest) GetServiceCount() int {
	return len(m.Services.S3) + len(m.Services.Redis) + len(m.Services.Postgres)
}
