FROM golang:1.12.5-alpine3.9 as builder

ENV GO111MODULE "on"
WORKDIR /builder
RUN apk add --update --no-cache git

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch

COPY --from=builder /builder/spanner-operator /spanner-operator

ENV MOCK_DATA_PATH ""

ENTRYPOINT ["/spanner-operator"]