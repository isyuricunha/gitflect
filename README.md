# gitflect

a lightweight go container that automatically and safely mirrors your github repositories to gitlab. it clones repositories, pulls large file storage objects, and pushes them seamlessly while keeping your environment highly optimized.

## how it works

the application is completely stateless. it does not save partial progress or cache clones on your disk avoiding storage bloat over time. upon each execution span it fetches your remote repository fresh into a temporary folder, executes a mirror swap, pushes everything to the destination provider, and immediately deletes the temporary files.

regarding rate limits from github or gitlab APIs, the networking logic inside uses exponential read requests and clean pagination. if a huge amount of repositories exist, it handles the pagination smoothly.

## available registries

you can pull the application from two mirrored registries

* docker hub 
* github container registry

for docker hub use `isyuricunha/gitflect:latest`
for github registry use `ghcr.io/isyuricunha/gitflect:latest`

you can also lock your version using a specific tagged version like `gitflect:v1.0.0` instead of latest.

## configuration parameters

source_provider
the platform you are copying from (like github).

source_token
personal access token. for github we highly recommend using a fine-grained personal access token.
required scopes for github fine-grained token:
* `contents`: read-only
* `metadata`: read-only

source_user
your username on the origin platform.

dest_provider
the target platform (like gitlab).

dest_token
access token for the target platform registry creation. for gitlab a standard personal access token is required.
required scopes for gitlab token:
* `api`: complete read/write access to the api (required to create missing repositories)
* `read_repository`: read access to repositories
* `write_repository`: write access to repositories via git-over-http

dest_user
the target destination username.

dest_url
optional. provide an absolute domain if you are using a self hosted custom git instance instead of the public cloud.

repo_visibility
can be public or private. use all to mirror everything.

repo_include
optional comma separated list of exact repository names. if used only these will be synced.

repo_exclude
optional comma separated list. skips these exact repositories.

sync_interval
time frequency for synchronization scheduling. you can define time spans like 6h or 24h. leave empty to run once and exit.

## deployment with separate environment file

if you prefer keeping your secrets in an isolated file, create a file named `.env` and configure your settings

```bash
source_provider=github
source_token=ghp_yourtoken
source_user=youruser
dest_provider=gitlab
dest_token=glpat_yourtoken
dest_user=youruser
sync_interval=6h
```

then configure your compose file like this

```yaml
services:
  gitflect:
    image: ghcr.io/isyuricunha/gitflect:latest
    restart: always
    env_file: .env
```

## deployment with inline environment variables

if you prefer having everything in a single file without external references, you can pass the variables inside the environment block directly

```yaml
services:
  gitflect:
    image: ghcr.io/isyuricunha/gitflect:latest
    restart: always
    environment:
      source_provider: github
      source_token: ghp_yourtoken
      source_user: youruser
      dest_provider: gitlab
      dest_token: glpat_yourtoken
      dest_user: youruser
      repo_visibility: all
      sync_interval: 12h
```
