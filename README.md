# gitflect

a lightweight go application that seamlessly mirrors your github repositories to gitlab. it clones repositories, includes large file storage objects, and securely pushes them to the destination without exposing tokens in the runtime logs.

## deployment

you can pull the official container image from the github registry or docker hub. configure your environment variables and start the service.

## variables

source_provider
the origin platform where the repositories live.

source_token
personal access token required for reading fetching repositories.

source_user
your username on the origin platform.

dest_provider
the destination platform where mirrored projects will be created.

dest_token
access token for the target platform registry creation.

dest_user
the target destination username.

dest_url
optional parameter if you are using a custom git instance on your own server.

repo_visibility
choose between all public or private to filter what should be replicated.

sync_interval
timer configuration for cron jobs scheduling. you can define time spans like 6h or 24h.

## simple usage

```bash
docker compose up
```
