FROM debian:9.3

# Install utilities
RUN apt update && \
    apt install -y curl

COPY http_broker /root/http_broker

# Service
CMD /root/http_broker