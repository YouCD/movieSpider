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
* [TpbpirateProxy](https://thepiratebay.party/rss//top100/200)

如果要使用 TG bot， 请添加如下指令

Telegram bot 指令如下

```shell
movie_download - 下载电影  电影名  清晰度
report_download - 报告下载状态
report_feedvioes - 报告Feed资源
```

# 部署




## 创建豆瓣想看列表

例如 [我的想看列表](https://movie.douban.com/people/251312920/wish)


## 定义环境变量
```api
export Mysql_Password=P@ssw0rd
export Mysql_Port=3306
export Aria2_Password=whVi763s5QrctiiyUxIs
export Aria2_ConfigDir=$PWD/aria2/config
export Aria2_DataDir=$PWD/aria2/data
export Aria2_Port=6800
export UID=$UID
export GID=$GID
export MovieSpider_Dir=$PWD/movieSpider
# 豆瓣想看列表
export DoubanUrl=https://movie.douban.com/people/251312920/wish

mkdir -p ${MovieSpider_Dir}

```


## 准备配置文件
* IpProxyPool
```
cat >IpProxyPool/config.yaml<<EOF
# server configuration
system:
  appName: ProxyPool
  httpAddr: 0.0.0.0
  httpPort: 5010

database:
  host: moviespider_mysql
  port: ${Mysql_Port}
  dbName: IpProxyPool
  username: root
  password: ${Mysql_Password}
log:
  filePath: logs
  fileName: run.log
  level: debug
  mode: file
EOF

```



* moviespider

```shell
cat >${MovieSpider_Dir}/config.yaml<<EOF
MySQL:
  # 这个地址是docker里面的地址
  Host: moviespider_mysql
  Port: ${Mysql_Port}
  Database: movie
  User: root
  Password: ${Mysql_Password}

Douban:
  # 豆瓣电影想看清单
  DoubanUrl: ${DoubanUrl}
  Scheduling: "*/10 * * * *"
  # 豆瓣 Cookie
  # Cookie: ''
Feed:
  # 代理池 https://github.com/YouCD/IpProxyPool
  ProxyPool: "http://moviespider_proxy:5010"
  BTBT:
    Scheduling: "*/5 * * * *"
  EZTV:
    Scheduling: "*/5 * * * *"
  GLODLS:
    Scheduling: "*/3 * * * *"
  TGX:
    Scheduling: "*/3 * * * *"
  RARBG:
    - Scheduling: "*/3 * * * *"
      ResourceType: movie
    - Scheduling: "*/2 * * * *"
      ResourceType: tv
  TORLOCK:
    - Scheduling: "*/3 * * * *"
      ResourceType: movie
    - Scheduling: "*/2 * * * *"
      ResourceType: tv
  MAGNETDL:
    - Scheduling: "*/3 * * * *"
      ResourceType: movie
    - Scheduling: "*/2 * * * *"
      ResourceType: tv
  TPBPIRATEPROXY:
    Scheduling: "*/3 * * * *"
Global:
  LogLevel: info
  Report: true

# Downloader 下载
Downloader:
  Scheduling: "*/60 * * * *"
  # 使用哪个 Aria2 下载
  Aria2Label: "home"

# Aria2 下载服务器
Aria2cList:
  - Url: "http://moviespider_aria2:6800"
    Token: ${Aria2_Password}
    Label: home

# 如果没有Telegram 就请忽略
#TG:
  # Telegram 网络代理
#  Proxy:
#    Url: socks5://127.0.0.1:1080
#    Enable: false
  # Telegram 机器人 token
#  BotToken: "TOKEN"
#   能正常访问机器人的Telegram用户
#  TgIDs: [ 123456 ]
EOF




```
### 运行

```shell
docker-compose -p moviespider up
```

# quickStart

```shell
curl https://raw.githubusercontent.com/YouCD/movieSpider/main/deployment/quickStart.sh| bash
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
ExecStart=${WorkDir}/${name} -f config.yaml
Restart=always

[Install]
WantedBy=multi-user.target
EOF

```