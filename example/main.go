package main

import (
	"fmt"
	"github.com/ahmadmuzakki29/go-syncer"
	"net/http"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	syncer.Init(address)
	http.HandleFunc("/payment", HandlePayment)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println(err)
	}
}

// try to curl http://localhost:9000/payment?payment_id=123
// in multiple terminal in same time
func HandlePayment(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	payment_id := r.FormValue("payment_id")
	err := syncer.Lock(payment_id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("start processing payment ", payment_id)
	time.Sleep(time.Duration(2) * time.Second)
	fmt.Println("finish processing payment ", payment_id)
	syncer.Unlock(payment_id)
}
