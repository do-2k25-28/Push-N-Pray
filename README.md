# Push'N'Pray

A platform-as-a-service in Go. Declare your apps and services in a `pushnpray.toml` manifest, push to your repository, and Push'N'Pray handles the rest.

## CLI

Install the CLI with:

```sh
curl -fsSL https://raw.githubusercontent.com/do-2k25-28/Push-N-Pray/refs/heads/main/scripts/install.sh | bash
```

## Manifest

Every project needs a `pushnpray.toml` at the root of the repository.

```toml
project-id = "my-project"

[apps]

# Deploy a pre-built image
[[apps.docker]]
name = "api"
image = "docker.io/myorg/myapp:latest"

# Or build from a Dockerfile
[[apps.dockerfile]]
name = "api"
dockerfile = "Dockerfile"
context = "."

[services]

[[services.postgres]]
name = "db"
used-by = ["api"]

[[services.s3]]
name = "storage"
used-by = ["api"]
```

`used-by` controls which app containers receive the service's environment variables.

## Managed services

Services are provisioned before app containers are deployed. On first deploy each service is created; on subsequent deploys the existing service is reused.

### Postgres

Credentials are auto-generated and injected into apps listed in `used-by`:

| Variable | Value |
|---|---|
| `POSTGRES_<NAME>_HOST` | Container hostname |
| `POSTGRES_<NAME>_PORT` | `5432` |
| `POSTGRES_<NAME>_USER` | `postgres` |
| `POSTGRES_<NAME>_PASSWORD` | Auto-generated password |

### S3

Push'N'Pray creates a dedicated Ceph RGW user with auto-generated credentials and a bucket named `<name>-<project-id>`. No credentials go in the manifest.

```toml
[[services.s3]]
name = "storage"
used-by = ["api"]
```

The following variables are injected into apps listed in `used-by`:

| Variable | Value |
|---|---|
| `S3_<NAME>_ENDPOINT` | `http://ceph:8080` |
| `S3_<NAME>_ACCESS_KEY` | Auto-generated access key |
| `S3_<NAME>_SECRET_KEY` | Auto-generated secret key |
| `S3_<NAME>_BUCKET` | `<name>-<project-id>` |

Example using the injected variables from inside a container:

```sh
# Upload
curl -X PUT "$S3_STORAGE_ENDPOINT/$S3_STORAGE_BUCKET/hello.txt" \
  --aws-sigv4 "aws:amz:us-east-1:s3" \
  --user "$S3_STORAGE_ACCESS_KEY:$S3_STORAGE_SECRET_KEY" \
  --data "hello world"

# Download
curl "$S3_STORAGE_ENDPOINT/$S3_STORAGE_BUCKET/hello.txt" \
  --aws-sigv4 "aws:amz:us-east-1:s3" \
  --user "$S3_STORAGE_ACCESS_KEY:$S3_STORAGE_SECRET_KEY"
```

## Server setup

### `.env`

Docker Compose reads `.env` from the project root. Create it before running `docker compose up`:

```sh
# .env
POSTGRES_PASSWORD=<random>
CEPH_DEMO_ACCESS_KEY=<random>
CEPH_DEMO_SECRET_KEY=<random>
```

Generate it in one command:

```sh
printf 'POSTGRES_PASSWORD=%s\nCEPH_DEMO_ACCESS_KEY=%s\nCEPH_DEMO_SECRET_KEY=%s\n' \
  $(openssl rand -hex 16) \
  $(openssl rand -hex 16) \
  $(openssl rand -hex 32) > .env
```

Then start the stack and the server:

```sh
docker compose up -d
make run
```

### Server environment variables

| Variable | Default | Description |
|---|---|---|
| `CEPH_ENDPOINT` | `http://10.200.0.2:8080` | Ceph RGW endpoint reachable from the server host |
| `DB_PASSWORD` | `postgres` | Postgres password (must match `POSTGRES_PASSWORD` in `.env`) |
| `DB_HOST` | `localhost` | Postgres host |
| `HTTP_PORT` | `4000` | Server port |
