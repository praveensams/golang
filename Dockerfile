FROM golang:alpine

WORKDIR "/mnt"

COPY valie.go /mnt

RUN ["go","run","valie.go"]
