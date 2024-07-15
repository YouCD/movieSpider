package feedspider

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"database/sql"
	"fmt"
	"io"
	"movieSpider/internal/httpclient"
	"movieSpider/internal/magnetconvert"
	"movieSpider/internal/types"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/youcd/toolkit/log"
)

type TgxDump struct {
	BaseFeeder
}

func NewTgxDump(scheduling, siteURL string) *TgxDump {
	return &TgxDump{
		BaseFeeder{
			web:        "TgxDump",
			url:        siteURL,
			scheduling: scheduling,
		},
	}
}

func (t *TgxDump) Crawler() (videos []*types.FeedVideo, err error) {
	c := httpclient.NewHTTPClient()
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, t.url, nil)
	if err != nil {
		return nil, fmt.Errorf("TgxDump new request, err: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("TgxDump resp, err: %w", err)
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("TgxDump reader, err: %w", err)
	}

	var (
		buf bytes.Buffer
		wg  sync.WaitGroup
		mu  sync.Mutex
	)
	const maxSize = 10 * 1024 * 1024 // 10 MB
	limitedReader := io.LimitReader(reader, maxSize)

	if _, err := io.Copy(&buf, limitedReader); err != nil {
		log.Error(err)
		return nil, fmt.Errorf("TgxDump copy, err: %w", err)
	}

	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, "|")
		if len(split) < 4 {
			continue
		}
		wg.Add(1)
		go func(split []string) {
			defer wg.Done()
			video := new(types.FeedVideo) //nolint
			if strings.Contains(strings.ToLower(split[2]), "tv") || strings.Contains(strings.ToLower(split[2]), "movies") {
				if video = t.parser2Video(split[1], split[2], line, split[4]); video != nil {
					mu.Lock()
					videos = append(videos, video)
					mu.Unlock()
				}
			}
		}(split)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("TgxDump scanner err: %w", err)
	}
	wg.Wait()
	return videos, nil
}

func (t *TgxDump) parser2Video(name, typ, rowData, torrentURL string) *types.FeedVideo {
	torrentName := strings.ReplaceAll(name, " ", ".")
	compileRegex := regexp.MustCompile(`(.*)\.(\d{4})\.`)
	matchArr := compileRegex.FindStringSubmatch(torrentName)
	if len(matchArr) < 3 {
		return nil
	}
	year := matchArr[2]
	if len(matchArr) == 0 {
		tvReg := regexp.MustCompile(`(.*)(\.[Ss][0-9][0-9][eE][0-9][0-9])`)
		TVNameArr := tvReg.FindStringSubmatch(torrentName)
		// 如果 正则匹配过后 没有结果直接 过滤掉
		if len(TVNameArr) == 0 {
			return nil
		}
		name = TVNameArr[1]
	} else {
		name = matchArr[1]
	}
	magnet, err := magnetconvert.FetchMagnet(torrentURL)
	if err != nil {
		log.Error(err)
		return nil
	}
	if magnet == "" {
		log.Warnf("spider: %s , name :%s , magnet is empty", t.web, name)
		return nil
	}
	return &types.FeedVideo{
		Name:        name,
		TorrentName: torrentName,
		TorrentURL:  torrentURL,
		Magnet:      magnet,
		Year:        year,
		Type:        typ,
		RowData:     sql.NullString{String: rowData},
		Web:         t.web,
	}
}
