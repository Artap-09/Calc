# syntax=docker/dockerfile:1
FROM centos:latest
MAINTAINER Peter Artemiev
ENV PSQL_CONTAINER="PSQL"
ENV PSQL_PASS="qwerty"
ENV PSQL_USER="postgres"
ENV MONGO_CONTAINER="MONGO"
ADD readerdb .
EXPOSE 4979
CMD /readerdb