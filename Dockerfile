##Builder Image
FROM golang:1.14 as builder
ENV GO111MODULE=on
COPY . /transaction-file-transporter
WORKDIR /transaction-file-transporter
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/application

#s Run Image
FROM scratch
COPY --from=builder /transaction-file-transporter/bin/application application
ENTRYPOINT ["./application"]