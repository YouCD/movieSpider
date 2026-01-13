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




wget https://raw.githubusercontent.com/YouCD/movieSpider/refs/heads/main/config.yaml -O  ${MOVIESPIDER_DIR}/config.yaml


sleep 5
echo -e "\033[32m[*] movieSpider的配置文件为: ${MOVIESPIDER_DIR}/config.yaml"
clear


echo "开始 下载docker-compose.yaml 以及相关的配置文件"
wget -q https://gh-proxy.com/raw.githubusercontent.com/YouCD/movieSpider/refs/heads/main/deployment/docker-compose.yaml

echo "开始 clash 的配置文件"
wget https://gh-proxy.com/raw.githubusercontent.com/Barabama/FreeNodes/main/nodes/v2rayshare.yaml -O ${CLASH_DIR}/clash.yaml

echo "启动 moviespider"
docker-compose -p moviespider up -d
