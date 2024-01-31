#!/bin/bash
clear
echo -e "\033[31m[*] 请确保已经安装了docker和docker-compose"
sleep 5
clear

echo -e "\033[31m[*] 定义的环境变量如下:"
echo -e "\033[31mexport Mysql_Password=P@ssw0rd
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
export DoubanUrl=https://movie.douban.com/people/251312920/wish




sleep 5
echo -e "\033[73m[*] 创建目录： ${MovieSpider_Dir}"
mkdir -p ${MovieSpider_Dir}
clear





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
  DoubanUrl:
    - Url: ${DoubanUrl}
  Scheduling: "*/10 * * * *"
Feed:
  BTBT:
      Scheduling: "*/5 * * * *"
  EZTV:
    Scheduling: "*/5 * * * *"
    MirrorSite: "https://eztvx.to"
  GLODLS:
    Scheduling: "*/3 * * * *"
    MirrorSite: "https://gtso.cc"
  TGX:
    Scheduling: "*/3 * * * *"
    MirrorSite: "https://tgx.rs"
  TORLOCK:
    - Scheduling: "*/3 * * * *"
      ResourceType: movie
      MirrorSite: "https://torlock.unblockit.date"
    - Scheduling: "*/2 * * * *"
      ResourceType: tv
      MirrorSite: "https://torlock.unblockit.date"
  MAGNETDL:
    - Scheduling: "*/3 * * * *"
      ResourceType: movie
      MirrorSite: "https://magnetdl.abcproxy.org"
    - Scheduling: "*/2 * * * *"
      ResourceType: tv
      MirrorSite: "https://magnetdl.abcproxy.org"
  TPBPIRATEPROXY:
    Scheduling: "*/3 * * * *"
    MirrorSite: "https://thepiratebay10.info"
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


sleep 5
echo -e "\033[32m[*] movieSpider配置文件内容如下:"
cat ${MovieSpider_Dir}/config.yaml
clear



sleep 5
echo "开始下载docker-compose.yaml以及相关的Dockerfile"
wget -q https://raw.githubusercontent.com/YouCD/movieSpider/main/deployment/docker-compose.yaml
wget -q https://raw.githubusercontent.com/YouCD/movieSpider/main/deployment/moviespider_Dockerfile
echo "启动 moviespider"
docker-compose -p moviespider up
