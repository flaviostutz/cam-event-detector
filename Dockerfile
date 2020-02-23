FROM golang:1.13.8 AS BUILD

RUN mkdir /cam-event-detector
WORKDIR /cam-event-detector

ADD go.mod .
ADD go.sum .
RUN go mod download

#now build source code
ADD . ./
RUN go build -o /go/bin/cam-event-detector



FROM golang:1.13.8

EXPOSE 3000

ENV EVENT_POST_ENDPOINT ''

COPY --from=BUILD /go/bin/* /bin/
ADD /startup.sh /

ENTRYPOINT /startup.sh
