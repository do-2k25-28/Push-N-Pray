# Push'N'Pray

🙏 A platform as a service solution in Go.

## CLI

Push'N'Pray has a CLI that can installed by running this in your terminal.

```sh
curl -fsSL https://raw.githubusercontent.com/do-2k25-28/Push-N-Pray/refs/heads/main/scripts/install.sh | bash
```

## Managed services

Services are declared in `pushnpray.toml`:

```toml
[services]

[[services.postgres]]
id = "data"
version = "18"

[[services.redis]]
id = "cache"
version = "8"

[[services.s3]]
id = "s3"
```

Services are reconciled before application containers are deployed. Each service:

- is reachable from project applications through its `id` as hostname;
- uses a persistent Docker volume;
- cannot change type or version after creation;
- is stopped and removed with its volume when removed from the manifest.

Service IDs must contain only lowercase letters, numbers, and hyphens.
