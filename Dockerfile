FROM flaviostutz/opencv-golang:1.0.1

EXPOSE 3000

VOLUME [ "/data" ]

WORKDIR /tmp
RUN apk add curl && \
    curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py && \
    python get-pip.py

RUN pip install Pillow torchvision

ENV CAM_ID ''
ENV VIDEO_SOURCE_URL ''
ENV EVENT_POST_ENDPOINT ''
ENV EVENT_OBJECT_IMAGE_ENABLE 'true'
ENV EVENT_SCENE_IMAGE_ENABLE 'false'
ENV EVENT_MAX_KEYPOINTS '-1'

RUN mkdir /cam-event-detector
WORKDIR /cam-event-detector

ADD go.mod .
RUN go mod download

#now build source code
ADD . ./
RUN go build -o /usr/bin/cam-event-detector

ADD startup.sh /

# CMD [ "/startup.sh" ]
