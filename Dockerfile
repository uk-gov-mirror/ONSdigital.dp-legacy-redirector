FROM alpine

ADD ca-certificates.crt /etc/ssl/certs/
ADD build/dp-legacy-redirector dp-legacy-redirector

ENV BIND_ADDR :8080
EXPOSE 8080

CMD ["/dp-legacy-redirector"]
