FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY . .
RUN go mod tidy
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/web-app ./main.go

FROM alpine:3.17
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
RUN mkdir -p /usr/bin/static/supporting-documents
COPY --from=build /go/src/app/bin /go/bin
COPY config.toml /usr/bin/
COPY service/mailer/template/upload.html /usr/bin/
COPY service/mailer/template/result.html /usr/bin/
EXPOSE 80
ENTRYPOINT /go/bin/web-app --port 80