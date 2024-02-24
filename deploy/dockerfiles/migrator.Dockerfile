FROM ubuntu:20.04

ENV APP_DIR /opt/app

WORKDIR $APP_DIR

RUN apt-get update \
    && groupadd -r web \
    && useradd -d $APP_DIR -r -g web web \
    && chown web:web -R $APP_DIR \
    && apt-get install -y netcat-traditional \
    && apt-get install -y acl

ADD https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 /bin/goose

COPY /db/migrations/*.sql db/migrations/
COPY deploy/scripts/migrator-start.sh .

RUN setfacl -R -m u:web:rwx /bin/goose
RUN setfacl -R -m u:web:rwx $APP_DIR/migrator-start.sh

USER web

ENTRYPOINT ["bash", "migrator-start.sh"]