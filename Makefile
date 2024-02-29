COIN_TRACKER_USER=root
COIN_TRACKER_HOST=194.164.59.123
COIN_TRACKER_PATH=/var/www/coin-tracker

coin-tracker-linux:
	CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static" -o bin/coin-tracker-linux-amd64 ./cmd/

deploy-coin-tracker:
	rsync configs/config.yml cmd.sh bin/coin-tracker-linux-amd64 $(COIN_TRACKER_USER)@$(COIN_TRACKER_HOST):$(COIN_TRACKER_PATH)
	ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null $(COIN_TRACKER_USER)@$(COIN_TRACKER_HOST) "cd $(COIN_TRACKER_PATH) && bash ./cmd.sh"

