import numpy as np
import cv2

print("OPENING VIDEO CAPTURE")
# cap = cv2.VideoCapture("http://87.57.111.162/mjpg/video.mjpg")
# cap = cv2.VideoCapture("samples/cars1.mp4")
cap = cv2.VideoCapture("rtsp://dummy-rtsp-relay:8554/stream")

while(True):
    # Capture frame-by-frame
    ret, frame = cap.read()

    # Our operations on the frame come here
    gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)

    # Display the resulting frame
    print("SHOW RESULTS")
    cv2.imshow('frame',gray)
    if cv2.waitKey(1) & 0xFF == ord('q'):
        break

# When everything done, release the capture
cap.release()
cv2.destroyAllWindows()

