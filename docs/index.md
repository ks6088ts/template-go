# Setup logs

```bash
git clone https://github.com/ks6088ts/template-go.git
cd template-go

# Setup development environment on Docker
cp .env.sample .env
docker network create shared-network
docker-compose up -d
docker-compose exec go bash

# inside docker container

# cobra
make cobra-init
make cobra-add COBRA_CMD=hello
make build
./outputs/template-go --help
./outputs/template-go hello --help

# Cross Compile
# https://medium.com/@georgeok/cross-compile-in-go-golang-9f0d1261ee26
make build \
    GOBUILD='GOOS=windows GOARCH=amd64 go build' \
    BIN_PATH='outputs/windows/cli.exe'

# Docker image for release by multi-stage build
docker-compose -f docker-compose.release.yml build
docker-compose -f docker-compose.release.yml run --rm release ash -c "./cmd --help"

# push tag for release
git tag v0.0.0
git push origin v0.0.0

# follow updates for GitHub Actions specifications
# https://docs.github.com/en/free-pro-team@latest/actions/reference/workflow-commands-for-github-actions#setting-an-environment-variable
# https://docs.github.com/en/free-pro-team@latest/actions/reference/workflow-commands-for-github-actions#adding-a-system-path

# put version info into binary
# https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
# go build -ldflags="-X 'package_path.variable_name=new_value'"

# count system calls by strace
strace -e trace=write -c ./outputs/template-go

# repository architecture
# https://medium.com/veltra-engineering/how-to-release-golang-tools-with-go-get-1c856739f5e3
# github.com/<owner>/<repository>
# github.com/<owner>/<repository>/cmd/<repository>
```

# after creating repos from the template

```bash
# rename repository names via following command
make rename REPO_NAME=example-go
```
