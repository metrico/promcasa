NAME?=promcasa

all:
	#CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o $(NAME)
	CGO_ENABLED=1 GOOS=linux CGO_LDFLAGS="-lm -ldl" go build -a -ldflags '-extldflags "-static"' -tags netgo -installsuffix netgo -o $(NAME)
	#go build -a -ldflags '-extldflags "-static"' -o $(NAME)

debug:
	go build -o $(NAME) 

modules:
	go get ./...

docker:
	docker build -f scripts/Dockerfile -t promcasa .

.PHONY: clean
clean:
	rm -fr $(NAME)
