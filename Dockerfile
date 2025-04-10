FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE ${PORT}
CMD ["./main"]


FROM nginx:alpine AS frontend
COPY frontend/ /usr/share/nginx/html/
EXPOSE 5173
CMD ["nginx", "-g", "daemon off;"]