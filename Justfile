watch-tests:
	nodemon --exec "go test ./... -v" --signal SIGTERM -e go
