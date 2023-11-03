FROM golang:latest as builder
WORKDIR /app
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -ldflags="-w -s" -o api cmd/api/main.go

FROM scratch

COPY --from=builder ./app/api / 

CMD [ "/api" ]
