FROM postgres:14.5

RUN apt-get update \
    && apt-get install wget -y \
    && apt-get install postgresql-14-postgis-3 -y \
    && apt-get install postgis -y

CMD ["/usr/local/bin/docker-entrypoint.sh","postgres"]
