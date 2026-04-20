package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// HasWatchProviders 检查是否有可用的观看提供商
func HasWatchProviders(providers *WatchProvidersResponse) bool {
	if providers == nil || len(providers.Results) == 0 {
		return false
	}

	// 遍历所有国家/地区的提供商
	for _, result := range providers.Results {
		// 检查是否有流媒体订阅
		if result.Flatrate != nil && len(*result.Flatrate) > 0 {
			return true
		}
		// 检查是否有租赁
		if result.Rent != nil && len(*result.Rent) > 0 {
			return true
		}
		// 检查是否有购买
		if result.Buy != nil && len(*result.Buy) > 0 {
			return true
		}
	}

	return false
}

func (c *Client) fmtOptions(urlOptions map[string]string) string {
	options := ""
	if len(urlOptions) > 0 {
		for key, value := range urlOptions {
			options += fmt.Sprintf(
				"&%s=%s",
				key,
				url.QueryEscape(value),
			)
		}
	}
	return options
}

func (c *Client) get(ctx context.Context, url string, data any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("could not fetch the url: %s", err)
	}
	req.Header.Add("content-type", "application/json;charset=utf-8")
	if c.bearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+c.bearerToken)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 status code: %s", res.Status)
	}
	if err = json.NewDecoder(res.Body).Decode(data); err != nil {
		return fmt.Errorf("could not decode the data: %s", err)
	}
	return nil
}
