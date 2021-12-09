# syntax=docker/dockerfile:1
FROM alpine:3
MAINTAINER Peter Artemiev
ENV PSQL_CONTAINER="PSQL"
ENV PSQL_PASS="qwerty"
ENV PSQL_USER="postgres"
ENV MONGO_CONTAINER="MONGO"
ADD calculator .
RUN chmod u+x calculator
EXPOSE 4969
CMD /calculator