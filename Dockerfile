FROM ubuntu:16.04
LABEL maintainer="jademcosta@gmail.com" \
      version="1.1"

ENV LIBVIPS_VERSION 8.5.7

RUN apt-get update && \
  apt-get install -y build-essential libxml2-dev libfftw3-dev \
	libmagickwand-dev libopenexr-dev liborc-0.4-0 \
	libgsf-1-dev libexpat1-dev \
	libglib2.0-dev libjpeg-dev libtiff-dev zlib1g-dev liblcms2-dev \
  libpng-dev libmagickcore-dev libfreetype6-dev libpango1.0-dev \
  libfontconfig1-dev libice-dev gettext pkg-config libexif-gtk-dev \
  python-all-dev python-dev libmatio-dev libcfitsio-dev \
  libopenslide-dev libwebp-dev libgif-dev libpoppler-glib-dev librsvg2-dev \
  automake libtool swig gtk-doc-tools gcc git libc6-dev make \
  ca-certificates wget

RUN cd /tmp && \
  wget https://github.com/jcupitt/libvips/releases/download/v${LIBVIPS_VERSION}/vips-${LIBVIPS_VERSION}.tar.gz && \
  tar zvxf vips-${LIBVIPS_VERSION}.tar.gz && \
  cd vips-${LIBVIPS_VERSION} && \
  ./configure --enable-debug=no --without-python && \
  make && \
  make install && \
  ldconfig && \
  cd /tmp && \
  rm -rf vips-${LIBVIPS_VERSION} && \
  rm vips-${LIBVIPS_VERSION}.tar.gz

ENV GO_VERSION 1.8.3

RUN cd /tmp && \
  wget -O go${GO_VERSION}.tar.gz https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz && \
  tar -C /usr/local -xzf go${GO_VERSION}.tar.gz && \
  rm go${GO_VERSION}.tar.gz

ENV GOPATH /go
RUN mkdir -p "$GOPATH/src/github.com/jademcosta/melanite" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR /go/src/github.com/jademcosta/melanite

EXPOSE 8080
