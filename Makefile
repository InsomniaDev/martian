build:
	go build ./cmd/martian/martian.go

run:
	go run ./cmd/martian/martian.go

deploy:
	./scripts/deployDocker.sh