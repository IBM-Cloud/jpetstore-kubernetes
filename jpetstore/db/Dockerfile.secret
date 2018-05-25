FROM mysql:5.5
ADD mysql5/*.sql /docker-entrypoint-initdb.d/
COPY mysql5/docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

