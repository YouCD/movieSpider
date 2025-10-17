# movieSpider

自动化下载电影的爬虫工具

## 磁力站点

支持的`Torrent`站点
* ~~[[bt4g](https://bt4g.org)~~
* ~~[btbt](https://www.1lou.me/forum-1.htm)~~
* [eztv](https://eztv.re)
* [glodls](https://glodls.to)
* [knaben](https://rss.knaben.eu)
* ~~[rarbg](http://rarbg.to/)~~
* ~~[TGx](https://tgx.rs)~~
* [torlock](https://www.torlock.com)
* ~~[magnetdl](https://www.magnetdl.com)~~
* [ThePirateBay](https://thepiratebay.org/search.php?q=top100:200)
* [1337x](https://1337x.to/)
* ~~[rarbg2](https://en.rarbg2.xyz)~~
* [TheRARBG](https://therarbg.com/)
* ~~[EXT](https://extto.com/)~~ 效果不理想
* [uindex](https://uindex.org/)
* [Ilcorsaronero](https://ilcorsaronero.link/)

## 使用TG

如果要使用 TG bot， 请添加如下指令

Telegram bot 指令如下

```shell
report_download - 报告下载状态
report_feedvioes - 报告Feed资源
```

## 功能
* [x] 自动爬取上述站点的资源
* [x] 自动通过`Aria2`下载
* [x] 支持`Telegram`通知：
  * 下载时通知
  * 电影上映通知
  * 电影下载完成通知
* [x] 支持`Telegram`机器人命令：
  
  有白名单，只有白名单的用户才能使用
  ```
  report_download - 报告下载状态
  report_feedvioes - 报告Feed资源
  ```
* [x] 集成DHT网络：感谢[nbdy/dhtc](https://github.com/nbdy/dhtc)
* [x] ~~微调Qwen2.5-0.5B模型： 提高种子名规范化能力~~
* [x] 微调Qwen3-0.6B模型： 提高种子名规范化能力
* [x] ~~引入[CloudflareBypassForScraping](https://github.com/sarperavci/CloudflareBypassForScraping)： 自动处理Cloudflare防护~~


* [ ] 自动化下载字幕
* [ ] 基本网页展示


## 时序图

![时序图](https://www.plantuml.com/plantuml/png/bL9FRn915B_FfvYavruyqiGKqaYQwCdUcXucxeYRBBk4pkgrj4XRAOZf7nMGgbYtYO5IZIqnOhyCRsOu-HKSEhjPfWMJSY2Rztxlz-UzMLnd9C_yh8EOpHoAizvpcpxEsHRD8vHJBVk9dICkmxgsvFafuXeDULuI-ssZj4GOIP5rQZ8yeDJIB4RPybwuZaTkbfDvczmpAYjfX2PTuFdxeNW2o-ebl3xYMpz4iopgHN7m4mRd_B37ArsaCNXUmIV74pX4zLz2vT1alWMLGPktQWRjpO4eT093vnrideOm5CSUeqIx0AyQOl16V7klaDqwVtkXpKEj25-2CLe528SdqTS1OZBbwUe06gtGqywH-f112NQ3SDRKJIudCA-UZZOAK3s6e2o22dHgKAkWvF8CPXaKvTwMCXDZ9pbdHyxaFJytb_OS2yyzj3FFMChnAHvjrENLjhQhkFZnYJpx_FVt7dpTvAtcKyMinuVz3c_kHjIh-gRwI6kgDrxEXq-YqKLO_UB7OvnWkjZTDHQkHIASYXphHLl2CIkUTOlRLlNpZawY5won2__D1RtIgbPPun-tPRStxF-2o-_3VbaE1kD9JHmhX9juH38kPQFzFxnwCWsyaewyy8iTNZX3rBNGelhpiyu_)


# 部署




## 创建豆瓣想看列表

例如 [我的想看列表](https://movie.douban.com/people/251312920/wish)


## 定义环境变量
```shell
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


* moviespider

```shell
cat > ${MovieSpider_Dir}/config.yaml<<'EOF'
MySQL:
  # 这个地址是docker里面的地址
  Host: 127.0.0.1
  Port: 3306
  Database: movie
  User: root
  Password: P@ssw0rd

DouBan:
  # 豆瓣电影想看清单
  Scheduling: "*/10 * * * *"
  DouBanList:
    - Url: "https://movie.douban.com/people/251312920/wish"
    - Url: "https://movie.douban.com/people/271517237/wish"

ExcludeWords:
  - 720p
  - dvsux
  - 480p
  - 360p
  #- hdr
  - .dv.
  - .dolby.vision

Feed:
  EZTV:
    Scheduling: "*/5 * * * *"
    Url: "https://eztvx.to/ezrss.xml"
    UseIPProxy: true
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
    Url: "https://thepiratebay.party/rss/top100/200"
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
  # 免费的网络代理池 https://github.com/YouCD/IpProxyPool
  IPProxyPool: "http://127.0.0.1:3001"
  DHTThread: 0 # DHT网络爬虫线程数, 0关闭
  NameParserModel: http://127.0.0.1:8000 # 使用模型进行解析种子名称
  Timeout: 60
# Downloader 下载
Downloader:
  Scheduling: "*/60 * * * *"
  # 使用哪个 Aria2 下载
  Aria2Label: "home"

# Aria2 下载服务器
Aria2cList:
  - Url: "http://127.0.0.1e:6800"
    Token: P@ssw0rd
    Label: home

# 如果没有Telegram 就请忽略
TG:
  # Telegram 机器人 token
#  BotToken: "TOKEN"
#   能正常访问机器人的Telegram用户
#  TgIDs: [ 123456 ]
  # 独立的代理地址
#  ProxyUrl: socks5://192.168.1.188:20170
EOF

```
### 运行

```shell
docker-compose -p moviespider up
```

# quickStart

```shell
curl https://raw.githubusercontent.com/YouCD/movieSpider/refs/heads/main/deployment/quickStart.sh| bash
```




# systemd
```shell
name=core
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

## 截图
* 上映通知

  ![photo_2023-06-17_11-22-43.jpg](doc/image/photo_2023-06-17_11-30-45.jpg)

* 下载通知

  ![photo_2023-06-17_11-22-43.jpg](doc/image/photo_2023-06-17_11-22-43.jpg)

* 下载完毕

  ![photo_2023-06-17_11-22-43.jpg](doc/image/photo_2023-06-18_07-59-39.jpg)


## 其他

### Q&A 
* Q：通知渠道
  
  A：目前通知的渠道只适配了TG


* Q：清晰度
  
  A：目前只支持 1080p 2160p等资源









### RARBG
 RARBG 时代结束了
```
Hello guys,
We would like to inform you that we have decided to shut down our site.
The past 2 years have been very difficult for us - some of the people in our team died due to covid complications,
others still suffer the side effects of it - not being able to work at all.
Some are also fighting the war in Europe - ON BOTH SIDES.
Also, the power price increase in data centers in Europe hit us pretty hard.
Inflation makes our daily expenses impossible to bare.
Therefore we can no longer run this site without massive expenses that we can no longer cover out of pocket.
After an unanimous vote we've decided that we can no longer do it.
We are sorry :(
Bye


```