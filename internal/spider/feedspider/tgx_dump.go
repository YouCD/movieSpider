package feedspider

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"
	"movieSpider/internal/magnetconvert"
	"movieSpider/internal/types"
	"strings"
	"sync"

	"github.com/youcd/toolkit/log"
)

type TgxDump struct {
	BaseFeeder
}

func NewTgxDump(scheduling, siteURL string, useIPProxy bool) *TgxDump {
	return &TgxDump{
		BaseFeeder{
			web:      "TgxDump",
			BaseFeed: types.BaseFeed{Url: siteURL, Scheduling: scheduling, UseIPProxy: useIPProxy},
		},
	}
}

func (t *TgxDump) Crawler() (videos []*types.FeedVideoBase, err error) {
	resp, err := t.HTTPRequest(t.Url)
	if err != nil {
		return nil, fmt.Errorf("TgxDump resp, err: %w", err)
	}

	reader, err := gzip.NewReader(bytes.NewReader(resp))
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
			var typ string
			typFiled := strings.ToLower(split[2])
			switch {
			case strings.Contains(typFiled, "tv"):
				typ = "tv"
			case strings.Contains(typFiled, "movie"):
				typ = "movie"
			}

			if strings.Contains(strings.ToLower(split[2]), "tv") || strings.Contains(strings.ToLower(split[2]), "movie") {
				video := &types.FeedVideoBase{
					TorrentName: split[1],
					TorrentURL:  split[4],
					Type:        typ,
					RowData:     sql.NullString{String: line},
					Web:         t.web,
				}

				if magnet, err := t.parser2Video(split[4]); err == nil {
					video.Magnet = magnet
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

func (t *TgxDump) parser2Video(torrentURL string) (string, error) {
	magnet, err := magnetconvert.FetchMagnetWithHTTPClient(torrentURL, t.HTTPClientDynamic())
	if err != nil {
		return "", fmt.Errorf("parser2Video magnet convert err: %w", err)
	}
	if magnet == "" {
		return "", ErrMagnetIsEmpty
	}
	return magnet, nil
}
