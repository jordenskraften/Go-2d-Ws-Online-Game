run:
	cd cmd && go run main.go

docker:
	docker build -t my-go-app .