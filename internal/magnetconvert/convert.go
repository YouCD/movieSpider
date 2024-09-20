package magnetconvert

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"errors"

	"github.com/anacrolix/torrent/metainfo"
)

var (
	ErrRespIsNil  = errors.New("response is nil")
	ErrMagnetMeta = errors.New("读取磁链meta信息错误")
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
		return "", ErrMagnetMeta
	}
	m2, err := mi.MagnetV2()
	if err != nil {
		return "", fmt.Errorf("转换失败，err: %w", err)
	}
	return m2.String(), nil
}

func FetchMagnetWithHTTPClient(url string, httpClient *http.Client) (magnet string, err error) {
	request, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("磁链获取错误,err: %w", err)
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("磁链获取错误,err: %w", err)
	}
	if resp == nil {
		return "", ErrRespIsNil
	}
	defer resp.Body.Close()

	magnet, err = IO2Magnet(resp.Body)
	if err != nil {
		return "", fmt.Errorf("磁链转换错误,err: %w", err)
	}
	return magnet, nil
}
