FROM golang:1.21 as builder

ENV CGO_ENABLED=0
WORKDIR /src
COPY ingest.go ingest.go
RUN go build -o oc_ingest ingest.go

FROM scratch
COPY --from=builder /src/oc_ingest /oc_ingest
ENTRYPOINT ["/oc_ingest"]
