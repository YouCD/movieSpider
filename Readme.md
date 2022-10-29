# movieSpider
  
 自动化下载电影的爬虫工具

支持的`Torrent`站点
* [bt4g](https://bt4g.org)
* [btbt](https://www.btbtt12.com/)
* [eztv](https://eztv.re)
* [glodls](https://glodls.to)
* [knaben](https://rss.knaben.eu)
* [rarbg](http://rarbg.to/)
* [TGx](https://tgx.rs)
* [torlock](https://www.torlock.com)
* [magnetdl](https://www.magnetdl.com)

如果要使用 TG bot， 请添加如下指令

Telegram bot 指令如下

```shell
movie_download - 下载电影  电影名  清晰度
report_download - 报告下载状态
report_feedvioes - 报告Feed资源
```

# 部署

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
* `jhao104/proxy_pool` 免费代理

运行`redis`
```shell
docker run -d --name redis --restart=always -p 6379:6379  redis
```

运行 `jhao104/proxy_pool`
有Redis密码
```shell
Redis_Password=123456
Redis_Host=127.0.0.1:6379
docker run -d --name proxy_pool --env DB_CONN=redis://:${Redis_Password}@${Redis_Host}/0 -p 5010:5010 jhao104/proxy_pool:latest
```
无Redis密码
```shell
Redis_Host=127.0.0.1:6379
docker run -d --name proxy_pool --env DB_CONN=redis://${Redis_Host}/0 -p 5010:5010 jhao104/proxy_pool:latest
```


* 一把梭
在一把梭之前请准备好配置文件
```shell
mkdir $PWD/movieSpiderCore
cp config.yaml $PWD/movieSpiderCore
```



```shell
export Mysql_Password=P@ssw0rd
export Mysql_Port=3307
export Mysql_Database=movie
export Aria2_Password=P@ssw0rd
export Aria2_ConfigDir=$PWD/movieSpiderCore/aria2/config
export Aria2_DataDir=$PWD/movieSpiderCore/aria2/data
export Aria2_Port=6801
export UID=$UID
export GID=$GID
export MovieSpider_Dir=$PWD/movieSpiderCore/
docker-compose up -d 

```








* 创建豆瓣想看列表

例如 [我的想看列表](https://movie.douban.com/people/251312920/wish)

# 命令帮助

* 生成配置

```sh
# movieSpiderCore config 
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
# movieSpiderCore -f config.yaml config --init.db 
config file is config.yaml.
2022-10-16T20:49:40.494+0800    INFO    cmd/config.go:94        db: movie 数据库初始化完毕. 
```

* 仅运行爬虫

```shell
# movieSpiderCore-linux -f config.yaml
```

* 同时运行TG机器人
需提前配置好 TG，见配置文件
```shell
# movieSpiderCore-linux -f config.yaml --run.bot
```

# 自行编译
* 安装 golang 1.18 

参考 [一键安装](https://github.com/Jrohy/go-install)

* git clone

```shell
git clone https://github.com/YouCD/movieSpider.git
```
* 编译

```shell
# cd movieSpiderCore
# make build 
```
# systemd
```shell
name=movieSpiderCore
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