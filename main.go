package main

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"transaction-transporter/pkg"
	"transaction-transporter/pkg/sender"
)


func main() {
	ctx := cloudevents.ContextWithTarget(context.Background(), os.Getenv("BROKER_URL"))
	hp, _ := cloudevents.NewHTTP()
	client,_ := cloudevents.NewClient(hp, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	sender := sender.NewSender(client,ctx)
	data, err := ioutil.ReadFile(os.Getenv("FILE_PATH"))
	if err != nil {
		log.Printf("Failed to read file %v", os.Getenv("FILE_PATH"))
	}
	var wg sync.WaitGroup
	for _, line := range strings.Split(string(data), "\n") {
		wg.Add(1)
		t,err := pkg.FromLine(line)
		if err != nil{
			log.Printf("Failed to parse transaction %v", err)
		}
		log.Printf("Failed to read file %v", t)
		go sender.Send(t,&wg)
	}
	wg.Wait()
}
