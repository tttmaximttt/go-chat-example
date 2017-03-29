FROM ubuntu:14.04
MAINTAINER Maksym Radko tttmaximttt@gmail.com
RUN apt-get update
RUN apt-get install -y ca-certificates
CMD curl https://www.google.com
ADD ./ .
EXPOSE 8080
ENTRYPOINT ["/gochat"]