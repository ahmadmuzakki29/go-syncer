# go-syncer
Let's say you have few machines that run the same Go binary but you want to do atomic transaction with unique ID once at a time. Then you need to have some distribution locking mechanism for that purpose.
go-syncer will synchronize go process on multiple machines. It make  use of mutex locking that identified with an ID. The service should live before any client try to synchronize. The connection use [gRPC](http://www.grpc.io/) with long-live http2 connection to lock process with `syncer.Lock()` until `syncer.Unlock()` called so the race condition technically impossible.
