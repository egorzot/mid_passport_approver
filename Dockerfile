FROM golang:1.19.3-alpine as build

WORKDIR /src

COPY . ./

RUN go mod download

RUN CGO_ENABLED=0 go build -o main .

FROM chromedp/headless-shell:109.0.5396.2 as final

WORKDIR /app
COPY --from=build /src/main ./

CMD ["./main"]
