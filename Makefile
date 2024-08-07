INPUT := $(arg)
run_server:
	go run server/*.go
run_client:
	go run client/*.go $(shell pwd)/$(INPUT)