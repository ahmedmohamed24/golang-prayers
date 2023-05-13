FROM golang:latest

RUN apt-get update
RUN apt-get upgrade -y

ENV GOBIN /go/bin

WORKDIR /app

RUN go mod init github.com/ahmedmohamed24/app

RUN go mod download

RUN go install -mod=mod github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build -buildvcs=false -o runserver ." --command=./runserver
