FROM busybox

ARG firstArg
ARG secondArg


RUN echo "firstArg=${firstArg}" >> message.txt
RUN echo "secondArg=${secondArg}" >> message.txt

CMD [ "cat", "message.txt" ]
