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

## S3 services

Add an S3 service to your `pushnpray.toml`, with the credentials of your Ceph user:

```toml
[services]
[[services.s3]]
name = "storage"
access-key = "my-access-key"
secret-key = "my-secret-key"
```

On first deploy, Push'N'Pray creates a bucket named `storage-<project-id>` in Ceph and injects the following environment variables into every app:

| Variable | Description |
|---|---|
| `S3_STORAGE_ACCESS_KEY` | Access key |
| `S3_STORAGE_SECRET_KEY` | Secret key |
| `S3_STORAGE_BUCKET` | Bucket name |

Use them in your app to talk to the bucket. Example with curl:

```sh
# Upload a file
curl -X PUT "$S3_STORAGE_ENDPOINT/$S3_STORAGE_BUCKET/hello.txt" \
  -H "Host: $(echo $S3_STORAGE_ENDPOINT | sed 's|http://||')" \
  --aws-sigv4 "aws:amz:us-east-1:s3" \
  --user "$S3_STORAGE_ACCESS_KEY:$S3_STORAGE_SECRET_KEY" \
  --data "hello world"

# Download it back
curl "$S3_STORAGE_ENDPOINT/$S3_STORAGE_BUCKET/hello.txt" \
  --aws-sigv4 "aws:amz:us-east-1:s3" \
  --user "$S3_STORAGE_ACCESS_KEY:$S3_STORAGE_SECRET_KEY"
```
