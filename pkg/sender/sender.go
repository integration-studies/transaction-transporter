package sender

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"log"
	"sync"
	"transaction-transporter/pkg"
)

type Sender struct {
	client cloudevents.Client
	ctx    context.Context
}

func NewSender(client cloudevents.Client, ctx context.Context) *Sender {
	return &Sender{
		client: client,
		ctx:    ctx,
	}
}

func (s *Sender) Send(t *pkg.Transaction, wg *sync.WaitGroup) {
	defer wg.Done()
	if result := s.client.Send(s.ctx, t.CloudEvent()); cloudevents.IsUndelivered(result) {
		log.Printf("Failed to send: %s", result.Error())
	} else if cloudevents.IsACK(result) {
		log.Printf("Sent: %v", t)
	} else if cloudevents.IsNACK(result) {
		log.Printf("Sent but not accepted: %s", result.Error())
	}
}
