pm2 stop coin-tracker
pm2 start coin-tracker-linux-amd64 --name=coin-tracker -- -config=./config.yml
