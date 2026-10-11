FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN go test ./...
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/server ./server/cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/server /server
ENV PORT=8000
EXPOSE 8000
ENTRYPOINT ["/server"]
