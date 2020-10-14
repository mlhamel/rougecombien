workdir: $GOPATH/src/github.com/mlhamel/rougecombien
observe: *.go
pubsub: gcloud beta emulators pubsub start --project=rougecombien --host-port=localhost:$PORT
web: go run -race ./cmd/rougecombien