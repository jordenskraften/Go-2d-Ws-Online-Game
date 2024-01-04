run:
	cd cmd && go run main.go

docker:
	docker build -t go-ws-2d-game .