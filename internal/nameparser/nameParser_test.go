package nameparser

import (
	"context"
	"fmt"
	"movieSpider/internal/config"
	"testing"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
}
func TestNameParserModelHandler(t *testing.T) {
	a := []string{"Australian Survivor S13",
		"Din.Don.9.Paesani.spaesati.2025.1080.mkv",
		"E1-4 of 8 [2025, WEBRip-AVC]",
		"Freakier Friday (2025) [1080p] [WEBRip] [x265] [10bit] [5.1]",
		"Freakier Friday (2025) [2160p] [WEBRip] [x265] [10bit] [5.1]",
		"Gabby's Dollhouse: The Movie (2025) [1080p] [WEBRip] [5.1] New",
		"Haul Out the Halloween (2025) [1080p] [WEBRip]",
		"Him (2025) [2160p] [WEBRip] [x265] [10bit] [5.1]",
		"[I.S.S.] Kamen Rider x Kamen Rider - Drive & Gaim Movie Taisen Full Throttle SD",
		"La donna della cabina numero 10 - The Woman in Cabin 10 (2025) Up... New",
		"Maintenance Required (2025) [2160p] [WEBRip] [x265] [10bit] [5.1] New",
		"Pursued (2025) [1080p] [WEBRip]",
		"Slow / Tu man nieko neprimeni [2023, WEB-DL 1080p] DVO (datynet) + Sub Rus, Eng + Original Lit",
		"The Derbyshire Auction House Season 2",
		"The Heatwave Lasted Four Days (1975) [EXTENDED.CUT.BLURAY] [1080p] [BluRay] [YTS.MX]",
		"The Partisan (2024) [1080p] [WEBRip]",
		"The Woman in Cabin 10 (2025) [1080p] [WEBRip] [x265] [10bit] [5.1]",
		"(COXC-1115) Akira Miyagawa Presents Space Battleship Yamato 2199 Concert 2015 | 宮川彬良 Presents『宇宙戦艦ヤマト2199』コンサート2015",
	}
	resulut, err := ModelHandler(context.Background(), a...)
	if err != nil {
		t.Error(err)
		return
	}
	for key, result := range resulut {
		fmt.Printf("%#v    ,%#v,%#v,%#v,%#v\n", key, result.TypeStr, result.NewName, result.Year, result.Resolution)
	}
}
