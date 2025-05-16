# 使用官方 MySQL 镜像作为基础镜像
FROM mysql:8.0.35

# 将 SQL 初始化脚本复制到容器中
COPY *.sql /docker-entrypoint-initdb.d/

# 暴露 MySQL 端口
EXPOSE 3306
