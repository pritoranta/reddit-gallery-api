FROM golang:1.24.6-bookworm AS base

# initialize work dir
WORKDIR /build

# download modules
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY *.go ./

# build
RUN go build -o reddit-gallery-api

# expose network port
EXPOSE 9361

# run
CMD ["/build/reddit-gallery-api"]
