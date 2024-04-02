FROM golang:1.22 as build

WORKDIR /app

COPY . . 
RUN go mod download && go mod verify
RUN CGO_ENABLED=0 go build -o ./coursehunt

FROM gcr.io/distroless/static-debian12
COPY --from=build /app /

EXPOSE 8080

CMD ["./coursehunt"]
