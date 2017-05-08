FROM alpine

ADD ca-certificates.crt /etc/ssl/certs/
ADD build/dp-ness-wda-redirector dp-ness-wda-redirector

ENV BIND_ADDR :8080
EXPOSE 8080

CMD ["/dp-ness-wda-redirector"]
