FROM golang:1.24.2-alpine3.21 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd/albums -o /gin-app


FROM  golang:1.24.2-alpine3.21 

COPY --from=build /gin-app /gin-app

EXPOSE 8090

ENTRYPOINT ["/gin-app"]



