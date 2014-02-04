FROM ubuntu:12.04
RUN apt-get update -q
RUN apt-get install -qy build-essential curl git
RUN curl -s https://go.googlecode.com/files/go1.2.src.tar.gz | tar -v -C /usr/local -xz
RUN cd /usr/local/go/src && ./make.bash --no-clean 2>&1
ENV PATH /usr/local/go/bin:$PATH
ENV GOPATH /go
RUN echo "deb http://archive.ubuntu.com/ubuntu precise universe" >> /etc/apt/sources.list
RUN apt-get update -q
RUN apt-get install -y mercurial
RUN go get github.com/Clever/riemann-gearman
RUN cp $GOPATH/bin/riemann-gearman /usr/local/bin/
