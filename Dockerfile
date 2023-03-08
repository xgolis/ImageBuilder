FROM golang:latest
COPY . .
CMD ["ls", "-la"]