# build stage 
FROM golang:1.24.5-alpine AS builder
WORKDIR /app

COPY . .

RUN go build -o /app/main ./main.go

FROM alpine 
WORKDIR /app 

COPY --from=builder /app/main .
CMD [ "/app/main" ]

EXPOSE 8080


