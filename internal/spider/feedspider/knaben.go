package feedspider

import (
	"encoding/xml"
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/types"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
)

//nolint:revive
type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Atom    string   `xml:"atom,attr"`
	Channel struct {
		Text string `xml:",chardata"`
		Link struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Title struct {
			Text string `xml:",chardata"`
		} `xml:"title"`
		Description struct {
			Text string `xml:",chardata"`
		} `xml:"description"`
		Language struct {
			Text string `xml:",chardata"`
		} `xml:"language"`
		LastBuildDate struct {
			Text string `xml:",chardata"`
		} `xml:"lastBuildDate"`
		WebMaster struct {
			Text string `xml:",chardata"`
		} `xml:"webMaster"`
		Ttl struct {
			Text string `xml:",chardata"`
		} `xml:"ttl"`
		Item []struct {
			Text string `xml:",chardata"`
			Guid struct {
				Text        string `xml:",chardata"`
				IsPermaLink string `xml:"isPermaLink,attr"`
			} `xml:"guid"`
			Title struct {
				Text string `xml:",chardata"`
			} `xml:"title"`
			Link struct {
				Text string `xml:",chardata"`
			} `xml:"link"`
			Category []struct {
				Text string `xml:",chardata"`
			} `xml:"category"`
			Author struct {
				Text string `xml:",chardata"`
			} `xml:"author"`
			Description struct {
				Text string `xml:",chardata"`
			} `xml:"description"`
			PubDate struct {
				Text string `xml:",chardata"`
			} `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}
type Knaben struct {
	BaseFeeder
}

func NewFeedKnaben() *Knaben {
	return &Knaben{
		BaseFeeder: BaseFeeder{
			BaseFeed: types.BaseFeed{
				Scheduling: config.Config.Feed.Knaben.Scheduling,
				Url:        config.Config.Feed.Knaben.Url,
				UseIPProxy: config.Config.Feed.Knaben.UseIPProxy,
			},
			web: "knaben",
		},
	}
}

//nolint:nakedret
func (k *Knaben) Crawler() (videos []*types.FeedVideoBase, err error) {
	fd, err := k.FeedParser().ParseURL(k.Url)
	if err != nil {
		return nil, ErrFeedParseURL
	}
	log.Debugf("%s Data: %#v", strings.ToUpper(k.web), fd.String())

	resp, err := k.HTTPRequest(k.Url)
	if err != nil {
		return nil, fmt.Errorf("btbt new request,url: %s, err: %w", k.Url, err)
	}

	var respData Rss
	err = xml.Unmarshal(resp, &respData)
	if err != nil {
		return nil, fmt.Errorf("xml.Unmarshal,err: %w", err)
	}
	for _, s := range respData.Channel.Item {
		video := &types.FeedVideoBase{
			TorrentName: s.Title.Text,
			Web:         k.web,
		}
		// 标题
		// 种子连接
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(s.Description.Text))
		if err != nil {
			return nil, fmt.Errorf("getMovies goquery,err: %w", err)
		}
		var videoType string
		doc.Each(func(_ int, selection *goquery.Selection) {
			for i, s2 := range strings.Split(selection.Text(), "\n") {
				if i == 1 {
					videoType = strings.ToLower(s2)
				}
			}
		})

		switch {
		case strings.Contains(videoType, "tv"):
			video.Type = types.VideoTypeTV.String()
		case strings.Contains(videoType, "movies"):
			video.Type = types.VideoTypeMovie.String()
		default:
			continue
		}

		doc.Find("a").Each(func(_ int, selection *goquery.Selection) {
			magnet := selection.AttrOr("href", "")
			if strings.HasPrefix(magnet, "magnet") {
				video.Magnet = magnet
				videos = append(videos, video)
			}
		})
	}
	return
}
