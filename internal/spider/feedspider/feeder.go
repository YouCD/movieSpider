package feedspider

import (
	"context"
	"errors"
	"fmt"
	"io"
	"movieSpider/internal/httpclient"
	"movieSpider/internal/types"
	"net/http"

	"github.com/mmcdole/gofeed"
)

var (
	ErrNoFeedData   = errors.New("no feed data")
	ErrFeedParseURL = errors.New("feed url解析失败")
)

//nolint:inamedparam
type Feeder interface {
	Scheduling() string
	WebName() string
	URL() string
	Crawler() ([]*types.FeedVideoBase, error)
}
type Crawler func() ([]*types.FeedVideo, error)

type BaseFeeder struct {
	types.BaseFeed
	web string
}

// func (b *BaseFeeder) SetHttpClient(HTTPRequest *http.Client) {
// 	b.mux.Lock()
// 	defer b.mux.Unlock()
// 	b.HTTPClient = HTTPRequest
// }

func (b *BaseFeeder) HTTPClient() *http.Client {
	return httpclient.HTTPClient
}

func (b *BaseFeeder) HTTPClientIPProxyPool() *http.Client {
	var count int
	for {
		count++
		proxyClient, _ := httpclient.NewIPProxyPoolHTTPClient(b.URL())
		if proxyClient == nil {
			if count > 3 {
				return b.HTTPClient()
			}
			continue
		}
		return proxyClient
	}
}
func (b *BaseFeeder) HTTPClientDynamic() *http.Client {
	if b.UseIPProxy {
		return b.HTTPClientIPProxyPool()
	}
	return b.HTTPClient()
}
func (b *BaseFeeder) FeedParser() *gofeed.Parser {
	fp := gofeed.NewParser()
	fp.Client = b.HTTPClientDynamic()
	return fp
}
func (b *BaseFeeder) FeedParserUserAgent(userAgent string) *gofeed.Parser {
	fp := gofeed.NewParser()
	fp.Client = b.HTTPClientDynamic()
	fp.UserAgent = userAgent
	return fp
}

func (b *BaseFeeder) Crawler() ([]*types.FeedVideo, error) {
	return nil, nil
}

func (b *BaseFeeder) URL() string {
	return b.BaseFeed.Url
}

func (b *BaseFeeder) Scheduling() string {
	return b.BaseFeed.Scheduling
}

func (b *BaseFeeder) WebName() string {
	return b.web
}

type FeederAbstractFactory interface {
	CreateFeeder(args ...interface{}) Feeder
}

func (b *BaseFeeder) HTTPRequest(url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("HTTPRequest new request,err: %w", err)
	}
	resp, err := b.HTTPClientDynamic().Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTPRequest do,err: %w", err)
	}
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()
	//nolint:wrapcheck
	return io.ReadAll(resp.Body)
}
