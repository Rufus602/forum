run:
	go run ./cmd/web/main.go

build:
	docker image build -t forum .

d-run:
	docker run -d -p 8000:8000 --name forumapp forum 
d-id:
	docker ps -a
d-stop:
	docker stop $(docker ps -aq)

