FROM registry.access.redhat.com/ubi8/python-39

# expecting args to already be url encoded
ARG nexusUrlWithAuth
ARG nexusHost

COPY requirements.txt /app/
COPY docker/run.sh /app/

WORKDIR /app

RUN if [ ! -z ${nexusUrlWithAuth} ]; \
    then \
      pip install -i ${nexusUrlWithAuth}/repository/pypi-all/simple --trusted-host ${nexusHost} --upgrade pip && \
      pip install -i ${nexusUrlWithAuth}/repository/pypi-all/simple --trusted-host ${nexusHost} -r requirements.txt; \
    else \
      pip install --upgrade pip && \
      pip install -r requirements.txt; \
    fi && \
    pip check

COPY src /app

USER root
RUN chown -R 1001:0 /app && \
    chmod -R g=u /app && \
    chmod +x /app/run.sh

USER 1001

CMD ["./run.sh"]
