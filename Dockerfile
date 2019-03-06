FROM golang:1.11.5

RUN mkdir -p /go/src/github.com/mozz100/tohora
WORKDIR /go/src/github.com/mozz100/tohora
COPY . .

RUN GOBIN=/usr/local/bin go install -v ./tohora.go

EXPOSE 80
ENTRYPOINT ["tohora", "80"]