FROM golang:1.11.5

RUN mkdir -p /go/src/github.com/mozz100/screenserve
WORKDIR /go/src/github.com/mozz100/screenserve
COPY . .

RUN GOBIN=/usr/local/bin go install -v ./screenserve.go

EXPOSE 80
ENTRYPOINT ["screenserve", "80"]