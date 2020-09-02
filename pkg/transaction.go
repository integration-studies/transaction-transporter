package pkg

import (
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"log"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	Type        string    `json:"type"`
	SubType     string    `json:"subType"`
	FromAccount string    `json:"fromAccount"`
	ToAccount   string    `json:"toAccount"`
	Value       float64   `json:"value"`
	Time        time.Time `json:"time"`
	DeviceType  string    `json:"deviceType"`
}

func FromLine(line string) (*Transaction, error) {
	s := []rune(line)
	val, errVal := strconv.ParseFloat(strings.TrimSpace(string(s[159:])), 32)
	if errVal != nil {
		log.Printf("Failed to read value : %v", errVal)
		return nil, errVal
	}
	date, errDate := time.Parse(time.RFC1123, string(s[80:109]))
	if errDate != nil {
		log.Printf("Failed to read date : %v", errDate)
		return nil, errDate
	}
	return &Transaction{
		Type:        strings.TrimSpace(string(s[0:10])),
		SubType:     strings.TrimSpace(string(s[11:20])),
		FromAccount: strings.TrimSpace(string(s[21:50])),
		ToAccount:   strings.TrimSpace(string(s[51:80])),
		Value:       val,
		Time:        date,
		DeviceType:  strings.TrimSpace(string(s[110:140])),
	}, nil
}

func (t *Transaction) CloudEvent() v2.Event {
	e := cloudevents.NewEvent()
	e.SetType(eventType(t))
	e.SetSource("tech.claudioed.transaction.file")
	e.SetDataContentType(cloudevents.ApplicationJSON)
	uuid, _ := uuid.NewUUID()
	e.SetID(uuid.String())
	e.SetSubject("new-transaction." + uuid.String())
	d, err := json.Marshal(t)
	if err != nil {
		log.Printf("Failed to parse transaction to json %v", err)
	}
	_ = e.SetData(cloudevents.ApplicationJSON, string(d))
	return e
}

func eventType(t *Transaction) string {
	if "DOC" == t.Type {
		return "tech.claudioed.transaction.doc.create"
	} else if "TED" == t.Type {
		return "tech.claudioed.transaction.ted.create"
	} else if "DOC" == t.Type && "VISA" == t.SubType {
		return "tech.claudioed.transaction.card.visa.create"
	} else {
		return "tech.claudioed.transaction.card.master.create"
	}
}
