FROM golang:latest as build
WORKDIR /usr/local/bin
COPY . .
WORKDIR /usr/local/bin/cmd/ImageBuilder
RUN go build -o ../../ImageBuilderApp
EXPOSE 8080
ENTRYPOINT [ "/usr/local/bin/ImageBuilderApp" ]