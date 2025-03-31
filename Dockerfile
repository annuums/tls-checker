FROM golang:1.23 AS builder

WORKDIR /builder

COPY go.mod ./
RUN go mod download

COPY . .

RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "x86_64" ]; then echo "amd64" > .arch; \
    elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then echo "arm64" > .arch; \
    else echo "unsupported arch: $ARCH" && exit 1; fi

RUN ARCH=$(cat .arch) && make build-linux.$ARCH
RUN ARCH=$(cat .arch) && mv deploy/bin/main.linux.$ARCH deploy/bin/main

FROM alpine:latest

WORKDIR /app

COPY --from=builder /builder/deploy/bin/main ./main

ENTRYPOINT ["./main"]
