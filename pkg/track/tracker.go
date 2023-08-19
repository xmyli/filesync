package track

import (
	"path/filepath"
	"time"
)

type Tracker struct {
	Path           string
	updateInterval time.Duration // If no changes have been detected for this duration since the last scan, do an update.
	waitDuration   time.Duration // Do a scan if it has been this duration since the last update.
	ignoredNames   []string
	currentState   FSData
	nextState      FSData
	updateTimer    *time.Timer
	scanTimer      *time.Timer
	onUpdate       func(FSData)
	onScan         func(FSData)
}

func (t *Tracker) GetState() FSData {
	return t.currentState
}

func (t *Tracker) SetOnUpdate(onUpdate func(FSData)) {
	t.onUpdate = onUpdate
}

func (t *Tracker) SetOnScan(onScan func(FSData)) {
	t.onScan = onScan
}

func (t *Tracker) Start() {
	t.updateTimer = time.NewTimer(t.updateInterval)
	t.scanTimer = time.NewTimer(0)

	go t.watch()

	go func() {
		for {
			<-t.updateTimer.C
			t.update()
		}
	}()

	go func() {
		for {
			<-t.scanTimer.C
			t.scan()
		}
	}()

	select {}
}

func NewTracker(path string, updateInterval time.Duration, waitDuration time.Duration, ignoredNames []string) Tracker {
	tracker := Tracker{}
	tracker.Path = filepath.Clean(path)
	tracker.updateInterval = updateInterval
	tracker.waitDuration = waitDuration
	tracker.ignoredNames = ignoredNames
	tracker.currentState = FSData{
		Objects: map[string]FSObject{},
	}
	tracker.nextState = FSData{
		Objects: map[string]FSObject{},
	}
	tracker.onUpdate = func(_ FSData) {}
	tracker.onScan = func(_ FSData) {}
	return tracker
}
