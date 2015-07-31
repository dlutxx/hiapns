package hiapns

import (
	"github.com/dlutxx/apns"
)

type Request struct {
	Notif *apns.Notification
	App   string
}

type Worker struct {
	hub   *Hub
	ReqCh chan Request
}

func NewWorker(hub *Hub) *Worker {
	w := &Worker{
		hub,
		make(chan Request, 16),
	}
	go w.runLoop()
	return w
}

func (w *Worker) runLoop() {
	for req := range w.ReqCh {
		req.Notif.Identifier = w.hub.cnt.Next()
		w.hub.Send(req.App, req.Notif)
	}
}
