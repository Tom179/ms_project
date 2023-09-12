chcp 65001
cd project-user
docker build -t project-user:latest . #自动寻找当前目录下的dockerFile构建project-user的镜像，标签名为latest
cd ..
docker-compose up -d