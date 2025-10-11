FROM --platform=$BUILDPLATFORM golang:1.25-trixie

# cpu architecture
ARG TARGETARCH

# initialize work dir
WORKDIR /build

# download modules
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY *.go ./

# build
RUN GOOS=linux GOARCH=$TARGETARCH go build -o reddit-gallery-api

# expose network port
EXPOSE 9361

# run
CMD ["/build/reddit-gallery-api"]

# example multi-platform build & push:
# docker build --platform linux/amd64,linux/arm64,linux/arm -t pritoranta/reddit-gallery-api:1.3.0-trixie .
# docker push pritoranta/reddit-gallery-api:1.3.0-trixie
