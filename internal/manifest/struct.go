package manifest

type App struct {
	Name string
}

type DockerFileApp struct {
	App
	Dockerfile string
	Context    string
}

type DockerApp struct {
	App
	Image string
}

type Service struct {
	ID string `toml:"id"`
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
	ProjectId     string `toml:"project-id"`
	RepositoryUrl string `toml:"repository-url"`
	Server        string `toml:"server,omitempty"`

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

func (m *Manifest) GetServiceCount() int {
	return len(m.Services.S3) + len(m.Services.Redis) + len(m.Services.Postgres)
}
