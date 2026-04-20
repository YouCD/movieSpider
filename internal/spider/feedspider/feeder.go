package feedspider

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"movieSpider/internal/httpclient"
	"movieSpider/internal/types"

	"github.com/mmcdole/gofeed"
)

//nolint:inamedparam
type Feeder interface {
	Scheduling() string
	WebName() string
	URL() string
	Crawler(ctx context.Context) ([]*types.FeedVideoBase, error)
}
type Crawler func() ([]*types.FeedVideo, error)

type BaseFeeder struct {
	types.BaseFeed
	web string
}

func (b *BaseFeeder) HTTPClient() *http.Client {
	return httpclient.HTTPClient
}

func (b *BaseFeeder) HTTPClientIPProxyPool(ctx context.Context) *http.Client {
	return httpclient.NewProxyHTTPClient(ctx)
}

func (b *BaseFeeder) HTTPClientDynamic(ctx context.Context) *http.Client {
	if b.UseIPProxy {
		return b.HTTPClientIPProxyPool(ctx)
	}
	return b.HTTPClient()
}

func (b *BaseFeeder) FeedParser(ctx context.Context) *gofeed.Parser {
	fp := gofeed.NewParser()
	fp.Client = b.HTTPClientDynamic(ctx)
	return fp
}

func (b *BaseFeeder) FeedParserUserAgent(ctx context.Context, userAgent string) *gofeed.Parser {
	fp := gofeed.NewParser()
	fp.Client = b.HTTPClientDynamic(ctx)
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

func (b *BaseFeeder) HTTPRequest(ctx context.Context, urlStr string) ([]byte, error) {
	var req *http.Request

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("HTTPRequest new request,err: %w", err)
	}

	resp, err := b.HTTPClientDynamic(ctx).Do(req)
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
