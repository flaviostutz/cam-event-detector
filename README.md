# cam-event-detector

Camera image event detector for objects entering, moving, stopping or exiting the scene.

## What it does

In general, the event detector performs the following:

1. Reads continuously a camera image

2. Detects and tracks an object that is moving in the scene

3. Generates events related to the lifecycle of the tracked object: ENTERED, MOVED, STOPPED, EXITED

4. Schedules the event to be sent to a HTTP POST endpoint

5. Posts all the queued events json as soon as possible to the target HTTP endpoint. May have delays in case of network outages or insufficient bandwidth.

## Usage

* create docker-compose.yml:

```yml
version: '3.7'

services:

  cam-event-detector:
    build: .
    image: flaviostutz/cam-event-detector
    ports:
      - 3000:3000
    restart: always
    environment:
      - CAM_ID=cam1
      - VIDEO_SOURCE_URL=http://195.252.80.186:9000/mjpg/video.mjpg
      - EVENT_POST_TARGET=http://localhost:8080/cam1/event
```

## ENVs

* CAM_ID - 'cam_id' for the generated events

* VIDEO_SOURCE_URL - any video stream URL supported by OpenCV. ex.: http://cam1/v.mpeg; rtsp://cam2:554/stream

* EVENT_POST_TARGET - an HTTP POST endpoint whose detected events json contents will be sent asynchronously.

  * example: POST http://anotherserver/cam1/event

```json
{
    "uuid": "ABDC-2342-BCCA",
    "cam_id": "cam1", //as in CAM_ID ENV
    "time": "2020-02-15T15:34:32",
    "type": "MOVED", //ENTERED, STOPPED, EXITED
    "bbox":{"x1":12,"y1":65,"x2":3,"y2":91},
    "speed": {"x":-5,"y":3},
    "image": "ABC3378BAAASA1", //JPEG image encoded in BASE64. if enabled in ENV EVENT_OBJECT_IMAGE_ENABLE
    "tracking": {
        "uuid": "ABDC-2342-BCCA",
        "keypoints": [
            {"ts":"2020-02-15T15:34:01", "bbox":{"x1":12,"y1":45,"x2":23,"y2":81}, "speed":{"x":-2,"y":1}},
            {"ts":"2020-02-15T15:34:01", "bbox":{"x1":12,"y1":65,"x2":3,"y2":91}, "speed":{"x":-5,"y":3}}
        ]
    },
    "scene": {
        "width": 800,
        "height": 600,
        "image": "AAC73621AAACCCDDD" //if enabled in ENV EVENT_SCENE_IMAGE_ENABLE
    },
}
```

* EVENT_OBJECT_IMAGE_ENABLE - if 'true', the cropped detected object image will be included in event payload. defaults to 'true'

* EVENT_SCENE_IMAGE_ENABLE - if 'true', the full scene image will be included in event payload. defaults to 'false'

* EVENT_MAX_KEYPOINTS - the max number of keypoints that will be included in event payload. defaults to '-1' (no limit)

## Samples

* You can test this with video samples from http://github.com/flaviostutz/camera-samples
