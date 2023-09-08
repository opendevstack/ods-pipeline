FROM busybox

ARG nexusUrl
ARG nexusUsername
ARG nexusPassword
ARG nexusAuth
ARG nexusUrlWithAuth
ARG nexusHost


RUN echo "nexusUrl=${nexusUrl}" >> message.txt
RUN echo "nexusUsername=${nexusUsername}" >> message.txt
RUN echo "nexusPassword=${nexusPassword}" >> message.txt
RUN echo "nexusAuth=${nexusAuth}" >> message.txt
RUN echo "nexusUrlWithAuth=${nexusUrlWithAuth}" >> message.txt
RUN echo "nexusHost=${nexusHost}" >> message.txt

CMD [ "cat", "message.txt" ]
