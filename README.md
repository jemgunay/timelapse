# Timelapse Capture

Captures frames from a camera at fixed intervals for a specified duration, then stitches the produced frames into a video.

## Usage

Produce a timelapse over the period of 1 hour, capturing a frame every minute:
```bash
go build
./timelapse -duration=1h -interval=1m
```

Capture frames only (i.e. don't stitch them into a video afterwards):
```bash
./timelapse -cmd="capture" -duration=1m -interval=2s
```

Stitch already captured frames into a video:
```bash
./timelapse -cmd=stitch -stitch_dir="./Thu-17-Jan-2019_22-54-01/frames/" -fps=30 -frame_step=2
```