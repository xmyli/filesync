package track

import (
	"fmt"
	"log"
	"os"

	"github.com/rjeczalik/notify"
)

func (t *Tracker) watch() {
	for {
		c := make(chan notify.EventInfo, 1)

		watchPath := fmt.Sprint(t.Path, "/...")

		if err := notify.Watch(watchPath, c, notify.All); err != nil {
			log.Fatal(err)
		}
		defer notify.Stop(c)

		for {
			ei := <-c

			t.update()

			info, err := os.Stat(ei.Path())
			if err == nil && info.IsDir() {
				notify.Stop(c)
				break
			}
		}
	}
}
