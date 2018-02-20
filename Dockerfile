FROM debian:9.3

# Install utilities
RUN apt update && \
    apt install -y curl

COPY http_broker /root/http_broker
COPY conf   /root/conf

RUN mkdir /logs

# Service
CMD /root/http_broker                           \
    -port               $PORT                   \
    -runningMode        prod                    \
    -loggerConfigFile   /root/conf/logger.json  \
    -databaseAddress    $DATABASE_ADDRESS       \
    -databaseName       $DATABASE_NAME          \
    -databaseUsername   $DATABASE_USERNAME      \
    -databasePassword   $DATABASE_PASSWORD