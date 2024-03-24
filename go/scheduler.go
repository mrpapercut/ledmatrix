package main

import (
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	ticker        *time.Ticker
	ticks         int64
	canvas        *Canvas
	config        *Config
	lastShownJob  string
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
}

func (s *Scheduler) Stop() {
	cleanupOnce.Do(func() {
		log.Println("\nCTRL+C signal received, cleaning up...")
		s.canvas.Close()
		s.ticker.Stop()
	})
}

func (s *Scheduler) runJobs() {
	if !DuringWorkingHours() {
		return
	}

	jobs := Jobs{}

	// Check if a priority 1 message is waiting. If so, display it
	messageReady, messageToDraw := jobs.GetHighPriorityMessage()
	if messageReady && s.lastShownJob != messageToDraw.Type {
		jobs.DrawFeedMessage(messageToDraw)

		sqlite := getSQLiteInstance()
		sqlite.LowerPriorityAndDisplayTimes(messageToDraw)

		s.updateLastShown(messageToDraw.Type)
	} else {
		// Every 10 minutes
		if s.ticks%10 == 0 {
			messageReady, messageToDraw := jobs.GetLowPriorityMessage()
			if messageReady && s.lastShownJob != messageToDraw.Type {
				jobs.DrawFeedMessage(messageToDraw)

				sqlite := getSQLiteInstance()
				sqlite.LowerPriorityAndDisplayTimes(messageToDraw)

				s.updateLastShown(messageToDraw.Type)
			} else {
				jobs.DrawIdleScreen()

				s.updateLastShown("idle")
			}
		}
	}

	// Background jobs (don't use display)
	// Every 5 minutes
	if s.ticks%5 == 0 {
		youtube := getYoutubeInstance(s.config)
		youtube.GetVideos()
	}

	// Every 6 hours
	// if s.ticks%(6*60) == 0 {}

	// Reset ticks every 24 hours
	if s.ticks%(24*60) == 0 {
		s.ticks = 0
	} else {
		s.ticks = s.ticks + 1
	}
}

func (s *Scheduler) updateLastShown(jobName string) {
	s.lastShownJob = jobName
	s.lastShownTime = time.Now().Unix()
}
