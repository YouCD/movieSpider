package feedSpider

import (
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"testing"
)

var d = `

<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
<channel>
<title>RARBG</title>
<description>RARBG rss feed direct download</description>
<link>https://rarbg.to</link>
<lastBuildDate>Thu, 20 Oct 2022 12:40:40 CEST</lastBuildDate>
<copyright>(c) 2022 RARBG</copyright>

  <item>
  <title>Documentary.Now.S04E03.1080p.WEB.H264-WHOSNEXT[rartv]</title>
  <description>Documentary.Now.S04E03.1080p.WEB.H264-WHOSNEXT[rartv]</description>
  <link>magnet:?xt=urn:btih:296c5c7dc10cea07509b5c7b2e2b9a02f485584d&amp;dn=Documentary.Now.S04E03.1080p.WEB.H264-WHOSNEXT%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2850&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2780&amp;tr=udp%3A%2F%2Ftracker.tallpenguin.org%3A15790&amp;tr=udp%3A%2F%2Ftracker.thinelephant.org%3A12720</link>
  <guid>296c5c7dc10cea07509b5c7b2e2b9a02f485584d</guid>
  <pubDate>Thu, 20 Oct 2022 12:39:56 +0200</pubDate>
  </item>
    <item>
  <title>Documentary.Now.S04E01.1080p.WEB.H264-WHOSNEXT[rartv]</title>
  <description>Documentary.Now.S04E01.1080p.WEB.H264-WHOSNEXT[rartv]</description>
  <link>magnet:?xt=urn:btih:27aefaf0f1339d379169df135f6c82330d1e6f00&amp;dn=Documentary.Now.S04E01.1080p.WEB.H264-WHOSNEXT%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2890&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2780&amp;tr=udp%3A%2F%2Ftracker.tallpenguin.org%3A15790&amp;tr=udp%3A%2F%2Ftracker.thinelephant.org%3A12750</link>
  <guid>27aefaf0f1339d379169df135f6c82330d1e6f00</guid>
  <pubDate>Thu, 20 Oct 2022 12:39:05 +0200</pubDate>
  </item>
    <item>
  <title>Dark.Side.Of.Comedy.S01E10.WEBRip.x264-ION10</title>
  <description>Dark.Side.Of.Comedy.S01E10.WEBRip.x264-ION10</description>
  <link>magnet:?xt=urn:btih:8e9166a48ffdfb2b446f5c8dc595ca7651d539d6&amp;dn=Dark.Side.Of.Comedy.S01E10.WEBRip.x264-ION10&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2850&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2750&amp;tr=udp%3A%2F%2Ftracker.thinelephant.org%3A12730&amp;tr=udp%3A%2F%2Ftracker.slowcheetah.org%3A14770</link>
  <guid>8e9166a48ffdfb2b446f5c8dc595ca7651d539d6</guid>
  <pubDate>Thu, 20 Oct 2022 12:36:23 +0200</pubDate>
  </item>
    <item>
  <title>Question.Everything.S02E04.720p.HDTV.x264-ORENJI[rartv]</title>
  <description>Question.Everything.S02E04.720p.HDTV.x264-ORENJI[rartv]</description>
  <link>magnet:?xt=urn:btih:f46b13d64db14d303c50a9261c4206257159b6b9&amp;dn=Question.Everything.S02E04.720p.HDTV.x264-ORENJI%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2800&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2850&amp;tr=udp%3A%2F%2Ftracker.fatkhoala.org%3A13710&amp;tr=udp%3A%2F%2Ftracker.slowcheetah.org%3A14710</link>
  <guid>f46b13d64db14d303c50a9261c4206257159b6b9</guid>
  <pubDate>Thu, 20 Oct 2022 12:29:19 +0200</pubDate>
  </item>
    <item>
  <title>Question.Everything.S02E04.HDTV.x264-FQM[rartv]</title>
  <description>Question.Everything.S02E04.HDTV.x264-FQM[rartv]</description>
  <link>magnet:?xt=urn:btih:adcfee00abdbc77a6817bb97816e547268cfab4d&amp;dn=Question.Everything.S02E04.HDTV.x264-FQM%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2770&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2970&amp;tr=udp%3A%2F%2Ftracker.tallpenguin.org%3A15770&amp;tr=udp%3A%2F%2Ftracker.thinelephant.org%3A12800</link>
  <guid>adcfee00abdbc77a6817bb97816e547268cfab4d</guid>
  <pubDate>Thu, 20 Oct 2022 12:24:55 +0200</pubDate>
  </item>
    <item>
  <title>The.Challenge.S38E02.1080p.WEB.h264-BAE[rartv]</title>
  <description>The.Challenge.S38E02.1080p.WEB.h264-BAE[rartv]</description>
  <link>magnet:?xt=urn:btih:77d14622556a88a0f32b0c741c502c314b5f608c&amp;dn=The.Challenge.S38E02.1080p.WEB.h264-BAE%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2880&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2860&amp;tr=udp%3A%2F%2Ftracker.thinelephant.org%3A12730&amp;tr=udp%3A%2F%2Ftracker.fatkhoala.org%3A13770</link>
  <guid>77d14622556a88a0f32b0c741c502c314b5f608c</guid>
  <pubDate>Thu, 20 Oct 2022 12:23:43 +0200</pubDate>
  </item>
    <item>
  <title>Love.at.First.Lie.S01E03.1080p.WEB.h264-BAE[rartv]</title>
  <description>Love.at.First.Lie.S01E03.1080p.WEB.h264-BAE[rartv]</description>
  <link>magnet:?xt=urn:btih:4ec103928bc56229cfd6f1f552723d056c112cdd&amp;dn=Love.at.First.Lie.S01E03.1080p.WEB.h264-BAE%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2730&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2880&amp;tr=udp%3A%2F%2Ftracker.slowcheetah.org%3A14750&amp;tr=udp%3A%2F%2Ftracker.tallpenguin.org%3A15720</link>
  <guid>4ec103928bc56229cfd6f1f552723d056c112cdd</guid>
  <pubDate>Thu, 20 Oct 2022 12:23:43 +0200</pubDate>
  </item>
    <item>
  <title>VICE.News.Tonight.2022.10.19.1080p.WEB.h264-BAE[rartv]</title>
  <description>VICE.News.Tonight.2022.10.19.1080p.WEB.h264-BAE[rartv]</description>
  <link>magnet:?xt=urn:btih:be38bbc726338e3965a2b079ec4a04b8b0cc2b87&amp;dn=VICE.News.Tonight.2022.10.19.1080p.WEB.h264-BAE%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2730&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2900&amp;tr=udp%3A%2F%2Ftracker.thinelephant.org%3A12740&amp;tr=udp%3A%2F%2Ftracker.tallpenguin.org%3A15780</link>
  <guid>be38bbc726338e3965a2b079ec4a04b8b0cc2b87</guid>
  <pubDate>Thu, 20 Oct 2022 12:23:21 +0200</pubDate>
  </item>
    <item>
  <title>The.Challenge.S38E02.720p.WEB.h264-BAE[rartv]</title>
  <description>The.Challenge.S38E02.720p.WEB.h264-BAE[rartv]</description>
  <link>magnet:?xt=urn:btih:bbbd50473021732f149fe39bbb938ec57e245add&amp;dn=The.Challenge.S38E02.720p.WEB.h264-BAE%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2800&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2770&amp;tr=udp%3A%2F%2Ftracker.thinelephant.org%3A12730&amp;tr=udp%3A%2F%2Ftracker.tallpenguin.org%3A15720</link>
  <guid>bbbd50473021732f149fe39bbb938ec57e245add</guid>
  <pubDate>Thu, 20 Oct 2022 12:23:16 +0200</pubDate>
  </item>
    <item>
  <title>Dark.Side.Of.Comedy.S01E10.1080p.WEB.h264-BAE[rartv]</title>
  <description>Dark.Side.Of.Comedy.S01E10.1080p.WEB.h264-BAE[rartv]</description>
  <link>magnet:?xt=urn:btih:88f09c3785fb1f15dc938385b76693c98ed9c028&amp;dn=Dark.Side.Of.Comedy.S01E10.1080p.WEB.h264-BAE%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2860&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2840&amp;tr=udp%3A%2F%2Ftracker.thinelephant.org%3A12730&amp;tr=udp%3A%2F%2Ftracker.slowcheetah.org%3A14730</link>
  <guid>88f09c3785fb1f15dc938385b76693c98ed9c028</guid>
  <pubDate>Thu, 20 Oct 2022 12:23:03 +0200</pubDate>
  </item>
    <item>
  <title>Icons.Unearthed.S02E03.1080p.WEB.h264-BAE[rartv]</title>
  <description>Icons.Unearthed.S02E03.1080p.WEB.h264-BAE[rartv]</description>
  <link>magnet:?xt=urn:btih:1b94a15423305fda58eba115f7843727ba9c4954&amp;dn=Icons.Unearthed.S02E03.1080p.WEB.h264-BAE%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2880&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2850&amp;tr=udp%3A%2F%2Ftracker.fatkhoala.org%3A13730&amp;tr=udp%3A%2F%2Ftracker.tallpenguin.org%3A15760</link>
  <guid>1b94a15423305fda58eba115f7843727ba9c4954</guid>
  <pubDate>Thu, 20 Oct 2022 12:22:59 +0200</pubDate>
  </item>
    <item>
  <title>Love.at.First.Lie.S01E03.720p.WEB.h264-BAE[rartv]</title>
  <description>Love.at.First.Lie.S01E03.720p.WEB.h264-BAE[rartv]</description>
  <link>magnet:?xt=urn:btih:4f1863b79e0e78b9afc941d03108489a46698bab&amp;dn=Love.at.First.Lie.S01E03.720p.WEB.h264-BAE%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2960&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2810&amp;tr=udp%3A%2F%2Ftracker.slowcheetah.org%3A14780&amp;tr=udp%3A%2F%2Ftracker.tallpenguin.org%3A15770</link>
  <guid>4f1863b79e0e78b9afc941d03108489a46698bab</guid>
  <pubDate>Thu, 20 Oct 2022 12:22:44 +0200</pubDate>
  </item>
    <item>
  <title>One.of.Us.Is.Lying.S02E04.720p.WEB.H264-GLHF[rartv]</title>
  <description>One.of.Us.Is.Lying.S02E04.720p.WEB.H264-GLHF[rartv]</description>
  <link>magnet:?xt=urn:btih:6c5689ba4fb23863c92a3abcc3df6f887e40d909&amp;dn=One.of.Us.Is.Lying.S02E04.720p.WEB.H264-GLHF%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2990&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2900&amp;tr=udp%3A%2F%2Ftracker.thinelephant.org%3A12740&amp;tr=udp%3A%2F%2Ftracker.fatkhoala.org%3A13730</link>
  <guid>6c5689ba4fb23863c92a3abcc3df6f887e40d909</guid>
  <pubDate>Thu, 20 Oct 2022 12:21:33 +0200</pubDate>
  </item>
    <item>
  <title>One.of.Us.Is.Lying.S02E03.720p.WEB.H264-GLHF[rartv]</title>
  <description>One.of.Us.Is.Lying.S02E03.720p.WEB.H264-GLHF[rartv]</description>
  <link>magnet:?xt=urn:btih:17290adcce566c729387cdeba635b218dfbc0da3&amp;dn=One.of.Us.Is.Lying.S02E03.720p.WEB.H264-GLHF%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2940&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2950&amp;tr=udp%3A%2F%2Ftracker.fatkhoala.org%3A13760&amp;tr=udp%3A%2F%2Ftracker.tallpenguin.org%3A15770</link>
  <guid>17290adcce566c729387cdeba635b218dfbc0da3</guid>
  <pubDate>Thu, 20 Oct 2022 12:21:17 +0200</pubDate>
  </item>
    <item>
  <title>A.Friend.of.the.Family.S01E06.WEBRip.x264-ION10</title>
  <description>A.Friend.of.the.Family.S01E06.WEBRip.x264-ION10</description>
  <link>magnet:?xt=urn:btih:1a400e1e53deff031d42931e9cbf1288093f05f4&amp;dn=A.Friend.of.the.Family.S01E06.WEBRip.x264-ION10&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2840&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2830&amp;tr=udp%3A%2F%2Ftracker.fatkhoala.org%3A13770&amp;tr=udp%3A%2F%2Ftracker.tallpenguin.org%3A15720</link>
  <guid>1a400e1e53deff031d42931e9cbf1288093f05f4</guid>
  <pubDate>Thu, 20 Oct 2022 12:21:15 +0200</pubDate>
  </item>
    <item>
  <title>One.of.Us.Is.Lying.S02E08.720p.WEB.H264-GLHF[rartv]</title>
  <description>One.of.Us.Is.Lying.S02E08.720p.WEB.H264-GLHF[rartv]</description>
  <link>magnet:?xt=urn:btih:7e17ca057e6569f086eb356410ad36f381806618&amp;dn=One.of.Us.Is.Lying.S02E08.720p.WEB.H264-GLHF%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2730&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2790&amp;tr=udp%3A%2F%2Ftracker.thinelephant.org%3A12720&amp;tr=udp%3A%2F%2Ftracker.fatkhoala.org%3A13730</link>
  <guid>7e17ca057e6569f086eb356410ad36f381806618</guid>
  <pubDate>Thu, 20 Oct 2022 12:20:19 +0200</pubDate>
  </item>
    <item>
  <title>One.of.Us.Is.Lying.S02E05.720p.WEB.H264-GLHF[rartv]</title>
  <description>One.of.Us.Is.Lying.S02E05.720p.WEB.H264-GLHF[rartv]</description>
  <link>magnet:?xt=urn:btih:1ad5309421b9a4e5af98fb96bc2f4290e9249b4c&amp;dn=One.of.Us.Is.Lying.S02E05.720p.WEB.H264-GLHF%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2730&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2860&amp;tr=udp%3A%2F%2Ftracker.fatkhoala.org%3A13760&amp;tr=udp%3A%2F%2Ftracker.tallpenguin.org%3A15730</link>
  <guid>1ad5309421b9a4e5af98fb96bc2f4290e9249b4c</guid>
  <pubDate>Thu, 20 Oct 2022 12:19:52 +0200</pubDate>
  </item>
    <item>
  <title>One.of.Us.Is.Lying.S02E07.1080p.WEB.H264-GLHF[rartv]</title>
  <description>One.of.Us.Is.Lying.S02E07.1080p.WEB.H264-GLHF[rartv]</description>
  <link>magnet:?xt=urn:btih:12f7a0137de4241999f8d9672b3e8f30035bff86&amp;dn=One.of.Us.Is.Lying.S02E07.1080p.WEB.H264-GLHF%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2950&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2980&amp;tr=udp%3A%2F%2Ftracker.fatkhoala.org%3A13750&amp;tr=udp%3A%2F%2Ftracker.tallpenguin.org%3A15740</link>
  <guid>12f7a0137de4241999f8d9672b3e8f30035bff86</guid>
  <pubDate>Thu, 20 Oct 2022 12:19:04 +0200</pubDate>
  </item>
    <item>
  <title>One.of.Us.Is.Lying.S02E01.720p.WEB.H264-GLHF[rartv]</title>
  <description>One.of.Us.Is.Lying.S02E01.720p.WEB.H264-GLHF[rartv]</description>
  <link>magnet:?xt=urn:btih:e06855299f78c24d01f9f5d36c5013e1e96391fc&amp;dn=One.of.Us.Is.Lying.S02E01.720p.WEB.H264-GLHF%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2910&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2800&amp;tr=udp%3A%2F%2Ftracker.fatkhoala.org%3A13750&amp;tr=udp%3A%2F%2Ftracker.thinelephant.org%3A12720</link>
  <guid>e06855299f78c24d01f9f5d36c5013e1e96391fc</guid>
  <pubDate>Thu, 20 Oct 2022 12:18:32 +0200</pubDate>
  </item>
    <item>
  <title>One.of.Us.Is.Lying.S02E06.1080p.WEB.H264-GLHF[rartv]</title>
  <description>One.of.Us.Is.Lying.S02E06.1080p.WEB.H264-GLHF[rartv]</description>
  <link>magnet:?xt=urn:btih:f20f4e63ad7997a2dbf95316aee1fb7d3133698b&amp;dn=One.of.Us.Is.Lying.S02E06.1080p.WEB.H264-GLHF%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2970&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2860&amp;tr=udp%3A%2F%2Ftracker.fatkhoala.org%3A13760&amp;tr=udp%3A%2F%2Ftracker.slowcheetah.org%3A14780</link>
  <guid>f20f4e63ad7997a2dbf95316aee1fb7d3133698b</guid>
  <pubDate>Thu, 20 Oct 2022 12:16:33 +0200</pubDate>
  </item>
    <item>
  <title>One.of.Us.Is.Lying.S02E02.720p.WEB.H264-GLHF[rartv]</title>
  <description>One.of.Us.Is.Lying.S02E02.720p.WEB.H264-GLHF[rartv]</description>
  <link>magnet:?xt=urn:btih:390e951e875d6bc97f7bd8f7e6e99d5aa87db5bf&amp;dn=One.of.Us.Is.Lying.S02E02.720p.WEB.H264-GLHF%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2950&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2980&amp;tr=udp%3A%2F%2Ftracker.tallpenguin.org%3A15760&amp;tr=udp%3A%2F%2Ftracker.slowcheetah.org%3A14710</link>
  <guid>390e951e875d6bc97f7bd8f7e6e99d5aa87db5bf</guid>
  <pubDate>Thu, 20 Oct 2022 12:16:32 +0200</pubDate>
  </item>
    <item>
  <title>Vampire.Academy.S01E09.1080p.WEB.H264-GGEZ[rartv]</title>
  <description>Vampire.Academy.S01E09.1080p.WEB.H264-GGEZ[rartv]</description>
  <link>magnet:?xt=urn:btih:befa8d4a4c456b61cc24ef7891b00732eae93f7d&amp;dn=Vampire.Academy.S01E09.1080p.WEB.H264-GGEZ%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2820&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2910&amp;tr=udp%3A%2F%2Ftracker.thinelephant.org%3A12760&amp;tr=udp%3A%2F%2Ftracker.tallpenguin.org%3A15790</link>
  <guid>befa8d4a4c456b61cc24ef7891b00732eae93f7d</guid>
  <pubDate>Thu, 20 Oct 2022 12:16:31 +0200</pubDate>
  </item>
    <item>
  <title>One.of.Us.Is.Lying.S02E06.720p.WEB.H264-GLHF[rartv]</title>
  <description>One.of.Us.Is.Lying.S02E06.720p.WEB.H264-GLHF[rartv]</description>
  <link>magnet:?xt=urn:btih:24603a09c2186f5c4ddb1b9818d55a7a5277d66c&amp;dn=One.of.Us.Is.Lying.S02E06.720p.WEB.H264-GLHF%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2910&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2940&amp;tr=udp%3A%2F%2Ftracker.tallpenguin.org%3A15790&amp;tr=udp%3A%2F%2Ftracker.thinelephant.org%3A12760</link>
  <guid>24603a09c2186f5c4ddb1b9818d55a7a5277d66c</guid>
  <pubDate>Thu, 20 Oct 2022 12:16:29 +0200</pubDate>
  </item>
    <item>
  <title>One.of.Us.Is.Lying.S02E08.1080p.WEB.H264-GLHF[rartv]</title>
  <description>One.of.Us.Is.Lying.S02E08.1080p.WEB.H264-GLHF[rartv]</description>
  <link>magnet:?xt=urn:btih:ca2df6e4265ad2c4b15d4c91a3bee37819ae77db&amp;dn=One.of.Us.Is.Lying.S02E08.1080p.WEB.H264-GLHF%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2850&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2970&amp;tr=udp%3A%2F%2Ftracker.fatkhoala.org%3A13740&amp;tr=udp%3A%2F%2Ftracker.tallpenguin.org%3A15770</link>
  <guid>ca2df6e4265ad2c4b15d4c91a3bee37819ae77db</guid>
  <pubDate>Thu, 20 Oct 2022 12:16:12 +0200</pubDate>
  </item>
    <item>
  <title>One.of.Us.Is.Lying.S02E07.720p.WEB.H264-GLHF[rartv]</title>
  <description>One.of.Us.Is.Lying.S02E07.720p.WEB.H264-GLHF[rartv]</description>
  <link>magnet:?xt=urn:btih:fa7461bbe860c42ef6d153abb8c71b207bf3d624&amp;dn=One.of.Us.Is.Lying.S02E07.720p.WEB.H264-GLHF%5Brartv%5D&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2760&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2910&amp;tr=udp%3A%2F%2Ftracker.fatkhoala.org%3A13720&amp;tr=udp%3A%2F%2Ftracker.slowcheetah.org%3A14720</link>
  <guid>fa7461bbe860c42ef6d153abb8c71b207bf3d624</guid>
  <pubDate>Thu, 20 Oct 2022 12:16:09 +0200</pubDate>
  </item>
  </channel></rss>`

func Test_feed(t *testing.T) {
	//infos := feed(rarbgTVShows)
	//for _, v := range infos {
	//	fmt.Println(v)
	//}
	//compileRegex := regexp.MustCompile("(.*)\\.[sS][0-9][0-9]|[Ee][0-9][0-9]?\\.")
	//v := "Deadliest.Catch.The.Viking.Returns.S01E05.WEBRip.x264-ION10"
	//v2 := "Deadliest.Catch.The.Viking.Returns.S01.WEBRip.x264-ION10"
	//matchArr := compileRegex.FindStringSubmatch(v)
	//matchArr2 := compileRegex.FindStringSubmatch(v2)
	//
	//// 片名
	//fmt.Println("matchArr ", matchArr)
	//fmt.Println("matchArr2", matchArr2)
}

func Test_rarbg_Crawler(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")
	model.NewMovieDB()

}

func Test_rarbg_switchClient(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")

	model.NewMovieDB()
	//feedRarbg := NewFeedRarbg("https://rarbg.to/rssdd.php?categories=18;19;41", "*/2 * * * *", types.ResourceTV)
	////feedRarbg.useProxyClient()
	//feedRarbg.switchClient()
	//
	//feedRarbg.switchClient()

}
