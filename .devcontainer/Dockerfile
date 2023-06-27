FROM mcr.microsoft.com/devcontainers/go:0-1.20 as builder

RUN apt-get update && apt-get install -y \
    unzip \
    cmake \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN curl -SL https://github.com/opencv/opencv/archive/4.7.0.zip -o opencv-4.7.0.zip \
    && curl -SL https://github.com/opencv/opencv_contrib/archive/refs/tags/4.7.0.zip -o opencv_contrib-4.7.0.zip \
    && unzip opencv-4.7.0.zip \
    && unzip opencv_contrib-4.7.0.zip -d ./opencv-4.7.0/
RUN cd opencv-4.7.0 \
    && mkdir build \
    && cd build \
    && cmake -DOPENCV_GENERATE_PKGCONFIG=ON -DOPENCV_EXTRA_MODULES_PATH=../opencv_contrib-4.7.0/modules .. \
    && cmake --build . \
    && make \
    && make install \
    && ldconfig

FROM mcr.microsoft.com/devcontainers/go:0-1.20

RUN apt-get update && apt-get install -y \
    lsb-release \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/local/ /usr/local/
RUN ldconfig

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>
RUN curl -sSf https://atlasgo.sh | sh

# [Optional] Uncomment the next lines to use go get to install anything else you need
USER vscode
# RUN go get -x <your-dependency-or-tool>
RUN go install github.com/kyleconroy/sqlc/cmd/sqlc@latest \
    && go install github.com/cosmtrek/air@latest
# USER root

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1
