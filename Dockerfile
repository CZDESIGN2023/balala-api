ARG GO_VERSION=1.24.1

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base

# 配置go代理
#ENV GOPROXY=https://goproxy.cn

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go mod download -x

FROM --platform=$BUILDPLATFORM base AS build

ARG TARGETOS
ARG TARGETARCH

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /src/bin/go-cs ./cmd/go-cs

FROM scratch AS binary
COPY --from=build /src/bin/go-cs .

FROM alpine:latest AS prod

# 安装时区数据
RUN apk --no-cache add tzdata ffmpeg

WORKDIR /app

# 从build阶段拷贝go程序
COPY --from=build /src/bin/go-cs ./bin/
# 从build阶段拷贝配置
COPY --from=build /src/configs/config.docker_compose.yaml.dist ./configs/config.yaml
COPY --from=build /src/configs/qqwry.dat ./configs/
COPY --from=build /src/resource ./resource

ENTRYPOINT ["./bin/go-cs"]
