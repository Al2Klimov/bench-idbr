package main

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
	"time"
)

type event struct {
	state interface{}
	typs  []interface{}
}

func main() {
	host := flag.String("host", "localhost", "HOST")
	port := flag.Uint("port", 6380, "PORT")
	flag.Parse()

	client := redis.NewClient(&redis.Options{
		Addr:         net.JoinHostPort(*host, strconv.FormatUint(uint64(*port), 10)),
		DialTimeout:  time.Minute / 2,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		PoolTimeout:  time.Minute,
	})

	xRead := make(chan map[string]interface{})

	go func() {
		xra := &redis.XReadArgs{
			Streams: []string{"icinga:dump", "0-0"},
			Block:   5 * time.Second,
		}

		for {
			res, errXR := client.XRead(xra).Result()
			if errXR != nil && errXR != redis.Nil {
				log.WithFields(log.Fields{"host": *host, "port": *port, "error": errXR.Error()}).Fatal(
					"Couldn't XREAD " + xra.Streams[0] + " from " + xra.Streams[1],
				)
			}

			for _, stream := range res {
				for _, message := range stream.Messages {
					xra.Streams[1] = message.ID
					xRead <- message.Values
				}
			}
		}
	}()

	start := time.Now()
	fps := time.NewTicker(time.Second / 30)
	var events []event

	for {
		select {
		case message := <-xRead:
			if message["state"] == "wip" {
				start = time.Now()
			}

			state := message["state"]
			if len(events) < 1 || events[len(events)-1].state != state {
				events = append(events, event{state, nil})
			}

			ev := &events[len(events)-1].typs
			*ev = append(*ev, message["type"])
		case now := <-fps.C:
			secs := float64(now.Sub(start)) / float64(time.Second)
			for _, ev := range events {
				fmt.Printf("%.02f  %s: %v\n", secs, ev.state, ev.typs)
			}

			events = nil
			fmt.Printf("%.02f\r", diffSeconds(now, start))
		}
	}
}

func diffSeconds(lhs, rhs time.Time) float64 {
	return float64(lhs.Sub(rhs)) / float64(time.Second)
}
