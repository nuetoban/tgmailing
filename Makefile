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