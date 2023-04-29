FROM golang:1.19-alpine AS builder 
WORKDIR /app 
COPY . .
RUN apk add build-base && go build -o forum cmd/web/main.go 

FROM alpine:latest 
WORKDIR /app 
COPY --from=builder /app .

EXPOSE 8000 
LABEL name = "FORUM"
LABEL authors = "bshayakh, rakhmeto, dizdibay"
LABEL release_date = "2023.04.30"
CMD ["/app/forum"]
