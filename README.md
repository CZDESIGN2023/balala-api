## Install Kratos
```
go env -w GOPROXY=https://goproxy.io,direct
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=off
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest

#ubuntu
sudo apt install protobuf-compiler

#mac
brew install go

# Error: Cannot install in Homebrew on ARM processor in Intel default prefix (/usr/local)!
(cd /opt
sudo mkdir homebrew
sudo git clone https://github.com/Homebrew/brew homebrew
eval "$(homebrew/bin/brew shellenv)"
sudo chown -R $(whoami) /opt/homebrew
brew update --force --quiet
chmod -R go-w "$(brew --prefix)/share/zsh"
)

brew install protobuf
go get -u github.com/golang/protobuf/proto
go get -u github.com/golang/protobuf/protoc-gen-go

go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/google/wire/cmd/wire@latest
go install github.com/golang/mock/gomock@latest
go install github.com/golang/mock/mockgen@latest

export GOROOT="$(brew --prefix golang)/libexec"
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

## docker安装一下服务
```
docker run -d -it --name nats-server -p 4222:4222 -p 6222:6222 -p 8222:8222 nats
docker run -d -it --name redis-svr -p 16379:6379 redis
docker run -d --name etcd \
    -p 2379:2379 -p 2380:2380 \
    -e ETCD_ADVERTISE_CLIENT_URLS=http://0.0.0.0:2379 \
    -e ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379 \
    quay.io/coreos/etcd:v3.5.1

```

## Install dart
```
brew tap dart-lang/dart
brew install dart
dart pub global activate protoc_plugin 20.0.1

vi ~/.zprofile
export PATH=$PATH:~/.pub-cache/bin
```

## Create a service
```
# Create a template project
kratos new server

cd server
# Add a proto template
kratos proto add api/server/server.proto
# Generate the proto code
kratos proto client api/server/server.proto
# Generate the source code of service by proto file
kratos proto server api/server/server.proto -t internal/service

go generate ./...
go build -o ./bin/ ./...
./bin/server -conf ./configs
```
## Generate other auxiliary files by Makefile
```
# Download and update dependencies
make init
# Generate API files (include: pb.go, http, grpc, validate, swagger) by proto file
make api
# Generate all files
make all
```
## Automated Initialization (wire)
```
# install wire
go get github.com/google/wire/cmd/wire

# generate wire
cd cmd/server
wire
```

## Docker
```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```

## 安裝Ginkgo測試框架命令行
```
go install github.com/onsi/ginkgo/ginkgo
```
