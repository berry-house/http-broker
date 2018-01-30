FROM debian:9.3

# Install utilities
RUN apt update && \
    apt install -y curl

COPY http-broker /root/http-broker

# Service
CMD /root/http-broker