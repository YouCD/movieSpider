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

# quickStart

```shell
curl https://gh-proxy.com/raw.githubusercontent.com/YouCD/movieSpider/refs/heads/main/deployment/quickStart.sh| bash
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