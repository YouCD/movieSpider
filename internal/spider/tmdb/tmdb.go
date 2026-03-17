package tmdb

import (
	"context"
	"fmt"
	"movieSpider/internal/httpclient"
	"movieSpider/internal/types"
	"strconv"

	"github.com/cyruzin/golang-tmdb"
)

// Client TMDB 客户端
type Client struct {
	client    *tmdb.Client
	accountID int
}

// NewClient 创建 TMDB 客户端
// accountID: TMDB 账户 ID
// apiToken: TMDB API Token (Bearer Token)
// sessionID: TMDB 会话 ID (可选，获取 watchlist 需要)
func NewClient(accountID int, apiToken string) (*Client, error) {
	client, err := tmdb.InitV4(apiToken)
	if err != nil {
		return nil, fmt.Errorf("初始化 TMDB 客户端失败: %w", err)
	}

	client.SetClientConfig(*httpclient.NewProxyHTTPClient(context.Background()))
	return &Client{
		client:    client,
		accountID: accountID,
	}, nil
}

// GetWatchlistMovies 获取账户的 watchlist 电影清单
// page: 页码，从 1 开始
func (c *Client) GetWatchlistMovies(ctx context.Context, page int) (*types.WatchlistMoviesResponse, error) {
	options := map[string]string{
		"page": strconv.Itoa(page),
	}

	result, err := c.client.GetMovieWatchlist(c.accountID, options)
	if err != nil {
		return nil, fmt.Errorf("获取 watchlist 电影失败: %w", err)
	}

	response := &types.WatchlistMoviesResponse{
		Page:         int(result.Page),
		TotalPages:   int(result.TotalPages),
		TotalResults: int(result.TotalResults),
		Results:      make([]types.WatchlistMovieItem, 0, len(result.Results)),
	}

	for _, movie := range result.Results {
		response.Results = append(response.Results, types.WatchlistMovieItem{
			Adult:            movie.Adult,
			BackdropPath:     movie.BackdropPath,
			GenreIDs:         movie.GenreIDs,
			ID:               int(movie.ID),
			OriginalLanguage: movie.OriginalLanguage,
			OriginalTitle:    movie.OriginalTitle,
			Overview:         movie.Overview,
			Popularity:       movie.Popularity,
			PosterPath:       movie.PosterPath,
			ReleaseDate:      movie.ReleaseDate,
			Title:            movie.Title,
			Video:            movie.Video,
			VoteAverage:      float64(movie.VoteAverage),
			VoteCount:        int(movie.VoteCount),
		})
	}

	return response, nil
}

// GetWatchlistTV 获取账户的 watchlist 电视剧清单
// page: 页码，从 1 开始
func (c *Client) GetWatchlistTV(ctx context.Context, page int) (*types.WatchlistTVResponse, error) {
	options := map[string]string{
		"page": strconv.Itoa(page),
	}

	result, err := c.client.GetTVShowsWatchlist(c.accountID, options)
	if err != nil {
		return nil, fmt.Errorf("获取 watchlist 电视剧失败: %w", err)
	}

	response := &types.WatchlistTVResponse{
		Page:         int(result.Page),
		TotalPages:   int(result.TotalPages),
		TotalResults: int(result.TotalResults),
		Results:      make([]types.WatchlistTVItem, 0, len(result.Results)),
	}

	for _, tv := range result.Results {
		genreIDs := make([]int, 0, len(tv.GenreIDs))
		for _, id := range tv.GenreIDs {
			genreIDs = append(genreIDs, int(id))
		}

		response.Results = append(response.Results, types.WatchlistTVItem{
			BackdropPath:     tv.BackdropPath,
			GenreIDs:         genreIDs,
			ID:               int(tv.ID),
			OriginCountry:    tv.OriginCountry,
			OriginalLanguage: tv.OriginalLanguage,
			OriginalName:     tv.OriginalName,
			Overview:         tv.Overview,
			Popularity:       tv.Popularity,
			PosterPath:       tv.PosterPath,
			FirstAirDate:     tv.FirstAirDate,
			Name:             tv.Name,
			VoteAverage:      float64(tv.VoteAverage),
			VoteCount:        int(tv.VoteCount),
		})
	}

	return response, nil
}

// GetMovieDetails 获取电影详情
// movieID: 电影 ID
func (c *Client) GetMovieDetails(ctx context.Context, movieID int) (*types.TmDBMovieDetailData, error) {
	result, err := c.client.GetMovieDetails(movieID, nil)
	if err != nil {
		return nil, fmt.Errorf("获取电影详情失败: %w", err)
	}

	genres := make([]struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}, 0, len(result.Genres))
	for _, g := range result.Genres {
		genres = append(genres, struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{
			ID:   int(g.ID),
			Name: g.Name,
		})
	}

	productionCountries := make([]struct {
		Iso31661 string `json:"iso_3166_1"`
		Name     string `json:"name"`
	}, 0, len(result.ProductionCountries))
	for _, pc := range result.ProductionCountries {
		productionCountries = append(productionCountries, struct {
			Iso31661 string `json:"iso_3166_1"`
			Name     string `json:"name"`
		}{
			Iso31661: pc.Iso3166_1,
			Name:     pc.Name,
		})
	}

	spokenLanguages := make([]struct {
		EnglishName string `json:"english_name"`
		Iso6391     string `json:"iso_639_1"`
		Name        string `json:"name"`
	}, 0, len(result.SpokenLanguages))
	for _, sl := range result.SpokenLanguages {
		spokenLanguages = append(spokenLanguages, struct {
			EnglishName string `json:"english_name"`
			Iso6391     string `json:"iso_639_1"`
			Name        string `json:"name"`
		}{
			EnglishName: "", // SDK 没有 EnglishName 字段
			Iso6391:     sl.Iso639_1,
			Name:        sl.Name,
		})
	}

	return &types.TmDBMovieDetailData{
		Adult:               result.Adult,
		BackdropPath:        result.BackdropPath,
		BelongsToCollection: result.BelongsToCollection,
		Budget:              int(result.Budget),
		Genres:              genres,
		Homepage:            result.Homepage,
		ID:                  int(result.ID),
		ImdbID:              result.IMDbID,
		OriginalLanguage:    result.OriginalLanguage,
		OriginalTitle:       result.OriginalTitle,
		Overview:            result.Overview,
		Popularity:          float64(result.Popularity),
		PosterPath:          result.PosterPath,
		ProductionCompanies: nil,
		ProductionCountries: productionCountries,
		ReleaseDate:         result.ReleaseDate,
		Revenue:             int(result.Revenue),
		Runtime:             result.Runtime,
		SpokenLanguages:     spokenLanguages,
		Status:              result.Status,
		Tagline:             result.Tagline,
		Title:               result.Title,
		Video:               result.Video,
		VoteAverage:         float64(result.VoteAverage),
		VoteCount:           int(result.VoteCount),
	}, nil
}

// GetTVDetails 获取电视剧详情
// tvID: 电视剧 ID
func (c *Client) GetTVDetails(ctx context.Context, tvID int) (*types.TmDBTVDetailData, error) {
	result, err := c.client.GetTVDetails(tvID, nil)
	if err != nil {
		return nil, fmt.Errorf("获取电视剧详情失败: %w", err)
	}

	return convertTVDetails(result), nil
}

func convertTVDetails(result *tmdb.TVDetails) *types.TmDBTVDetailData {
	if result == nil {
		return nil
	}

	createdBy := make([]struct {
		ID          int    `json:"id"`
		CreditID    string `json:"credit_id"`
		Name        string `json:"name"`
		Gender      int    `json:"gender"`
		ProfilePath string `json:"profile_path"`
	}, 0, len(result.CreatedBy))
	for _, cb := range result.CreatedBy {
		createdBy = append(createdBy, struct {
			ID          int    `json:"id"`
			CreditID    string `json:"credit_id"`
			Name        string `json:"name"`
			Gender      int    `json:"gender"`
			ProfilePath string `json:"profile_path"`
		}{
			ID:          int(cb.ID),
			CreditID:    cb.CreditID,
			Name:        cb.Name,
			Gender:      cb.Gender,
			ProfilePath: cb.ProfilePath,
		})
	}

	genres := make([]struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}, 0, len(result.Genres))
	for _, g := range result.Genres {
		genres = append(genres, struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{
			ID:   int(g.ID),
			Name: g.Name,
		})
	}

	networks := make([]struct {
		ID            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	}, 0, len(result.Networks))
	for _, n := range result.Networks {
		networks = append(networks, struct {
			ID            int    `json:"id"`
			LogoPath      string `json:"logo_path"`
			Name          string `json:"name"`
			OriginCountry string `json:"origin_country"`
		}{
			ID:            int(n.ID),
			LogoPath:      n.LogoPath,
			Name:          n.Name,
			OriginCountry: n.OriginCountry,
		})
	}

	productionCompanies := make([]struct {
		ID            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	}, 0, len(result.ProductionCompanies))
	for _, pc := range result.ProductionCompanies {
		productionCompanies = append(productionCompanies, struct {
			ID            int    `json:"id"`
			LogoPath      string `json:"logo_path"`
			Name          string `json:"name"`
			OriginCountry string `json:"origin_country"`
		}{
			ID:            int(pc.ID),
			LogoPath:      pc.LogoPath,
			Name:          pc.Name,
			OriginCountry: pc.OriginCountry,
		})
	}

	productionCountries := make([]struct {
		Iso31661 string `json:"iso_3166_1"`
		Name     string `json:"name"`
	}, 0, len(result.ProductionCountries))
	for _, pc := range result.ProductionCountries {
		productionCountries = append(productionCountries, struct {
			Iso31661 string `json:"iso_3166_1"`
			Name     string `json:"name"`
		}{
			Iso31661: pc.Iso3166_1,
			Name:     pc.Name,
		})
	}

	seasons := make([]struct {
		AirDate      string  `json:"air_date"`
		EpisodeCount int     `json:"episode_count"`
		ID           int     `json:"id"`
		Name         string  `json:"name"`
		Overview     string  `json:"overview"`
		PosterPath   string  `json:"poster_path"`
		SeasonNumber int     `json:"season_number"`
		VoteAverage  float64 `json:"vote_average"`
	}, 0, len(result.Seasons))
	for _, s := range result.Seasons {
		seasons = append(seasons, struct {
			AirDate      string  `json:"air_date"`
			EpisodeCount int     `json:"episode_count"`
			ID           int     `json:"id"`
			Name         string  `json:"name"`
			Overview     string  `json:"overview"`
			PosterPath   string  `json:"poster_path"`
			SeasonNumber int     `json:"season_number"`
			VoteAverage  float64 `json:"vote_average"`
		}{
			AirDate:      s.AirDate,
			EpisodeCount: s.EpisodeCount,
			ID:           int(s.ID),
			Name:         s.Name,
			Overview:     s.Overview,
			PosterPath:   s.PosterPath,
			SeasonNumber: s.SeasonNumber,
			VoteAverage:  float64(s.VoteAverage),
		})
	}

	// TVDetails 没有 SpokenLanguages 字段，使用空切片
	spokenLanguages := make([]struct {
		EnglishName string `json:"english_name"`
		Iso6391     string `json:"iso_639_1"`
		Name        string `json:"name"`
	}, 0)

	return &types.TmDBTVDetailData{
		Adult:               false, // TVDetails 没有 Adult 字段
		BackdropPath:        result.BackdropPath,
		CreatedBy:           createdBy,
		EpisodeRunTime:      convertEpisodeRunTime(result.EpisodeRunTime),
		FirstAirDate:        result.FirstAirDate,
		Genres:              genres,
		Homepage:            result.Homepage,
		ID:                  int(result.ID),
		InProduction:        result.InProduction,
		Languages:           result.Languages,
		LastAirDate:         result.LastAirDate,
		Name:                result.Name,
		Networks:            networks,
		NumberOfEpisodes:    result.NumberOfEpisodes,
		NumberOfSeasons:     result.NumberOfSeasons,
		OriginCountry:       result.OriginCountry,
		OriginalLanguage:    result.OriginalLanguage,
		OriginalName:        result.OriginalName,
		Overview:            result.Overview,
		Popularity:          float64(result.Popularity),
		PosterPath:          result.PosterPath,
		ProductionCompanies: productionCompanies,
		ProductionCountries: productionCountries,
		Seasons:             seasons,
		SpokenLanguages:     spokenLanguages,
		Status:              result.Status,
		Tagline:             result.Tagline,
		Type:                result.Type,
		VoteAverage:         float64(result.VoteAverage),
		VoteCount:           int(result.VoteCount),
	}
}

func convertEpisodeRunTime(runtimes []int) []interface{} {
	result := make([]interface{}, len(runtimes))
	for i, r := range runtimes {
		result[i] = r
	}
	return result
}

// WatchProvider 观看提供商
type WatchProvider struct {
	LogoPath        string `json:"logo_path"`
	ProviderID      int    `json:"provider_id"`
	ProviderName    string `json:"provider_name"`
	DisplayPriority int    `json:"display_priority"`
}

// WatchProviderResult 国家/地区的观看提供商结果
type WatchProviderResult struct {
	Link     string           `json:"link"`
	Flatrate *[]WatchProvider `json:"flatrate"` // 流媒体订阅
	Rent     *[]WatchProvider `json:"rent"`     // 租赁
	Buy      *[]WatchProvider `json:"buy"`      // 购买
}

// WatchProvidersResponse 观看提供商响应
type WatchProvidersResponse struct {
	ID      int                            `json:"id"`
	Results map[string]WatchProviderResult `json:"results"` // key 是国家代码 (如 "CN", "US")
}

// GetMovieWatchProviders 获取电影的观看提供商
// movieID: 电影 ID
// 返回按国家/地区分组的观看提供商信息（流媒体、租赁、购买等）
func (c *Client) GetMovieWatchProviders(ctx context.Context, movieID int) (*WatchProvidersResponse, error) {
	result, err := c.client.GetMovieWatchProviders(movieID, nil)
	if err != nil {
		return nil, fmt.Errorf("获取电影观看提供商失败: %w", err)
	}

	response := &WatchProvidersResponse{
		ID:      int(result.ID),
		Results: make(map[string]WatchProviderResult),
	}

	for country, provider := range result.Results {
		response.Results[country] = convertWatchProviderResult(provider)
	}

	return response, nil
}

// GetTVWatchProviders 获取电视剧的观看提供商
// tvID: 电视剧 ID
// 返回按国家/地区分组的观看提供商信息（流媒体、租赁、购买等）
func (c *Client) GetTVWatchProviders(ctx context.Context, tvID int) (*WatchProvidersResponse, error) {
	result, err := c.client.GetTVWatchProviders(tvID, nil)
	if err != nil {
		return nil, fmt.Errorf("获取电视剧观看提供商失败: %w", err)
	}

	response := &WatchProvidersResponse{
		ID:      int(result.ID),
		Results: make(map[string]WatchProviderResult),
	}

	for country, provider := range result.Results {
		response.Results[country] = convertWatchProviderResult(provider)
	}

	return response, nil
}

func convertWatchProviderResult(provider tmdb.WatchProviderResult) WatchProviderResult {
	result := WatchProviderResult{
		Link: provider.Link,
	}

	if provider.Flatrate != nil {
		flatrate := make([]WatchProvider, len(*provider.Flatrate))
		for i, p := range *provider.Flatrate {
			flatrate[i] = WatchProvider{
				LogoPath:        p.LogoPath,
				ProviderID:      p.ProviderID,
				ProviderName:    p.ProviderName,
				DisplayPriority: p.DisplayPriority,
			}
		}
		result.Flatrate = &flatrate
	}

	if provider.Rent != nil {
		rent := make([]WatchProvider, len(*provider.Rent))
		for i, p := range *provider.Rent {
			rent[i] = WatchProvider{
				LogoPath:        p.LogoPath,
				ProviderID:      p.ProviderID,
				ProviderName:    p.ProviderName,
				DisplayPriority: p.DisplayPriority,
			}
		}
		result.Rent = &rent
	}

	if provider.Buy != nil {
		buy := make([]WatchProvider, len(*provider.Buy))
		for i, p := range *provider.Buy {
			buy[i] = WatchProvider{
				LogoPath:        p.LogoPath,
				ProviderID:      p.ProviderID,
				ProviderName:    p.ProviderName,
				DisplayPriority: p.DisplayPriority,
			}
		}
		result.Buy = &buy
	}

	return result
}
