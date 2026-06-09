# Push'N'Pray

A platform-as-a-service in Go. Declare your apps and services in a `pushnpray.toml` manifest, push to your repository, and Push'N'Pray handles the rest.

## CLI

Install the CLI with:

```sh
curl -fsSL https://raw.githubusercontent.com/do-2k25-28/Push-N-Pray/refs/heads/main/scripts/install.sh | bash
```

## Manifest

Services are declared in `pushnpray.toml`. You can find an example [here](./pushnpray.toml.example).

### Applications

Applications are user provided programs that run. All applications listening on port 80 will be exposed to the world
using HTTPS.

> [!NOTE]
> All files saved to the file system are not kept when updating/redeploying your app. If your app needs data persistence look at [managed services](#managed-services).

#### Dockerfile

The server will be build the given dockerfile in the given context and then deploy it.

Example:

```toml
[apps]
[[apps.dockerfile]]
name = 'my-app'
dockerfile = 'src/app1/Dockerfile'
context = 'src/app1/'
```

### Docker

The server will pull the given image and deploy it.

> [!IMPORTANT]
> The image needs to be publicly accessible.

Example:

```toml
[apps]
[[apps.docker]]
name = 'my-app'
image = 'ghcr.io/jdoe/my-app:latest'
```

### Managed Services

Services are applications that your app may depend on such as Postgres or Redis. They are managed by us (deployment, data persistence).

All managed services require a kind of authentication and some details. For example, Postgres require a database name and a default user/password. These are generated for you and can be injected to your app using environment variables. To avoid leaking secrets, you have to explicitly tell which apps can access these secrets. This is done using the `used-by` which take the list of application names.

In the provided manifest example, we can see that the service `db` (postgres service) is used by the app `backend`. Therefore only the `backend` app will have `db` secrets injected.

> [!NOTE]
> Since you can deploy multiple instances of the same service, environment variables are following this format `{SERVICE_TYPE}_{SERVICE_NAME}`. For example `POSTGRES_MY_DB`.

### Postgres

The `postgres` service is Postgres version 18.

Example:

```toml
[services]
[[services.postgres]]
name = 'db'
used-by = ['backend']
```

Injected variables are:

| Name                       | Description                           |
| -------------------------- | ------------------------------------- |
| `POSTGRES_{NAME}_USER`     | Postgres user to use.                 |
| `POSTGRES_{NAME}_PASSWORD` | Postgres password to use.             |
| `POSTGRES_{NAME}_HOST`     | Where to reach the Postgres instance. |
| `POSTGRES_{NAME}_PORT`     | Port Postgres is listening on.        |

### S3

Push'N'Pray creates a dedicated Ceph RGW user with auto-generated credentials and a bucket named `<name>-<project-id>`. No credentials go in the manifest.

```toml
[[services.s3]]
name = "storage"
used-by = ["api"]
```

The following variables are injected into apps listed in `used-by`:

| Variable                 | Description            |
|--------------------------|------------------------|
| `S3_{NAME}_ENDPOINT`     | Ceph container GW      |
| `S3_{NAME}_ACCESS_KEY`   | Ceph bucket access key |
| `S3_{NAME}_SECRET_KEY`   | Ceph bucket secret key |
| `S3_{NAME}_BUCKET`       | Ceph bucket name       |

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

