# Swarmex Pack

Helm-like stack packaging CLI for Docker Swarm — template + values → deploy.

Part of [Swarmex](https://github.com/ccvass/swarmex) — enterprise-grade orchestration for Docker Swarm.

## What It Does

Packages Docker Swarm stacks as reusable templates with configurable values, similar to Helm charts. Supports Go template syntax in compose files, value overrides, and full lifecycle management (install, upgrade, uninstall).

## Labels

No service labels. This is a CLI tool. Commands:

- `install` — Render template and deploy the stack
- `upgrade` — Update an existing stack with new values
- `uninstall` — Remove the stack
- `render` — Preview the rendered compose file without deploying

## How It Works

1. Reads a pack file (templated compose YAML) and a values file.
2. Merges `--set` overrides on top of the values file.
3. Renders the final compose file using Go template engine.
4. Deploys (or previews) the rendered stack via `docker stack deploy`.

## Quick Start

```bash
# Preview rendered output
swarmex-pack render --pack-file pack.yaml --values values.yaml --set replicas=3

# Install a stack
swarmex-pack install my-stack --pack-file pack.yaml --values values.yaml

# Upgrade with overrides
swarmex-pack upgrade my-stack --pack-file pack.yaml --set image=app:v2

# Remove
swarmex-pack uninstall my-stack
```

## Verified

Render with value overrides produced correct output. Install deployed stack with 2/2 services running.

## License

Apache-2.0
