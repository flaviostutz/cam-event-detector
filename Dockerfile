FROM flaviostutz/opencv-golang:1.0.2

RUN apk add --no-cache git

#dependency for lib github.com/flaviostutz/sort
WORKDIR /tmp
RUN apk add --no-cache cmake && \
    wget -O metis.tar.gz http://glaros.dtc.umn.edu/gkhome/fetch/sw/metis/metis-5.1.0.tar.gz && \
    tar -xvf metis.tar.gz && \
    cd metis-5.1.0 && \
    make config && make install && \
    rm -rf /tmp/*

# RUN apk add --no-cache py-pip jpeg-dev python-dev
# RUN pip install Pillow torchvision
# RUN pip install torch torchvision -f https://download.pytorch.org/whl/torch_stable.html

EXPOSE 3000
VOLUME [ "/data" ]

ENV CAM_ID ''
ENV VIDEO_SOURCE_URL ''
ENV EVENT_POST_ENDPOINT ''
ENV IMAGE_POST_ENDPOINT ''
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
