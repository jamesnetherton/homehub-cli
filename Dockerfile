FROM golang:alpine as compile

RUN apk --no-cache add make git

COPY . src/github.com/jamesnetherton/homehub-cli

ENV CGO_ENABLED=0
RUN cd src/github.com/jamesnetherton/homehub-cli && \
    make build && \
    cp build/homehub-cli /go/bin/homehub-cli


FROM scratch

COPY --from=compile /go/bin/homehub-cli /homehub-cli

ENTRYPOINT ["/homehub-cli"]
