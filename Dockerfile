FROM golang:1.19.2-alpine AS builder

RUN apk update && apk upgrade && \
    apk add --no-cache make bash

WORKDIR /src

COPY . .

# Build
RUN make build

# Using a distroless image from https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static-debian11

COPY --from=builder /src/bin/app /

EXPOSE 8000

CMD ["/app"]