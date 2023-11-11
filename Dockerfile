# Build the application from source
FROM golang:1.21.1-alpine AS build

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app

# Deploy the application binary into a lean image
FROM gcr.io/distroless/static-debian11

ENV TZ="Asia/Bangkok"

COPY --from=build /bin/app /bin
COPY .env.prod /bin

EXPOSE 3000

ENTRYPOINT ["/bin/app", "/bin/.env.prod"]