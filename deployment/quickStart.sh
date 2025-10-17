#!/bin/bash
clear
echo -e "\033[31m[*] 请确保已经安装了docker和docker-compose"
sleep 5
clear

mkdir movieSpider_stack && cd movieSpider_stack

echo -e "\033[31m[*] 定义的环境变量如下:"
echo -e "\033[31m export Mysql_Password=P@ssw0rd
export Mysql_Port=3306
export Aria2_Password=whVi763s5QrctiiyUxIs
export Aria2_ConfigDir=$PWD/aria2/config
export Aria2_DataDir=$PWD/aria2/data
export Aria2_Port=6800
export MovieSpider_Dir=$PWD/movieSpider
export DoubanUrl=https://movie.douban.com/people/251312920/wish"
clear





export Mysql_Password=P@ssw0rd
export Mysql_Port=3306
export Aria2_Password=whVi763s5QrctiiyUxIs
export Aria2_ConfigDir=$PWD/aria2/config
export Aria2_DataDir=$PWD/aria2/data
export Aria2_Port=6800
export MovieSpider_Dir=$PWD/movieSpider
export IpProxyPool_Dir=$PWD/IpProxyPool
export UID=$UID
export GID=$GID
export DoubanUrl=https://movie.douban.com/people/251312920/wish




sleep 5
echo -e "\033[73m[*] 创建目录： ${MovieSpider_Dir}"
mkdir -p ${MovieSpider_Dir}
mkdir -p ${IpProxyPool_Dir}
clear





cat > ${MovieSpider_Dir}/config.yaml<<EOF
MySQL:
  # 这个地址是docker里面的地址
  Host: moviespider_mysql
  Port: ${Mysql_Port}
  Database: movie
  User: root
  Password: P@ssw0rd

DouBan:
  # 豆瓣电影想看清单
  Scheduling: "*/10 * * * *"
  DouBanList:
    - Url: "${DoubanUrl}"

ExcludeWords:
  - 720p
  - dvsux
  - 480p
  - 360p
  #- hdr
  - .dv.
  - .dolby.vision

Feed:
  BTBT:
    Scheduling: "*/5 * * * *"
    Url: "https://www.1lou.me/forum-1.htm"
  EZTV:
    Scheduling: "*/5 * * * *"
    Url: "https://eztvx.to/ezrss.xml"
  GLODLS:
    Scheduling: "*/3 * * * *"
    Url: "https://glodls.to/rss.php?cat=1,41"
    UseIPProxy: true
  TORLOCK:
    - Scheduling: "*/3 * * * *"
      ResourceType: movie
      Url: "https://www.torlock.com/movies/rss.xml"
      UseIPProxy: true
    - Scheduling: "*/2 * * * *"
      ResourceType: tv
      Url: "https://www.torlock.com/television/rss.xml"
      UseIPProxy: true
  Web1337x:
    - Scheduling: "*/3 * * * *"
      ResourceType: movie
      Url: "https://1337x.to/popular-movies"
      UseIPProxy: true
    - Scheduling: "*/2 * * * *"
      ResourceType: tv
      Url: "https://1337x.to/popular-tv"
      UseIPProxy: true
  ThePirateBay:
    Scheduling: "*/3 * * * *"
    Url: "https://thepiratebay.org/search.php?q=top100:200"
    UseIPProxy: true
  Knaben:
    Scheduling: "*/3 * * * *"
    Url: "https://rss.knaben.org////hidexxx"
  TheRarbg:
    - Scheduling: "*/3 * * * *"
      Url: "https://therarbg.to/api/v1/recommendation-list/tv/"
      ResourceType: tv
    - Scheduling: "*/3 * * * *"
      Url: "https://therarbg.to/api/v1/recommendation-list/Movies/"
      ResourceType: movie
  Uindex:
    - Scheduling: "*/3 * * * *"
      Url: "https://uindex.org/top.php?c=2"
      ResourceType: tv
    - Scheduling: "*/3 * * * *"
      Url: "https://uindex.org/top.php?c=1"
      ResourceType: movie
  Ilcorsaronero:
    - Scheduling: "*/3 * * * *"
      Url: "https://ilcorsaronero.link/cat/serie-tv"
      ResourceType: tv
    - Scheduling: "*/3 * * * *"
      Url: "https://ilcorsaronero.link/cat/film"
      ResourceType: movie
Global:
  LogLevel: info
  Report: true
  # 网络代理池
  IPProxyPool: "http://moviespider_proxy:3001"
  DHTThread: 0 # DHT网络爬虫线程数, 0关闭
  NameParserModel: http://moviespider_name_parser_model:8000 # 使用模型进行解析种子名称
  Timeout: 60

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
#  ProxyUrl: socks5://127.0.0.1:1080
#    Enable: false
  # Telegram 机器人 token
#  BotToken: "TOKEN"
#   能正常访问机器人的Telegram用户
#  TgIDs: [ 123456 ]
EOF


sleep 5
echo -e "\033[32m[*] movieSpider配置文件内容如下:"
cat ${MovieSpider_Dir}/config.yaml
clear



#sleep 5
echo "开始 构建镜像"
#docker build  -f https://raw.githubusercontent.com/YouCD/movieSpider/main/deployment/moviespider_proxy_Dockerfile -t moviespider_proxy .
#docker build  -f https://raw.githubusercontent.com/YouCD/movieSpider/main/deployment/moviespider_Dockerfile -t moviespider_core .

echo "开始 下载docker-compose.yaml 以及相关的配置文件"
wget -q https://raw.githubusercontent.com/YouCD/movieSpider/refs/heads/main/deployment/docker-compose.yaml
wget -q https://raw.githubusercontent.com/YouCD/IpProxyPool/refs/heads/main/conf/config.yaml -O ${IpProxyPool_Dir}/config.yaml
sed -i "s/127.0.0.1/moviespider_mysql/" ${IpProxyPool_Dir}/config.yaml

echo "启动 moviespider"
docker-compose -p moviespider up -d
