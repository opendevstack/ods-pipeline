FROM registry.access.redhat.com/ubi8/openjdk-11:1.13

WORKDIR /opt/sbt-sample-app
COPY dist /opt/sbt-sample-app

USER root
RUN chown -R jboss: /opt/sbt-sample-app
USER jboss

EXPOSE 9000

ENTRYPOINT ["/opt/sbt-sample-app/bin/sbt-sample-app"]
CMD []
