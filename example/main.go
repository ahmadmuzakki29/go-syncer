package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ahmadmuzakki29/go-syncer/client"
	"log"
)

const (
	address = "localhost:9999"
)

var syncer *client.Client

func main() {
	cfg := client.Config{
		EndPoint:    address,
		LockTimeout: time.Duration(10) * time.Second, // auto-unlock when timeout reach
	}

	var err error
	syncer, err = client.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/payment", HandlePayment)
	err = http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// try to execute this in terminal
// curl http://localhost:9000/payment?payment_id=123 & curl http://localhost:9000/payment?payment_id=123
func HandlePayment(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	payment_id := r.FormValue("payment_id")

	// Lock to start processing the payment in atomic way
	err := syncer.Lock(payment_id)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("start processing payment ", payment_id)
	// simulate process that happens with this payment
	time.Sleep(time.Second)
	fmt.Println("finish processing payment ", payment_id)

	//you have to call syncer.Unlock() otherwise future process of the ID will be blocked
	syncer.Unlock(payment_id)
}
