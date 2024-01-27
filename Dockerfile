FROM golang:1.20 as build
WORKDIR /app

COPY go.mod /app/
COPY go.sum /app/

ARG USERNAME_GITHUB
ARG TOKEN_GITHUB

RUN git config --global url."https://${USERNAME_GITHUB}:${TOKEN_GITHUB}@github.com".insteadOf "https://github.com"

RUN go mod download
RUN go mod tidy

COPY . /app/

ARG BUILD_TAG
ARG BUILD_TIMESTAMP

RUN CGO_ENABLED=0 go build -o /app/main -ldflags="-X 'main.buildTime=${BUILD_TIMESTAMP}' -X 'main.buildVer=${BUILD_TAG}'"

# Deploy

FROM alpine:3.16.0
WORKDIR /app

EXPOSE 7780
EXPOSE 7781

RUN apk update
RUN apk add --no-cache tzdata
ENV cp /usr/share/zoneinfo/Asia/Makassar /etc/localtime
RUN echo "Asia/Makassar" > /etc/timezone

COPY --from=build /app/conf /app/conf
COPY --from=build /app/main /app/sp-payment

CMD ["/app/sp-payment"]