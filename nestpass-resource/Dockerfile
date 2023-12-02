FROM golang:1.21.0

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target=/root/.cache/go-build make build

ARG PORT=2000
EXPOSE $PORT

CMD ["make", "run"]