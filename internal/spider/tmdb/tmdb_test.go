package tmdb

import (
	"context"
	"movieSpider/internal/config"
	"testing"
)

var (
	accountID int
	apiToken  string
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	accountID = config.Config.TMDB.AccountID
	apiToken = config.Config.TMDB.BearerToken
}

func TestGetWatchlistMovies(t *testing.T) {
	client, err := NewClient(accountID, apiToken)
	if err != nil {
		t.Fatalf("创建 TMDB 客户端失败: %v", err)
	}

	ctx := context.Background()
	result, err := client.GetWatchlistMovies(ctx, 1)
	if err != nil {
		t.Fatalf("获取 watchlist 电影失败: %v", err)
	}

	t.Logf("共 %d 条电影记录，第 %d/%d 页", result.TotalResults, result.Page, result.TotalPages)
	for _, movie := range result.Results {
		t.Logf("电影: %s (%s) - ID: %d", movie.Title, movie.ReleaseDate, movie.ID)
	}
}

func TestGetWatchlistTV(t *testing.T) {
	client, err := NewClient(accountID, apiToken)
	if err != nil {
		t.Fatalf("创建 TMDB 客户端失败: %v", err)
	}

	ctx := context.Background()
	result, err := client.GetWatchlistTV(ctx, 1)
	if err != nil {
		t.Fatalf("获取 watchlist 电视剧失败: %v", err)
	}

	t.Logf("共 %d 条电视剧记录，第 %d/%d 页", result.TotalResults, result.Page, result.TotalPages)
	for _, tv := range result.Results {
		t.Logf("电视剧: %s (%s) - ID: %d", tv.Name, tv.FirstAirDate, tv.ID)
	}
}

func TestGetMovieDetails(t *testing.T) {
	client, err := NewClient(accountID, apiToken) // 详情获取不需要 accountID 和 sessionID
	if err != nil {
		t.Fatalf("创建 TMDB 客户端失败: %v", err)
	}

	ctx := context.Background()
	// 使用一个真实的电影 ID 进行测试，例如：《阿凡达》ID: 19995
	movieID := 19995

	result, err := client.GetMovieDetails(ctx, movieID)
	if err != nil {
		t.Fatalf("获取电影详情失败: %v", err)
	}

	t.Logf("电影: %s", result.Title)
	t.Logf("原始标题: %s", result.OriginalTitle)
	t.Logf("简介: %s", result.Overview)
	t.Logf("上映日期: %s", result.ReleaseDate)
	t.Logf("时长: %d 分钟", result.Runtime)
	t.Logf("评分: %.2f (%d 票)", result.VoteAverage, result.VoteCount)
	t.Logf("IMDB ID: %s", result.ImdbID)
	t.Logf("状态: %s", result.Status)
}

func TestGetTVDetails(t *testing.T) {
	client, err := NewClient(accountID, apiToken) // 详情获取不需要 accountID 和 sessionID
	if err != nil {
		t.Fatalf("创建 TMDB 客户端失败: %v", err)
	}

	ctx := context.Background()
	// 使用一个真实的电视剧 ID 进行测试，例如：《权力的游戏》ID: 1399
	tvID := 111110

	result, err := client.GetTVDetails(ctx, tvID)
	if err != nil {
		t.Fatalf("获取电视剧详情失败: %v", err)
	}

	t.Logf("电视剧: %s", result.Name)
	t.Logf("原始名称: %s", result.OriginalName)
	t.Logf("简介: %s", result.Overview)
	t.Logf("首播日期: %s", result.FirstAirDate)
	t.Logf("季数: %d", result.NumberOfSeasons)
	t.Logf("集数: %d", result.NumberOfEpisodes)
	t.Logf("评分: %.2f (%d 票)", result.VoteAverage, result.VoteCount)
	t.Logf("状态: %s", result.Status)
	t.Logf("是否在制作中: %v", result.InProduction)

	for _, season := range result.Seasons {
		t.Logf("季 %d: %s (%d 集)", season.SeasonNumber, season.Name, season.EpisodeCount)
	}
}

func TestGetMovieWatchProviders(t *testing.T) {
	client, err := NewClient(accountID, apiToken)
	if err != nil {
		t.Fatalf("创建 TMDB 客户端失败: %v", err)
	}

	ctx := context.Background()
	// 使用一个真实的电影 ID 进行测试，例如：《阿凡达》ID: 19995
	movieID := 19995

	result, err := client.GetMovieWatchProviders(ctx, movieID)
	if err != nil {
		t.Fatalf("获取电影观看提供商失败: %v", err)
	}

	t.Logf("电影 ID: %d", result.ID)
	for country, providers := range result.Results {
		t.Logf("国家/地区: %s, 链接: %s", country, providers.Link)
		if providers.Flatrate != nil {
			t.Logf("  流媒体订阅:")
			for _, p := range *providers.Flatrate {
				t.Logf("    - %s (ID: %d)", p.ProviderName, p.ProviderID)
			}
		}
		if providers.Rent != nil {
			t.Logf("  租赁:")
			for _, p := range *providers.Rent {
				t.Logf("    - %s (ID: %d)", p.ProviderName, p.ProviderID)
			}
		}
		if providers.Buy != nil {
			t.Logf("  购买:")
			for _, p := range *providers.Buy {
				t.Logf("    - %s (ID: %d)", p.ProviderName, p.ProviderID)
			}
		}
	}
}

func TestGetTVWatchProviders(t *testing.T) {
	client, err := NewClient(0, apiToken)
	if err != nil {
		t.Fatalf("创建 TMDB 客户端失败: %v", err)
	}

	ctx := context.Background()
	// 使用一个真实的电视剧 ID 进行测试，例如：《权力的游戏》ID: 1399
	tvID := 1399

	result, err := client.GetTVWatchProviders(ctx, tvID)
	if err != nil {
		t.Fatalf("获取电视剧观看提供商失败: %v", err)
	}

	t.Logf("电视剧 ID: %d", result.ID)
	for country, providers := range result.Results {
		t.Logf("国家/地区: %s, 链接: %s", country, providers.Link)
		if providers.Flatrate != nil {
			t.Logf("  流媒体订阅:")
			for _, p := range *providers.Flatrate {
				t.Logf("    - %s (ID: %d)", p.ProviderName, p.ProviderID)
			}
		}
		if providers.Rent != nil {
			t.Logf("  租赁:")
			for _, p := range *providers.Rent {
				t.Logf("    - %s (ID: %d)", p.ProviderName, p.ProviderID)
			}
		}
		if providers.Buy != nil {
			t.Logf("  购买:")
			for _, p := range *providers.Buy {
				t.Logf("    - %s (ID: %d)", p.ProviderName, p.ProviderID)
			}
		}
	}
}

func TestClient_GetSearchMovies(t *testing.T) {
	client, err := NewClient(0, apiToken)
	if err != nil {
		t.Fatalf("创建 TMDB 客户端失败: %v", err)
	}

	ctx := context.Background()

	client.GetSearchMovies(ctx, "阿凡达")
}

func TestClient_GetSearchTVShow(t *testing.T) {
	client, err := NewClient(0, apiToken)
	if err != nil {
		t.Fatalf("创建 TMDB 客户端失败: %v", err)
	}

	ctx := context.Background()

	show, err := client.GetSearchTVShow(ctx, "One Piece")
	if err != nil {
		t.Fatalf("搜索电视剧失败: %v", err)
	}
	for _, result := range show {
		t.Logf("搜索结果: %+v", result)
	}
}

func TestClient_GetTVSeasonWatchProviders(t *testing.T) {
	client, err := NewClient(0, apiToken)
	if err != nil {
		t.Fatalf("创建 TMDB 客户端失败: %v", err)
	}

	ctx := context.Background()
	providers, err := client.GetTVSeasonWatchProviders(ctx, 111110, 1, nil)
	if err != nil {
		t.Fatalf("获取电视剧季的观看提供商失败: %v", err)
	}
	t.Logf("电视剧季的观看提供商: %+v", providers)
}

func TestClient_GetTVSeasonDetails(t *testing.T) {
	client, err := NewClient(0, apiToken)
	if err != nil {
		t.Fatalf("创建 TMDB 客户端失败: %v", err)
	}

	ctx := context.Background()
	season, err := client.GetTVSeasonDetails(ctx, 111110, 2)
	if err != nil {
		t.Fatalf("获取电视剧季详情失败: %v", err)
	}
	t.Logf("电视剧季详情: %+v", season)
}
