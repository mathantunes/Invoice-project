package main

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mathantunes/arex_project/queuer"
	"github.com/mathantunes/arex_project/services"
)

const (
	// CreateInvoiceQueue Queue Name
	CreateInvoiceQueue = "create_invoice"
	// UpdateInvoiceQueue Queue Name
	UpdateInvoiceQueue = "update_invoice"

	// QueueWorkers Number of concurrent routines listening to Queue
	QueueWorkers = 5
	// DBWorkers Number of concurrent rountines updating the database
	DBWorkers = 5
)

func main() {

	q := queuer.New()

	for i := 0; i < QueueWorkers; i++ {
		go QueueWorkerCreate(q)
	}

}

func QueueWorkerCreate(queue queuer.QueueManager) {
	for {
		url, err := queue.GetQueueURL(CreateInvoiceQueue)
		if err != nil {
			fmt.Println(err)
			continue
		}

		payload, err := queue.ReadFromQueue(url)
		if err != nil {
			fmt.Println(err)
			continue
		}

		invoice := &services.InternalInvoice{}
		err = proto.Unmarshal([]byte(payload), invoice)
		if err != nil {
			fmt.Println(err)
			continue
		}

	}

}
