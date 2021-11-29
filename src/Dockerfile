FROM golang:alpine AS build

RUN apk update && \
    apk upgrade && \
    apk add --no-cache --virtual .build-deps curl gcc g++ autoconf automake libtool pkgconfig alpine-sdk git bash git openssh && \
    rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/openvenues/libpostal /libpostal

WORKDIR /libpostal
RUN ./bootstrap.sh && \
    ./configure --prefix=/usr --datadir=/usr/share/libpostal && \
    make -j4 && \
    make install

RUN mkdir -p /src
WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ../publish/app

FROM alpine:latest

COPY --from=build /usr/share/libpostal /usr/share/libpostal
COPY --from=build /usr/lib/libpostal.so /usr/lib/libpostal.so
COPY --from=build /usr/lib/libpostal.so.1 /usr/lib/libpostal.so.1
COPY --from=build /usr/include/libpostal /usr/include/libpostal

WORKDIR /
COPY --from=build /publish .

EXPOSE 80
ENTRYPOINT ["./app"]