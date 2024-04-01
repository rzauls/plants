watch-tests:
	nodemon --exec "go test ./..." --signal SIGTERM -e go
