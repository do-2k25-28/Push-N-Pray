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
}

type PostgresService struct {
	Service
	Version string
}

type RedisService struct {
	Service
	Version string
}

type S3Service struct {
	Service
}

var Manifest struct {
	ProjectId     string `toml:"project-id"`
	RepositoryUrl string `toml:"repository-url"`

	Apps struct {
		Dockerfile []DockerFileApp
		Docker     []DockerApp
	}

	Services struct {
		Postgres []PostgresService
		Redis    []RedisService
		S3       []S3Service
	}
}
