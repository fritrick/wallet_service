.PHONY: build
build:
	go build -o ./bin/wallet_service ./src

.PHONY: install-modules
install-modules:
	go get -u ./...

.PHONY: rundb
rundb:
	./misc/rundb.sh

.PHONY: stopdb
stopdb:
	docker stop wallet-postgres && docker rm wallet-postgres

.PHONY: run
run: build
	./bin/wallet_service
