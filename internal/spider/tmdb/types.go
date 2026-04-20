package tmdb

type CrewItem struct {
	Job                string  `json:"job"`
	Department         string  `json:"department"`
	CreditId           string  `json:"credit_id"`
	Adult              bool    `json:"adult"`
	Gender             int     `json:"gender"`
	Id                 int     `json:"id"`
	KnownForDepartment string  `json:"known_for_department"`
	Name               string  `json:"name"`
	OriginalName       string  `json:"original_name"`
	Popularity         float64 `json:"popularity"`
	ProfilePath        *string `json:"profile_path"`
}
type GuestStar struct {
	Character          string  `json:"character"`
	CreditId           string  `json:"credit_id"`
	Order              int     `json:"order"`
	Adult              bool    `json:"adult"`
	Gender             int     `json:"gender"`
	Id                 int     `json:"id"`
	KnownForDepartment string  `json:"known_for_department"`
	Name               string  `json:"name"`
	OriginalName       string  `json:"original_name"`
	Popularity         float64 `json:"popularity"`
	ProfilePath        *string `json:"profile_path"`
}
type Episode struct {
	AirDate        string      `json:"air_date"`
	EpisodeNumber  int         `json:"episode_number"`
	EpisodeType    string      `json:"episode_type"`
	Id             int         `json:"id"`
	Name           string      `json:"name"`
	Overview       string      `json:"overview"`
	ProductionCode string      `json:"production_code"`
	Runtime        int         `json:"runtime"`
	SeasonNumber   int         `json:"season_number"`
	ShowId         int         `json:"show_id"`
	StillPath      string      `json:"still_path"`
	VoteAverage    float64     `json:"vote_average"`
	VoteCount      int         `json:"vote_count"`
	Crew           []CrewItem  `json:"crew"`
	GuestStars     []GuestStar `json:"guest_stars"`
}
type SeasonDetailResult struct {
	Id       string    `json:"_id"`
	AirDate  string    `json:"air_date"`
	Episodes []Episode `json:"episodes"`
	Name     string    `json:"name"`
	Networks []struct {
		Id            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	} `json:"networks"`
	Overview     string  `json:"overview"`
	Id1          int     `json:"id"`
	PosterPath   string  `json:"poster_path"`
	SeasonNumber int     `json:"season_number"`
	VoteAverage  float64 `json:"vote_average"`
}

// TVExternalIDs 电视剧外部ID
type TVExternalIDs struct {
	ImdbID      string `json:"imdb_id"`
	FreebaseMID string `json:"freebase_mid"`
	FreebaseID  string `json:"freebase_id"`
	TVDBID      int64  `json:"tvdb_id"`
	TVRageID    int64  `json:"tvrage_id"`
	WikidataID  string `json:"wikidata_id"`
	FacebookID  string `json:"facebook_id"`
	InstagramID string `json:"instagram_id"`
	TwitterID   string `json:"twitter_id"`
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
