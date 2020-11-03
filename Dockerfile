FROM golang:alpine AS builder

RUN apk update
RUN apk add --no-cache git
WORKDIR /app/

RUN go get github.com/pin/tftp
RUN go get github.com/go-git/go-git
RUN go get github.com/sparrc/go-ping
RUN go get github.com/docker/docker/client
RUN go get github.com/shirou/gopsutil/cpu

COPY main.go main.go
RUN CGO_ENABLED=0 go build -o /main


FROM scratch
COPY --from=builder /main /main
ENTRYPOINT ["/main"]
