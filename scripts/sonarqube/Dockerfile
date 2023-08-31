ARG from=""
FROM $from

USER root
ENV SQ_CONF_DIR="/opt/sonarqube/conf"
RUN chmod -R ugo+wrX "${SQ_CONF_DIR}"

USER sonarqube
RUN echo "sonar.forceAuthentication=true" >> conf/sonar.properties

# Uncomment the following line to see debug output, e.g. in case SonarQube is
# not starting up properly.
# RUN echo "sonar.log.level=DEBUG" >> conf/sonar.properties
