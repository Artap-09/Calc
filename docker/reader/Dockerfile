# syntax=docker/dockerfile:1
FROM centos:8
MAINTAINER Peter Artemiev
ENV PSQL_CONTAINER="PSQL"
ENV PSQL_PASS="qwerty"
ENV PSQL_USER="postgres"
ENV MONGO_CONTAINER="MONGO"
ADD readerdb .
RUN chmod u+x reader
EXPOSE 4979
CMD /readerdb