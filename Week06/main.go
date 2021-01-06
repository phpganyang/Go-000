package main

import (
	"Week06/roll"
	"log"
	"time"
)

var rd *roll.RollingDemo

func init() {
	rd = roll.NewBucket()
}

func main() {
	for i := 1; i <= 20; i++ {
		now := time.Now()
		time.Sleep(1 * time.Second)
		rd.Increment()
		log.Printf("avg: %d", rd.Avg(now))
	}

}
