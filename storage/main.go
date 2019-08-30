package main

import (
	"encoding/base64"
	"fmt"
	"log"

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
)

func main() {

	q := queuer.New()
	// Create Insertion QUEUE
	q.CreateQueue(CreateInvoiceQueue)
	// Create Update QUEUE
	q.CreateQueue(UpdateInvoiceQueue)

	for i := 0; i < QueueWorkers; i++ {
		go QueueWorkerCreate(q)
		go QueueWorkerUpdate(q)
	}

	select {}
}

// QueueWorkerUpdate Worker to Queue the SQS and process updates
func QueueWorkerUpdate(queue queuer.QueueManager) {
	for {
		url, err := queue.GetQueueURL(UpdateInvoiceQueue)
		if err != nil {
			fmt.Println(err)
			continue
		}

		payload, err := queue.ReadFromQueue(url)
		if err != nil {
			// fmt.Println(err)
			continue
		}

		invoice := &services.InternalInvoice{}
		decoded, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			log.Println(err)
			continue
		}
		err = proto.Unmarshal([]byte(decoded), invoice)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if invoice.GetType() == services.InvoiceType_AR {
			QueryInvoice(invoice, updateAR)
		} else if invoice.GetType() == services.InvoiceType_AP {
			QueryInvoice(invoice, updateAP)
		}
	}
}

// QueueWorkerCreate Worker to Queue the SQS Queue and process insertions
func QueueWorkerCreate(queue queuer.QueueManager) {
	for {
		url, err := queue.GetQueueURL(CreateInvoiceQueue)
		if err != nil {
			fmt.Println(err)
			continue
		}

		payload, err := queue.ReadFromQueue(url)
		if err != nil {
			// fmt.Println(err)
			continue
		}

		invoice := &services.InternalInvoice{}
		decoded, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			log.Println(err)
			continue
		}
		err = proto.Unmarshal([]byte(decoded), invoice)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if invoice.GetType() == services.InvoiceType_AR {
			QueryInvoice(invoice, insertAR)
		} else if invoice.GetType() == services.InvoiceType_AP {
			QueryInvoice(invoice, insertAP)
		}
	}

}
