# movieSpider
 主要自动化下载电影的爬虫工具

# 前置条件
* 安装MySQL
```
# docker run -d --restart=always  --name mysql -e MYSQL_ROOT_PASSWORD=P@ssw0rd -p 3306:3306 mysql:5.7.18

```
* 安装Aria2

```shell
# PASSWORD=123.com
# Data_dir=$PWD/aria2-download
# docker run -d \
    --network host \
    --name aria2-pro \
    --restart unless-stopped \
    --log-opt max-size=1m \
    -e PUID=$UID \
    -e PGID=$GID \
    -e RPC_SECRET=${PASSWORD} \
    -e UMASK_SET=022 \
    -e RPC_PORT=6800 \
    -e LISTEN_PORT=6888 \
    -v $PWD/aria2-config:/config \
    -e CUSTOM_TRACKER_URL=https://trackerslist.com/all_aria2.txt \
    -e TZ=Asia/Shanghai \
    -v $Data_dir:/downloads \
    p3terx/aria2-pro


```

* 创建豆瓣想看列表

例如 [我的想看列表](https://movie.douban.com/people/251312920/wish)

# 命令帮助

* 生成配置

```sh
# movieSpider config 
MySQL:
  Host: 127.0.0.1
  Port: 3306
  Database: movie
  User: root
  Password: PASSWORD
...
...
TG:
  BotToken: "TOKEN"
  TgIDs: [ 123456 ]
```
* 初始化数据库

请提前配置好数据库连接信息，见配置文件
```shell
# movieSpider -f config.yaml config --init.db 
config file is config.yaml.
2022-10-16T20:49:40.494+0800    INFO    cmd/config.go:94        db: movie 数据库初始化完毕. 
```

* 仅运行爬虫

```shell
# movieSpider-linux -f config.yaml
```

* 同时运行TG机器人
需提前配置好 TG，见配置文件
```shell
# movieSpider-linux -f config.yaml --run.bot
```

# 自行编译
* 安装 golang 1.18 略
* git clone

```shell
git clone https://github.com/YouCD/movieSpider.git
```
* 编译

```shell
# cd movieSpider
# make build 
```
# systemd
```shell
name=movieSpider
WorkDir="/home/ycd/btspidery_data/movieSpider"
cat >/etc/systemd/system/${name}.service<<EOF

[Unit]
Description=${name}
Documentation=${name}
Wants=network-online.target
After=network-online.target

[Service]
WorkingDirectory=${WorkDir}
ExecStart=${WorkDir}/${name} -f config.yaml --run.bot
Restart=always

[Install]
WantedBy=multi-user.target
EOF

```