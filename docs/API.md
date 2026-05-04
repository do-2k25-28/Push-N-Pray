# API

## Access

This document covers the internal HTTP REST API exposed by the backend.

It is accessible at `http(s)://api.<DOMAIN>/v1/`.

If using our Push'N'Pray instance it is `https://api.pushnpray.polydo.dev/v1/`.

## Authentication

Using a `Authorization` header that is a bearer token.

## App

### Health check

This route returns `204` if the app is running.

```http
GET /v1/health
```

```http
HTTP/1.1 204 No Content
```

## Authentication

### Register

Register an account on the platform.

```http
POST /v1/auth/register

Content-Type: application/json

{
  "email": "john.doe@acme.org",
  "username": "jdoe",
  "password": "superSecretPassword"
}
```

```http
HTTP/1.1 200 OK

Content-Type: application/json

{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.xTauSR2dlM1bJuIiwlRHy0Sj-66g5_7qL2RKWT2u5J4",
  "refreshToken": "8l86wueqyffatorwoy3kh9cy4v2665sw"
}
```

### Login

Login using email/password combo.

```http
POST /v1/auth/login

Content-Type: application/json

{
  "email": "john.doe@acme.org",
  "password": "superSecretPassword"
}
```

```http
HTTP/1.1 200 OK

Content-Type: application/json

{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.xTauSR2dlM1bJuIiwlRHy0Sj-66g5_7qL2RKWT2u5J4",
  "refreshToken": "jn2bqk5gm3153gzljs0f4krvagf5lhv3"
}
```

### Token

Use the refresh token to obtain a new access and refresh token.

```http
POST /v1/auth/token

Content-Type: application/json

{
  "refreshToken": "jn2bqk5gm3153gzljs0f4krvagf5lhv3"
}
```

```http
HTTP/1.1 200 OK

Content-Type: application/json

{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.xTauSR2dlM1bJuIiwlRHy0Sj-66g5_7qL2RKWT2u5J4",
  "refreshToken": "xlu1qrzcxjn7wlho8k10pvptccsjea35"
}
```

## PAT

### List PATs

Lists the personal access tokens linked to the account.

```http
GET /v1/tokens
```

```http
HTTP/1.1 200 OK

Content-Type: application/json

{
  "tokens": [
    {
      "id": "00b6fda9-52be-4ce9-8635-73784bb4a4fc",
      "name": "GitHub Actions",
      "expiresAt": "1785574098"
    }
  ]
}
```

### Create PAT

```http
POST /v1/tokens

Authorization: Bearer <accessToken>
Content-Type: application/json

{
  "name": "CLI for my laptop",
  "expiresAt": null
}
```

```http
HTTP/1.1 200 OK

Content-Type: application/json

{
  "id": "8eda2225-a0dd-41ef-a258-c9b3a6d34ce2",
  "token": "8bb8069ce18241d89dfaf1fbc635d642"
}
```

### Delete PAT

```http
DELETE /v1/tokens/:tokenId

Authorization: Bearer <accessToken>
```

```
HTTP/1.1 204 No Content
```

## Projects

### Create project

Create a project with the given name, slug and repository.

```http
POST /v1/projects

Authorization: Bearer <accessToken>
Content-Type: application/json

{
  "name": "My revolutionary project",
  "slug": "my-revolutionary-project",
  "repository": {
    "url": "https://github.com/RichardDorian/CitiesAPI.git",
    "branch": "main"
  }
}
```

```http
HTTP/1.1 200 OK

Content-Type: application/json

{
  "id": "p0ZoB1FwH6"
}
```

### Delete project

Delete the project and shut down any deployment.

```http
DELETE /v1/projects/:projectId

Authorization: Bearer <accessToken>
```

```http
HTTP/1.1 204 No Content
```

### Deploy project

Trigger a deployment, optionally on given tag/commit.
The API responds with a deployment identifier.

```http
POST /v1/projects/:projectId/deploy

Authorization: Bearer <accessToken>
Content-Type: application/json

{
  "tag": "v0.0.1"
}
```

or

```http
POST /v1/projects/:projectId/deploy

Authorization: Bearer <accessToken>
Content-Type: application/json

{
  "commit": "eee21577cf5d079e96142cae34df17d886e8bc97"
}
```

or

```http
POST /v1/projects/:projectId/deploy

Authorization: Bearer <accessToken>
```

```http
HTTP/1.1 200 OK

Content-Type: application/json

{
  "id": "ca9902d4-fd31-4a0b-a529-466c26760ab1"
}
```

### Deployment info

Get deployment status

```http
GET /v1/projects/:projectId/deployments/:deploymentId

Authorization: Bearer <accessToken>
```

```
HTTP/1.1 200 OK

Content-Type: application/json

{
  "status": "in-progress"
}
```

or

```
HTTP/1.1 200 OK

Content-Type: application/json

{
  "status": "success",
  "url": "https://my-revolutionary-project-p0ZoB1FwH6.pushnpray.polydo.dev"
}
```

```
HTTP/1.1 200 OK

Content-Type: application/json

{
  "status": "error",
  "message": "Service 'backend' is unhealthy."
}
```

| Status        | Description                          |
| ------------- | ------------------------------------ |
| `in-progress` | Deployment is in progress.           |
| `success`     | Project deployed successfuly.        |
| `error`       | An error happened during deployment. |

### Environment variables

Set environment variables that will be injected for all services running in the project.

```http
POST /v1/projects/:projectId/env

Authorization: Bearer <accessToken>
Content-Type: application/json

{
  "variables": [
    {
      "name": "MY_CUSTOM_ENV_VARIABLE",
      "value": "c2997d1b7f93405c957417141be22c73"
    }
  ]
}
```

```http
HTTP/1.1 204 OK
```