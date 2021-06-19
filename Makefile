.PHONY: test
test:
	go test ./...

profile.out:
	go test -v -coverpkg=./... -coverprofile=profile.out ./...

.PHONY: cov
cov: profile.out
	go tool cover -func profile.out

.PHONY: webcov
webcov: profile.out
	go tool cover -html profile.out

.PHONY: clean
clean:
	rm *.out

.PHONY: install
install:
	go install

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: docker/build
docker/build:
	docker build -t nuetoban/tgmailing .

.PHONY: docker/push
docker/push:
	docker push nuetoban/tgmailing

.PHONY: docker/build/push
docker/build/push: docker/build docker/push