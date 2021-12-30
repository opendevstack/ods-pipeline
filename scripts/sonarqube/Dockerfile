FROM sonarqube:8.4-community

RUN echo "sonar.forceAuthentication=true" >> conf/sonar.properties
RUN echo "sonar.search.javaAdditionalOpts=-Dbootstrap.system_call_filter=false" >> conf/sonar.properties
# see https://stackoverflow.com/questions/47210302/sonarqube-6-7-failed-to-start-because-config-seccomp-not-compiled-into-kernel

