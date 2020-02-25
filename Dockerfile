FROM golang:1.14-rc-alpine3.11 AS BUILD

RUN mkdir /cam-event-detector
WORKDIR /cam-event-detector

ADD go.mod .
RUN go mod download

#now build source code
ADD . ./
RUN go build -o /go/bin/cam-event-detector



FROM czentye/opencv-video-minimal:4.2-py3.7.5

RUN apt-get update && \
    apt-get install -y ssh openssh-server

EXPOSE 3000

VOLUME [ "/data" ]

ENV CAM_ID ''
ENV VIDEO_SOURCE_URL ''
ENV EVENT_POST_ENDPOINT ''
ENV EVENT_OBJECT_IMAGE_ENABLE 'true'
ENV EVENT_SCENE_IMAGE_ENABLE 'false'
ENV EVENT_MAX_KEYPOINTS '-1'

COPY --from=BUILD /go/bin/* /bin/
ADD /startup.sh /

ENTRYPOINT /startup.sh

