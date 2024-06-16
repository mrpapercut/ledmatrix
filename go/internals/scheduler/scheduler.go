package scheduler

import (
	"log"
	"sync"
	"time"

	"github.com/mrpapercut/ledmatrix/internals/canvas"
	"github.com/mrpapercut/ledmatrix/internals/config"
	"github.com/mrpapercut/ledmatrix/internals/jobs"
	"github.com/mrpapercut/ledmatrix/internals/spritesheet"
	"github.com/mrpapercut/ledmatrix/internals/sqlite"
	"github.com/mrpapercut/ledmatrix/internals/types"
	"github.com/mrpapercut/ledmatrix/internals/utils"
	"github.com/mrpapercut/ledmatrix/internals/youtube"
)

type Scheduler struct {
	ticker             *time.Ticker
	ticks              int64
	config             *config.Config
	lastShownJob       string
	lastShownTime      int64
	lastShowJobChannel chan string
}

var (
	cleanupOnce sync.Once
)
var schedulerLock = &sync.Mutex{}
var schedulerInstance *Scheduler

func GetSchedulerInstance(config *config.Config) *Scheduler {
	if schedulerInstance == nil {
		schedulerLock.Lock()
		defer schedulerLock.Unlock()

		if schedulerInstance == nil {
			schedulerInstance = &Scheduler{
				ticker:             time.NewTicker(1 * time.Minute),
				ticks:              0,
				config:             config,
				lastShownJob:       "nothing",
				lastShownTime:      0,
				lastShowJobChannel: make(chan string),
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

	go func() {
		for lastShownJob := range s.lastShowJobChannel {
			s.updateLastShown(lastShownJob)
		}
	}()
}

func (s *Scheduler) Stop() {
	cleanupOnce.Do(func() {
		log.Println("\nCTRL+C signal received, cleaning up...")

		canvas := canvas.GetCanvasInstance()
		canvas.Close()

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

		s.drawMessage(messageToDraw)

		sqlite := sqlite.GetSQLiteInstance()
		sqlite.LowerPriorityAndDisplayTimes(messageToDraw)

		s.lastShowJobChannel <- messageToDraw.Type
	} else if s.ticks > 0 && s.ticks%10 == 0 { // Every 10 minutes
		messageReady, messageToDraw := jobs.GetLowPriorityMessage()
		if messageReady && s.lastShownJob != messageToDraw.Type {
			log.Printf("Drawing low priority message of type %s\n", messageToDraw.Type)

			s.drawMessage(messageToDraw)

			sqlite := sqlite.GetSQLiteInstance()
			sqlite.LowerPriorityAndDisplayTimes(messageToDraw)

			s.lastShowJobChannel <- messageToDraw.Type
		} else {
			log.Println("Drawing idle screen")

			jobs.DrawIdleScreen()

			s.lastShowJobChannel <- "idle"
		}
	} else {
		s.lastShowJobChannel <- "nothing" // great hack
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

func (s *Scheduler) drawMessage(message types.FeedMessage) {
	jobs := jobs.Jobs{}

	if message.Type == "youtubeVideo" {
		ytlogo, _ := spritesheet.GetSpritesheetFromJson("./sprites/youtubeLogo.json")
		jobs.DrawFeedMessage(message, ytlogo)
	} else {
		jobs.DrawFeedMessage(message)
	}
}

func (s *Scheduler) updateLastShown(jobName string) {
	s.lastShownJob = jobName
	s.lastShownTime = time.Now().Unix()
}
