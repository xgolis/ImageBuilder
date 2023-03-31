FROM golang:latest as build
RUN CGO_ENABLE=0 go install -a -v -work ./cmd/...

FROM alpine:latest
COPY --from=build /gobin/golang-app /usr/local/bin/golang-app
EXPOSE 8080
ENTRYPOINT [ "/usr/local/bin/golang-app" ]