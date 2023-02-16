# docker build -t go_rest_api . -f Dockerfile-Builder
# docker run -dit --name go_rest_api -p 9000:9000 go_rest_api
# docker run -dit --name go_rest_api -p 9000:9000 6060:6060 go_rest_api

FROM golang:1.18 As builder

RUN apt-get update && apt-get install bash && apt-get install -y --no-install-recommends ca-certificates

WORKDIR /app

COPY . .
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w' -o go_rest_api

# Add a new user "dev" with user id 1000
RUN useradd -u 1001 dev
# Change to non-root privilege
USER dev

#FROM scratch
FROM alpine

WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/cmd/go_rest_api .

WORKDIR /app/resource
COPY --from=builder /app/resource/application.yml .

CMD ["/app/go_rest_api"]