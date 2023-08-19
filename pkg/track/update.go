package track

import "log"

func (t *Tracker) update() {
	log.Println("UPDATE")
	t.scanTimer.Reset(t.waitDuration)
	t.onUpdate(t.currentState)
}
