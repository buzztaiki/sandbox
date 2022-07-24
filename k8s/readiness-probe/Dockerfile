FROM golang

WORKDIR /go/src/app
COPY go.mod *.go ./

RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT ["readiness-probe"]
