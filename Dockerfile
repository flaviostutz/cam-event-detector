FROM golang:1.13.8 AS BUILD

RUN mkdir /cam-event-detector
WORKDIR /cam-event-detector

ADD go.mod .
ADD go.sum .
RUN go mod download

#now build source code
ADD . ./
RUN go build -o /go/bin/cam-event-detector



FROM czentye/opencv-video-minimal:4.2-py3.7.5

EXPOSE 3000

ENV CAM_ID ''
ENV VIDEO_SOURCE_URL ''
ENV EVENT_POST_ENDPOINT ''
ENV EVENT_OBJECT_IMAGE_ENABLE 'true'
ENV EVENT_SCENE_IMAGE_ENABLE 'false'
ENV EVENT_MAX_KEYPOINTS '-1'

COPY --from=BUILD /go/bin/* /bin/
ADD /startup.sh /

ENTRYPOINT /startup.sh

