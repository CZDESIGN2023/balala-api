# mysql-es 
- mysql同步es服务

```bash
docker run --rm \
  -v ./etc:/app/etc \
  -v mysql2es:/app \
  czdesign2023/mysql2es -log_level debug
```


# 本地部署

- docker没需要设置一下代理，没翻墙拉不到容器镜像
  docker配置 -》 Resources -》 Proxies -》 输入： http://127.0.0.1:7890 https://127.0.0.1:7890
  具体端口根据本地代理设置而定

- 把各容器关联到同一个网络 balala-net 下 
```shell
docker network create balala-net
docker network connect balala-net elasticsearch
docker network connect balala-net mysql
```


- 启动同步容器
```shell
docker run -d \
  --name balala-mysql2es \
  --net balala-net \
  -v $PWD/etc/river_local.toml:/app/etc/river.toml \
  -v balala-mysql2es:/app/var \
  czdesign2023/mysql2es -log_level debug
```
