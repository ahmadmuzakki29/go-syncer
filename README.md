# go-syncer
Let's say you have few machines that run the same Go binary and you want to do atomic transaction with unique ID in each server once at a time. You need to have kind of distribution locking mechanism for that purpose. With go-syncer you can synchronize go process on multiple machines depending on an ID. It make  use of mutex locking that identified with that ID. 

Go-syncer consists of server and client. The server must live before any client try to synchronize. The connection among server and client use [gRPC](http://www.grpc.io/) with long-live http2 connection to lock process with `syncer.Lock()` until `syncer.Unlock()` called, so the race condition technically impossible.

# Getting started
To install the library, run:

`go get -u github.com/ahmadmuzakki29/go-syncer`

Build the server binary :
```
cd $GOPATH/src/github.com/ahmadmuzakki29/go-syncer
go get -u
go build server/syncer.go
```

Run the server :

`./syncer [options]`

options :
 
--debug : Init debug mode

--port &lt;port&gt; : Specify port for the service. default is 9999



The following is a simple example which show how to implement it as a client. We gonna create dummy http server which will process an ID - in this example is `payment_id`. After our http server is listening we gonna hit the API twice at the same time to simulate multiple process in different machines . 
Even though the http service is executed in single machine but we can assure you that the locking mechanism will do well in real multi machines because it use http2 transport to communicate with the **go-syncer** server.
```go
package main

import (
	"fmt"
	"time"
	"net/http"
	"github.com/ahmadmuzakki29/go-syncer"
)

const (
	address = "localhost:9999"
)

func main() {
	syncer.Init(address)
	http.HandleFunc("/payment", HandlePayment)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println(err)
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
```

Then in terminal execute this `curl http://localhost:9000/payment?payment_id=123 & curl http://localhost:9000/payment?payment_id=123`. That command will hit the API twice asynchronously if you run **go-syncer** server in debug mode you'd see `suspending 123  buffer: 2` it means that **go-syncer** hold the incoming `123` request and process count at the moment is 2 processes.
