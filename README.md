# Push'N'Pray

🙏 A platform as a service solution in Go.

## CLI

Push'N'Pray has a CLI that can installed by running this in your terminal.

```sh
curl -fsSL https://raw.githubusercontent.com/do-2k25-28/Push-N-Pray/refs/heads/main/scripts/install.sh | bash
```

## Manifest

The manifest describes your applications and their dependencies (Postgres, ...). It is named `pushnpray.toml` and is located at the root of your repository.
You can find an example [here](./pushnpray.toml.example).

## Applications

Applications are user provided programs that run. All applications listening on port 80 will be exposed to the world
using HTTPS.

> [!NOTE]
> All files saved to the file system are not kept when updating/redeploying your app. If your app needs data persistence look at [managed services](#managed-services).

### Dockerfile

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

### Static website

Use this service to deploy a static website. The server will run a nginx instance.

In the configuration you need to precise in `path` the root folder of your static web app. All files in this folder wil be statically served.

Example:

```toml
[apps]
[[apps.static]]
name = 'my-app'
path = 'src/webapp'
```

## Managed Services

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
