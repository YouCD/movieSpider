package magnetconvert

import (
	"context"
	"fmt"
	"io"
	httpClient2 "movieSpider/internal/httpclient"
	"net/http"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/pkg/errors"
)

// FileToMagnet
//
//	@Description: 通过文件获取磁链
//	@param file
//	@return string
//	@return error
func FileToMagnet(file string) (string, error) {
	mi, err := metainfo.LoadFromFile(file)
	if err != nil {
		//nolint:goerr113
		return "", fmt.Errorf("cannot read the metainfo from file: %s. %s", file, err.Error())
	}
	m2, err := mi.MagnetV2()
	if err != nil {
		//nolint:goerr113
		return "", fmt.Errorf("转换失败，err: %w", err)
	}
	return m2.String(), nil
}

// IO2Magnet
//
//	@Description: 通过io.Reader获取磁链
//	@param r
//	@return string
//	@return error
func IO2Magnet(r io.Reader) (string, error) {
	mi, err := metainfo.Load(r)
	if err != nil {
		return "", errors.New("读取磁链meta信息错误")
	}
	m2, err := mi.MagnetV2()
	if err != nil {
		return "", fmt.Errorf("转换失败，err: %w", err)
	}
	return m2.String(), nil
}

// FetchMagnet
//
//	@Description: 获取磁链
//	@param url
//	@return magnet
//	@return err
func FetchMagnet(url string) (magnet string, err error) {
	request, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		return "", errors.WithMessage(err, "TGx: 磁链获取错误")
	}
	client := httpClient2.NewHTTPClient()
	resp, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("磁链获取错误,err: %w", err)
	}
	if resp == nil {
		return "", errors.New("response is nil")
	}
	defer resp.Body.Close()

	magnet, err = IO2Magnet(resp.Body)
	if err != nil {
		return "", fmt.Errorf("磁链转换错误,err: %w", err)
	}
	return magnet, nil
}
