all:
		GOOS=linux GOARCH=amd64 go build -o netsetgo cmd/netsetgo.go

clean:
		rm -f netsetgo
