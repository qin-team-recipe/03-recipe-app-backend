FROM aopontann/dev_go_opencv:latest

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>
WORKDIR /workspaces

# [Optional] Uncomment the next lines to use go get to install anything else you need
USER vscode
# RUN go get -x <your-dependency-or-tool>

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.20.0 \
    && go install github.com/cosmtrek/air@v1.44.0 \
    && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.16.2

COPY go.mod go.sum ./
RUN go mod download
# USER root

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1
