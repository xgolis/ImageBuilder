FROM golang:latest as build
WORKDIR /app
COPY . .
RUN cd cmd/ImageBuilder && \
    go build -o ../../ImageBuilder

FROM redhat/ubi8:latest
WORKDIR /app
COPY --from=build /app/ImageBuilder .
EXPOSE 8082
ENTRYPOINT [ "./ImageBuilder" ]