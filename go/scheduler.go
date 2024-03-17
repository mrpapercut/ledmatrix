package main

import (
	"fmt"
	"sync"
	"time"
)

type Scheduler struct {
	thirtySecondsTicker *time.Ticker
	minuteTicker *time.Ticker
	fiveMinuteTicker *time.Ticker
	oneHourTicker *time.Ticker
	sixHoursTicker *time.Ticker
	canvas *Canvas
	config *Config
	lastShownJob string
	lastShownTime int64
}

var (
	cleanupOnce sync.Once
)
var schedulerLock = &sync.Mutex{}
var schedulerInstance *Scheduler

func getSchedulerInstance(canvas *Canvas, config *Config) *Scheduler {
	if schedulerInstance == nil {
		schedulerLock.Lock()
		defer schedulerLock.Unlock()

		if schedulerInstance == nil {
			schedulerInstance = &Scheduler{
				thirtySecondsTicker: time.NewTicker(30 * time.Second),
				minuteTicker: time.NewTicker(1 * time.Minute),
				fiveMinuteTicker: time.NewTicker(5 * time.Minute),
				oneHourTicker: time.NewTicker(1 * time.Hour),
				sixHoursTicker: time.NewTicker(6 * time.Hour),
				canvas: canvas,
				config: config,
				lastShownJob: "",
				lastShownTime: 0,
			}
		}
	}

	return schedulerInstance
}

func (s *Scheduler) Start() {
	fmt.Println("Starting scheduler at", time.Now())

	go func() {
		for {
			select{
			case <- s.thirtySecondsTicker.C:
				s.run30SecondsJobs()
			case <- s.minuteTicker.C:
				s.run1MinuteJobs()
			case <- s.fiveMinuteTicker.C:
				s.run5MinuteJobs()
			case <- s.oneHourTicker.C:
				s.run1HourJobs()
			case <- s.sixHoursTicker.C:
				s.run6HoursJobs()
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	cleanupOnce.Do(func() {
		fmt.Println("\nCTRL+C signal received, cleaning up...")
		s.canvas.Close()
		s.thirtySecondsTicker.Stop()
		s.minuteTicker.Stop()
		s.fiveMinuteTicker.Stop()
		s.oneHourTicker.Stop()
		s.sixHoursTicker.Stop()
	})
}

func (s *Scheduler) run30SecondsJobs() {
	if !DuringWorkingHours() {
		fmt.Println("Outside working hours!")
		return
	}

	// Check if a priority 1 message is waiting. If so, display it
	jobs := Jobs{}
	messageReady, messageToDraw := jobs.GetHighPriorityMessage()

	if !messageReady || s.lastShownJob == messageToDraw.Type {
		return
	}

	jobs.DrawFeedMessage(messageToDraw)
	s.lastShownJob = messageToDraw.Type
	s.lastShownTime = time.Now().Unix()

	sqlite := getSQLiteInstance()
	sqlite.LowerPriorityAndDisplayTimes(messageToDraw)
}

func (s *Scheduler) run1MinuteJobs() {
	if !DuringWorkingHours() {
		return
	}

	youtube := getYoutubeInstance(s.config)
	youtube.GetVideos()
}

func (s *Scheduler) run5MinuteJobs() {
	if !DuringWorkingHours() {
		return
	}

	now := time.Now().Unix()
	if (now - s.lastShownTime) > 30 {
		jobs := Jobs{}

		messageReady, messageToDraw := jobs.GetLowPriorityMessage()

		if messageReady && s.lastShownJob != messageToDraw.Type {
			jobs.DrawFeedMessage(messageToDraw)
			s.lastShownJob = messageToDraw.Type
			s.lastShownTime = time.Now().Unix()

			sqlite := getSQLiteInstance()
			sqlite.LowerPriorityAndDisplayTimes(messageToDraw)
		} else {
			go jobs.DrawIdleScreen()

			s.lastShownJob = "idle"
			s.lastShownTime = 0
		}
	}
}

func (s *Scheduler) run1HourJobs() {
	// fmt.Println("Running 1 hour jobs")

}

func (s *Scheduler) run6HoursJobs() {
	// fmt.Println("Running 6 hours jobs")
	// Fetch weather data
}
