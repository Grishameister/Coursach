package check

import (
	"log"
	"time"
)

var PassToQueue = true

func Check(in chan interface{}) {
	for {
		ticker := time.NewTicker(10 * time.Second)
		select {
		case <-in:
			PassToQueue = true
		case <-ticker.C:
			log.Println("tick")
			PassToQueue = false
		}
	}
}
