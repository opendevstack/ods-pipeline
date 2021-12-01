FROM registry.access.redhat.com/ubi8/python-39

COPY app /app
COPY run.sh /app/

WORKDIR /app

RUN pip install --upgrade pip && \
    pip install -r requirements.txt && \
    pip check

USER root
RUN chown -R 1001:0 /app && \
    chmod -R g=u /app && \
    chmod +x /app/run.sh

USER 1001

CMD ["./run.sh"]
