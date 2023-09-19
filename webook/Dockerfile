# 基础镜像
FROM ubuntu:20.04
# 把编译好的 webook 打包进来这个镜像，放到工作目录 /app
COPY webook /app/webook
WORKDIR /app
ENTRYPOINT ["/app/webook"]