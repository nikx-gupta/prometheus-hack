FROM golang:1.17-alpine as build
WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go mod download && go build -o go_postgresapi

FROM golang:1.17-alpine as final
WORKDIR /app
COPY --from=build /app/* ./
RUN ls -a

EXPOSE 2112
CMD ["/app/go_postgresapi"]