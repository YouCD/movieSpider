#!/bin/bash
set -e
clear
# 检查 Docker 和 Docker Compose 是否已安装
check_dependencies() {
    echo -e "\033[33m[*] 正在检查依赖项...\033[0m"

    # 检查 docker 命令是否存在
    if ! command -v docker &> /dev/null; then
        echo -e "\033[31m[错误] Docker 未安装，请先安装 Docker\033[0m"
        echo "安装 Docker 请参考: https://docs.docker.com/get-docker/"
        exit 1
    fi

    # 检查 docker-compose 命令是否存在
    if ! command -v docker-compose &> /dev/null; then
        echo -e "\033[31m[错误] Docker Compose 未安装，请先安装 Docker Compose\033[0m"
        echo "安装 Docker Compose 请参考: https://docs.docker.com/compose/install/"
        exit 1
    fi

    # 输出版本信息
    echo -e "\033[32m[*] Docker 版本: $(docker --version)\033[0m"
    echo -e "\033[32m[*] Docker Compose 版本: $(docker-compose --version)\033[0m"
    echo -e "\033[32m[*] 所有依赖项检查通过\033[0m"
    sleep 2
}

# 调用检查函数
check_dependencies
clear

#mkdir movieSpider_stack && cd movieSpider_stack

# 检查 .env 文件是否存在，不存在则创建
if [ ! -f ".env" ]; then
    echo -e "\033[32m[*] .env 文件不存在，正在创建...\033[0m"
    # LLM 配置 - 从用户输入获取
    read -p "请输入 LLM API Key: " input_LLM_API_KEY
    if [ -z "$input_LLM_API_KEY" ]; then
        echo "警告: LLM API Key 为空，LLM 功能将不可用"
        export LLM_API_KEY=""
    else
        export LLM_API_KEY=$input_LLM_API_KEY
    fi

    read -p "请输入 LLM Model: " input_LLM_MODEL
    export LLM_MODEL=${input_LLM_MODEL}

    read -p "请输入 LLM Base URL: " input_LLM_BASE_URL
    export LLM_BASE_URL=${input_LLM_BASE_URL}

    cat > .env <<EOF
LLM_API_KEY="${input_LLM_API_KEY}"
LLM_MODEL="${input_LLM_MODEL}"
LLM_BASE_URL="${input_LLM_BASE_URL}"
MYSQL_PASSWORD=P@ssw0rd
MYSQL_PORT=3306
ARIA2_PASSWORD=whVi763s5QrctiiyUxIs
ARIA2_CONFIG_DIR=$PWD/aria2/config
ARIA2_DATA_DIR=$PWD/aria2/data
ARIA2_PORT=6800
MOVIESPIDER_DIR=$PWD/movieSpider
DOUBAN_URL=https://movie.douban.com/people/251312920/wish
CLASH_DIR=$PWD/clash
EOF
    echo -e "\033[32m[*] .env 文件创建完成: $(pwd)/.env\033[0m"
else
    echo -e "\033[33m[*] .env 文件已存在，跳过创建\033[0m"
fi

# 从 .env 文件加载环境变量
set -a
source .env

set +a

sleep 5
echo -e "\033[73m[*] 创建目录： ${MOVIESPIDER_DIR} ${CLASH_DIR}"
mkdir -p ${MOVIESPIDER_DIR} ${CLASH_DIR}




cat > ${MOVIESPIDER_DIR}/config.yaml<<EOF
MySQL:
  # 这个地址是docker里面的地址
  Host: moviespider_mysql
  Port: 3306
  Database: movie
  User: root
  Password: P@ssw0rd

DouBan:
  # 豆瓣电影想看清单
  Scheduling: "*/10 * * * *"
  DouBanList:
    - Url: "${DOUBAN_URL}"

ExcludeWords:
  - 720p
  - dvsux
  - 480p
  - 360p
  - .dv.
  - .dolby.vision

Feed:
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
      Url: "https://www.1337x.to/popular-movies"
      UseIPProxy: true
    - Scheduling: "*/2 * * * *"
      ResourceType: tv
      Url: "https://1337x.to/popular-tv"
      UseIPProxy: true
  ThePirateBay:
    Scheduling: "*/3 * * * *"
    Url: "https://thepiratebay.party/rss//top100/200"
  Knaben:
    Scheduling: "*/3 * * * *"
    Url: "https://rss.knaben.org////hidexxx"
  TheRarbg:
    - Scheduling: "*/3 * * * *"
      Url: "https://therarbg.to/get-posts/category:TV:time:10D/"
      ResourceType: tv
    - Scheduling: "*/3 * * * *"
      Url: "https://therarbg.to/get-posts/category:Movies:time:10D/"
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
  DHTThread: 0 # DHT网络爬虫线程数, 0关闭
  Timeout: 123
  ProxyUrl: http://moviespider_clash:7890

LLM:
  ApiKey: "${LLM_API_KEY}"
  Model: "${LLM_MODEL}"
  BaseUrl: "${LLM_BASE_URL}"
# Downloader 下载
Downloader:
  Scheduling: "*/60 * * * *"
  # 使用哪个 Aria2 下载
  Aria2Label: "home"

# Aria2 下载服务器
Aria2cList:
  - Url: "http://moviespider_aria2:6800"
    Token: P@ssw0rd
    Label: home

# 如果没有Telegram 就请忽略
TG:
  # Telegram 机器人 token
#  BotToken: "TOKEN"
#   能正常访问机器人的Telegram用户
#  TgIDs: [ 123456 ]
EOF


sleep 5
echo -e "\033[32m[*] movieSpider的配置文件为: ${MOVIESPIDER_DIR}/config.yaml"
clear


echo "开始 下载docker-compose.yaml 以及相关的配置文件"
wget -q https://gh-proxy.com/raw.githubusercontent.com/YouCD/movieSpider/refs/heads/main/deployment/docker-compose.yaml

echo "开始 clash 的配置文件"
wget https://gh-proxy.com/raw.githubusercontent.com/Barabama/FreeNodes/main/nodes/v2rayshare.yaml -O ${CLASH_DIR}/clash.yaml

echo "启动 moviespider"
docker-compose -p moviespider up -d
