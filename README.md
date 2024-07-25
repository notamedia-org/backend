## Build

```bash
make build
```

## Start

```bash
bin/beats
```


## Docker
docker-compose.dev.yml example

```yaml
version: '3.5'

services:
  beats:
    build:
      context: .
      dockerfile: Dockerfile
    healthcheck:
      test: wget --spider http://localhost:8887/health/readiness
      interval: 10s
      timeout: 10s
      retries: 10
```

# Environment

| env     | default value | description      |
|:--------|:--------------|:-----------------|
| PORT    | 8887          | Service port     |
| DEBUG   |               | Debug string     |
| APP_ENV | development   | Type of instance | 