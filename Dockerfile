FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/gmf-core ./cmd/gmf-core \
 && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/gmf-worker ./cmd/gmf-worker

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/gmf-core /gmf-core
COPY --from=build /out/gmf-worker /gmf-worker
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/gmf-core"]
