FROM gomicro/goose

WORKDIR /migrations/
ADD ./migrations/*.sql ./
ADD ./goose.sh ./

RUN chmod +x ./goose.sh

ENTRYPOINT "./goose.sh"