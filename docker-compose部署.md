# **工作目录：项目根目录**

## 多平台构建
### 创建构建器
```bash
docker buildx create --driver=docker-container --name=container
docker buildx use container
```
### 获取当前分支名
```bash
TAG_VER=$(git rev-parse --abbrev-ref HEAD)
TAG_VER=${TAG_VER#docker/}
echo $TAG_VER
```

### 构建 api 镜像
```bash
docker buildx build \
--platform=linux/arm64,linux/amd64 \
--tag=czdesign2023/balala-api:$TAG_VER \
--push \
.
```

### 构建数据库镜像
```bash
docker buildx build \
--platform=linux/arm64,linux/amd64 \
--tag=czdesign2023/balala-mysql:$TAG_VER \
--push \
./sql
```

### 构建mysql同步es镜像
```bash
docker buildx build \
--platform=linux/arm64,linux/amd64 \
--tag=czdesign2023/balala-mysql2es:$TAG_VER \
--push \
./mysql-es
```

### 构建前端镜像
```bash
docker buildx build \
--platform=linux/arm64,linux/amd64 \
--tag=czdesign2023/balala-web:$TAG_VER \
--push \
../balala_client_vue
```

## 当前平台构建
### 构建 api 镜像
```bash
docker build \
--tag=czdesign2023/balala-api:$TAG_VER \
.
```

### 构建数据库镜像
```bash
docker build \
--tag=czdesign2023/balala-mysql:$TAG_VER \
./sql
```

### 构建mysql2es镜像
```bash
docker build \
--tag=czdesign2023/balala-mysql2es:$TAG_VER \
./mysql-es
```

### 构建前端镜像
```bash
docker build \
--tag=czdesign2023/balala-web:$TAG_VER \
../balala_client_vue
```

### docker-compose 启动
```shell
docker-compose up -d
```

### 拉取指定platform镜像并导出
#### amd64
```shell
docker pull --platform=linux/amd64 czdesign2023/balala-api:$TAG_VER
docker pull --platform=linux/amd64 czdesign2023/balala-mysql:$TAG_VER
docker pull --platform=linux/amd64 czdesign2023/balala-mysql2es:$TAG_VER
docker pull --platform=linux/amd64 czdesign2023/balala-web:$TAG_VER
docker pull --platform=linux/amd64 redis:7.2.3
docker pull --platform=linux/amd64 elasticsearch:7.17.22

docker save -o ./build/balala-images-amd64.tar \
  czdesign2023/balala-api:$TAG_VER \
  czdesign2023/balala-mysql:$TAG_VER \
  czdesign2023/balala-mysql2es:$TAG_VER \
  czdesign2023/balala-web:$TAG_VER \
  redis:7.2.3 \
  elasticsearch:7.17.22

```

#### arm64
```shell
docker pull --platform=linux/arm64 czdesign2023/balala-api:$TAG_VER
docker pull --platform=linux/arm64 czdesign2023/balala-mysql:$TAG_VER
docker pull --platform=linux/arm64 czdesign2023/balala-mysql2es:$TAG_VER
docker pull --platform=linux/arm64 czdesign2023/balala-web:$TAG_VER
docker pull --platform=linux/arm64 redis:7.2.3
docker pull --platform=linux/arm64 elasticsearch:7.17.22

docker save -o ./build/balala-images-arm64.tar \
  czdesign2023/balala-api:$TAG_VER \
  czdesign2023/balala-mysql:$TAG_VER \
  czdesign2023/balala-mysql2es:$TAG_VER \
  czdesign2023/balala-web:$TAG_VER \
  redis:7.2.3 \
  elasticsearch:7.17.22

```
