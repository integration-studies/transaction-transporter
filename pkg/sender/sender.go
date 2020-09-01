package sender

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"log"
	"sync"
	"transaction-transporter/pkg"
)

type Sender struct {
	client cloudevents.Client
	ctx context.Context
}

func NewSender(client cloudevents.Client,ctx context.Context) *Sender {
	return &Sender{
		client: client,
		ctx:    ctx,
	}
}

func (s *Sender) Send(t *pkg.Transaction, wg *sync.WaitGroup) {
	defer wg.Done()
	res := s.client.Send(s.ctx, t.CloudEvent())
	if cloudevents.IsUndelivered(res) {
		log.Printf("Failed to send: %v", res)
	} else {
		var httpResult *cehttp.Result
		cloudevents.ResultAs(res, &httpResult)
		log.Printf("Sent %d with status code", httpResult.StatusCode)
	}
}




