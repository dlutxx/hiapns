package main

import (
	"flag"
	"github.com/dlutxx/hiapns"
	"log"
)

var (
	topic   = flag.String("topic", "apns", "NSQ topic name")
	channel = flag.String("channel", "sender", "NSQ channel name")
	nsqd    = flag.String("nsqd", "127.0.0.1:4150", "NSQD tcp address")
	cfgfile = flag.String("cfg", "/etc/hiapns.json", "config file path")
)

func main() {
	flag.Parse()

	hub, err := hiapns.NewHubFromCfgFile(*cfgfile)
	if err != nil {
		log.Fatal(err)
	}

	worker := hiapns.NewWorker(hub)
	if err = hiapns.FeedWorkerWithNSQ(worker, *topic, *channel, *nsqd); err != nil {
		log.Fatal(err)
	}

	// block forever
	select {}
}
