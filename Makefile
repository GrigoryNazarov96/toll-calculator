gateway:
	@go build -o bin/gateway ./gateway
	@./bin/gateway &

obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu &

receiver:
	@go build -o bin/receiver ./data_receiver
	@./bin/receiver &

calc:
	@go build -o bin/calc ./distance_calculator
	@./bin/calc &

aggregator:
	@go build -o bin/aggregator ./aggregator
	@./bin/aggregator &

proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto

all: proto gateway aggregator calc receiver obu

stop: 
	@pkill -f gateway
	@pkill -f aggregator
	@pkill -f receiver
	@pkill -f obu
	@pkill -f calc
	
.PHONY: obu aggregator gateway receiver calc