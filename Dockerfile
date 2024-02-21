FROM golang:1.19.2-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o expense-tracker cmd/expense-tracker/main.go

FROM alpine
WORKDIR /app
COPY --from=build /app/expense-tracker .
ENTRYPOINT [ "./expense-tracker"]