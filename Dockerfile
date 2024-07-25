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

COPY --from=build /app/bin/beats ./bin/
COPY --from=build /app/data/ ./data/


CMD ["/app/bin/beats"]
