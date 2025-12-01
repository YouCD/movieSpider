package job

import (
	"context"
	"encoding/json"
	"github.com/robfig/cron/v3"
	"github.com/youcd/toolkit/log"
	"movieSpider/internal/bus"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"os"
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
func (r *ReleaseTimeJob) Run() {
	if r.scheduling == "" {
		log.WithCtx(context.Background()).Error("ReleaseTimeJob: Scheduling is null")
		os.Exit(1)
	}
	log.WithCtx(context.Background()).Infof("ReleaseTimeJob: Scheduling is: [%s]", r.scheduling)
	c := cron.New()
	_, _ = c.AddFunc(r.scheduling, func() {
		log.WithCtx(context.Background()).Infof("ReleaseTimeJob: Check video date for published.", r.scheduling)

		videos, err := model.NewMovieDB().FetchThisYearVideo()
		if err != nil {
			log.WithCtx(context.Background()).Error(err)
		}
		for _, video := range videos {
			if video.IsDatePublished() {
				go func(v *types.DouBanVideo) {
					bus.DatePublishedChan <- v
				}(video)
				var names []string
				err := json.Unmarshal([]byte(video.Names), &names)
				if err != nil {
					log.WithCtx(context.Background()).Error(err)
				}
				log.WithCtx(context.Background()).Infof("Video: %s , DatePublished: %v", names[0], video.DatePublished)
			}
		}
	})
	c.Start()
}
