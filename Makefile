build-server: 
	go build -o ./out/server ./cmd/server

run-server: build-server
	docker-compose up --build -d

stop-server-image: 
	docker-compose down