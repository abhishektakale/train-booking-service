run-server: build-server
	docker-compose up --build -d

stop-server-image: 
	docker-compose down