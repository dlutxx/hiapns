package hiapns

import (
	"github.com/bitly/go-nsq"
	"log"
)

func FeedWorkerWithNSQ(worker *Worker, nsqTopic, nsqChannel, nsqdAddr string) error {
	nsqcfg := nsq.NewConfig()
	csm, err := nsq.NewConsumer(nsqTopic, nsqChannel, nsqcfg)
	if err != nil {
		return err
	}

	csm.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
		req, err := ParseRequestFromJSON(msg.Body)
		log.Println("Recv: ", string(msg.Body))
		if err == nil {
			worker.ReqCh <- req
		} else {
			log.Println("ERR:", err)
		}
		msg.Finish()
		return nil
	}))

	return csm.ConnectToNSQD(nsqdAddr)
}
