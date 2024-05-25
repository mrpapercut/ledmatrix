package scheduler

import (
	"log"
	"sync"
	"time"

	"github.com/mrpapercut/ledmatrix/internals/canvas"
	"github.com/mrpapercut/ledmatrix/internals/config"
	"github.com/mrpapercut/ledmatrix/internals/jobs"
	"github.com/mrpapercut/ledmatrix/internals/sqlite"
	"github.com/mrpapercut/ledmatrix/internals/utils"
	"github.com/mrpapercut/ledmatrix/internals/youtube"
)

type Scheduler struct {
	ticker        *time.Ticker
	ticks         int64
	canvas        *canvas.Canvas
	config        *config.Config
	lastShownJob  string
	lastShownTime int64
}

var (
	cleanupOnce sync.Once
)
var schedulerLock = &sync.Mutex{}
var schedulerInstance *Scheduler

func GetSchedulerInstance(canvas *canvas.Canvas, config *config.Config) *Scheduler {
	if schedulerInstance == nil {
		schedulerLock.Lock()
		defer schedulerLock.Unlock()

		if schedulerInstance == nil {
			schedulerInstance = &Scheduler{
				ticker:        time.NewTicker(1 * time.Minute),
				ticks:         0,
				canvas:        canvas,
				config:        config,
				lastShownJob:  "",
				lastShownTime: 0,
			}
		}
	}

	return schedulerInstance
}

func (s *Scheduler) Start() {
	log.Println("Starting scheduler at", time.Now())

	go func() {
		for range s.ticker.C {
			s.runJobs()
		}
	}()

	jobs := jobs.Jobs{}
	jobs.DrawAnimationByFilename("kirbyDance01")
	jobs.DrawClock()
}

func (s *Scheduler) Stop() {
	cleanupOnce.Do(func() {
		log.Println("\nCTRL+C signal received, cleaning up...")
		s.canvas.Close()
		s.ticker.Stop()
	})
}

func (s *Scheduler) runJobs() {
	if !utils.DuringWorkingHours() {
		return
	}

	jobs := jobs.Jobs{}

	// Check if a priority 1 message is waiting. If so, display it
	messageReady, messageToDraw := jobs.GetHighPriorityMessage()
	if messageReady && s.lastShownJob != messageToDraw.Type {
		log.Printf("Drawing high priority message of type %s\n", messageToDraw.Type)

		jobs.DrawFeedMessage(messageToDraw)

		sqlite := sqlite.GetSQLiteInstance()
		sqlite.LowerPriorityAndDisplayTimes(messageToDraw)

		s.updateLastShown(messageToDraw.Type)
	} else {
		// Every 10 minutes
		if s.ticks > 0 && s.ticks%10 == 0 {
			messageReady, messageToDraw := jobs.GetLowPriorityMessage()
			if messageReady && s.lastShownJob != messageToDraw.Type {
				log.Printf("Drawing low priority message of type %s\n", messageToDraw.Type)

				jobs.DrawFeedMessage(messageToDraw)

				sqlite := sqlite.GetSQLiteInstance()
				sqlite.LowerPriorityAndDisplayTimes(messageToDraw)

				s.updateLastShown(messageToDraw.Type)
			} else {
				log.Println("Drawing idle screen")

				jobs.DrawIdleScreen()

				s.updateLastShown("idle")
			}
		}

		s.updateLastShown("nothing") // great hack
	}

	// Background jobs (don't use display)
	// Every 5 minutes
	if s.ticks > 0 && s.ticks%5 == 0 {
		log.Println("Getting Youtube videos")

		youtube := youtube.GetYoutubeInstance(s.config)
		go youtube.GetVideos()
	}

	// Every 6 hours
	// if s.ticks%(6*60) == 0 {}

	// Reset ticks every 24 hours
	if s.ticks > 0 && s.ticks%(24*60) == 0 {
		s.ticks = 0
	} else {
		s.ticks = s.ticks + 1
	}
}

func (s *Scheduler) updateLastShown(jobName string) {
	s.lastShownJob = jobName
	s.lastShownTime = time.Now().Unix()
}
