# Deploy

This guide covers the deployment of a Push'N'Pray instance.

## Guide

Start by creating a folder somewhere on your server.

```bash
mkdir -p ~/pushnpray
cd ~/pushnpray
```

Download the `compose.production.yml` from the repository.

```bash
curl -fsSL https://raw.githubusercontent.com/do-2k25-28/Push-N-Pray/refs/heads/main/compose.yml > compose.yml
```

Save your Traefik config as `traefik.yml`. See [TLS providers](#tls-providers) more details.

You need to set up these environment variables:

```bash
export API_HOST=<API SERVER HOST>

export PG_USER=pushnpray
export PG_PASSWORD=<STRONG PASSWORD>
export PG_DB=pushnpray

export TRAEFIK_CONFIG_PATH="~/pushnpray/traefik.yml"
export TRAEFIK_NET=<AVAILABLE DOCKER NETWORK NAME>
```

If your Traefik config requires environment variables put them in `traefik.env` (for exemple Cloudflare API tokens).

## TLS providers

### Let's Encrypt + DNS Challenge + Cloudflare provider

We use Traefik for routing HTTPs traffic. The TLS configuration depends on your infrastructure and your DNS provider.

A premade config is available [here](https://raw.githubusercontent.com/do-2k25-28/Push-N-Pray/refs/heads/main/infrastructure/traefik/config.yml) that:
  - Uses Let's Encrypt as the certificate provider
  - Uses the DNS challenge
  - Uses the Cloudflare provider for the DNS challenge
  - For our domain `*.pushnpray.polydo.dev`

Get it using

```bash
curl -fsSL https://raw.githubusercontent.com/do-2k25-28/Push-N-Pray/refs/heads/main/infrastructure/traefik/config.yml > traefik.yml
```
