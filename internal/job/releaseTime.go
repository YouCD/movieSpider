package job

import (
	"context"
	"encoding/json"
	"os"

	"movieSpider/internal/bus"
	"movieSpider/internal/model"
	"movieSpider/internal/types"

	"github.com/robfig/cron/v3"
	"github.com/youcd/toolkit/log"
)

type ReleaseTimeJob struct {
	scheduling string
}

func NewReleaseTimeJob(scheduling string) *ReleaseTimeJob {
	if scheduling == "" {
		return &ReleaseTimeJob{scheduling: "0 9 * * *"}
	}
	return &ReleaseTimeJob{scheduling: scheduling}
}

func (r *ReleaseTimeJob) Run(ctx context.Context) {
	if r.scheduling == "" {
		log.WithCtx(ctx).Error("ReleaseTimeJob: Scheduling is null")
		os.Exit(1)
	}
	log.WithCtx(ctx).Infof("ReleaseTimeJob: Scheduling is: [%s]", r.scheduling)
	c := cron.New()
	_, _ = c.AddFunc(r.scheduling, func() {
		log.WithCtx(ctx).Infof("ReleaseTimeJob: Check video date for published.", r.scheduling)

		videos, err := model.NewMovieDB().FetchThisYearVideo(ctx)
		if err != nil {
			log.WithCtx(ctx).Error(err)
		}
		for _, video := range videos {
			if video.IsDatePublished() {
				go func(v *types.TMDBVideo) {
					bus.DatePublishedChan <- v
				}(video)
				var names []string
				err := json.Unmarshal([]byte(video.Names), &names)
				if err != nil {
					log.WithCtx(ctx).Error(err)
				}
				log.WithCtx(ctx).Infof("Video: %s , DatePublished: %v", names[0], video.DatePublished)
			}
		}
	})
	c.Start()
}
