FROM python:3

WORKDIR /
COPY entrypoint.sh entrypoint.sh
COPY main main
COPY conf/app.conf conf/app.conf
COPY static static
RUN chmod +x main entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]