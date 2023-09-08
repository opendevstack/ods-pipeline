FROM busybox

ARG firstArg
ARG secondArg

COPY message.txt message.txt

RUN echo "firstArg=${firstArg}" >> message.txt
RUN echo "secondArg=${secondArg}" >> message.txt

CMD [ "cat", "message.txt" ]
