#用于构建docker镜像
FROM alpine
# 使用 Alpine Linux 作为基础镜像。Alpine 是一个轻量级的 Linux 发行版，非常适合用作 Docker 镜像的基础
WORKDIR /Initial
#设置容器中的工作目录为 /Initial
COPY ./target/project-user .
#将宿主机上 ./target/project-user 文件复制到容器中的当前工作目录 /Initial
COPY ./config/app-docker.yaml .
# 将宿主机上的 ./config/app-docker.yaml 文件复制到容器中的当前工作目录 /Initial
RUN  mkdir config && mv app-docker.yaml config/app.yaml
#在容器中创建一个 config 目录，然后将 app-docker.yaml 文件移动到 config 目录并重命名为 app.yaml
EXPOSE 8080 8881
#指定容器暴露的端口号。这里指定了两个端口号，分别是 8080 和 8881。
ENTRYPOINT ["./project-user"]
# 设置容器的入口点，即在容器启动时执行的命令。这里设置为运行 ./project-user 程序。



#合起来，这个 Dockerfile 的目的是创建一个基于 Alpine Linux 的 Docker 镜像，将 project-user 可执行文件和配置文件复制到容器中，然后设置容器的入口点为运行 project-user 程序，并且暴露了两个端口号（8080 和 8881）。这样通过构建这个 Docker 镜像，你就可以在容器中运行 project-user 服务了。
