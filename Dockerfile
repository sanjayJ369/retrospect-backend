# build stage 
FROM golang:1.24.5-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .

RUN go build -o /app/main ./main.go


FROM alpine 
WORKDIR /app 

COPY --from=builder /app/main .
COPY --from=builder /go/bin/migrate .

COPY app.env .
COPY ./db/migration ./migration
COPY start.sh .

CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]

EXPOSE 8080
