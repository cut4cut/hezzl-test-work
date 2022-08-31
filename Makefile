run: ### Run docker-compose
	docker-compose up --build -d app && docker-compose logs -f
.PHONY: run

down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: down

protoc: ### Compile protobuf
	protoc --go_out=. --go_opt=paths=source_relative \
    	   --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		    ./internal/controller/rpc/proto/*.proto     
.PHONY: protoc

