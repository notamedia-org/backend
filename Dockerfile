FROM golang:1.21 as build

WORKDIR /app

ENV GOOS=linux
ENV CGO_ENABLED=0

RUN apt update && apt install -y jq

COPY . .

RUN export VERSION=$(cat package.json | jq -r .version) && \
    make build


FROM alpine:latest

WORKDIR /app

COPY --from=build /app/bin/v1 ./bin/
COPY --from=build /app/data/ ./data/

RUN curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
RUN echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
RUN apt-get update
RUN apt-get install -y migrate

RUN migrate -source ./migrations -database postgres://postgres:root@localhost:5432/template2?sslmode=disable up

CMD ["/app/bin/beats"]
