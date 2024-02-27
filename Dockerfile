FROM amazoncorretto:17

COPY --from=golang:1.21-alpine /usr/local/go/ /usr/local/go/
 
 RUN yum update -y && yum install graphviz -y

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

EXPOSE 9888
CMD ["app"]