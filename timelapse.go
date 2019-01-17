package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gocv.io/x/gocv"
)

const (
	// output image formats
	imageFormat = "jpg"
	videoFormat = "avi"
	// camera resolution
	camWidth  = 1920
	camHeight = 1080
)

// captureFrames captures frames from a camera at fixed intervals for a specified period of time and store them in the
// provided directory.
func captureFrames(dirName string, deviceID uint, duration, interval time.Duration) error {
	// validate params
	if duration.Nanoseconds() == 0 || interval.Nanoseconds() == 0 {
		return errors.New("duration and interval must be non-zero")
	}
	if duration <= interval {
		return errors.New("timelapse duration must be greater than the interval")
	}

	// determine final number of frames created upon successful timelapse capture
	frameCount, err := strconv.ParseInt(fmt.Sprintf("%d", duration/interval), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to determine frame count: %s", err)
	}
	// count 0th frame
	frameCount++
	log.Printf("  A total of %d frames will be captured over a period of %s.", frameCount, duration)
	log.Printf("  Estimated completion date is %s.", time.Now().Add(duration).Format(time.RFC1123))

	// open web cam
	cam, err := gocv.OpenVideoCapture(int(deviceID))
	if err != nil {
		return fmt.Errorf("failed to open video capture: %s", err)
	}
	defer cam.Close()

	// set camera resolution
	cam.Set(gocv.VideoCaptureFrameWidth, camWidth)
	cam.Set(gocv.VideoCaptureFrameHeight, camHeight)

	// create directory to store timelapse images
	log.Println("  Writing frames to \"" + dirName + "\" directory.")
	if err := os.MkdirAll(dirName, 0755); err != nil {
		return fmt.Errorf("failed to create directory for timelapse images: %s", err)
	}

	img := gocv.NewMat()
	defer img.Close()

	timeElapsed := time.Duration(0)
	for {
		// capture camera frame
		if ok := cam.Read(&img); !ok {
			log.Printf("cannot read from device")
			continue
		}
		if img.Empty() {
			continue
		}

		// write image frame to file to be stitched together later (allowing the restitching with different parameters
		// at a later date)
		imgFileName := dirName + strconv.FormatInt(time.Now().Unix(), 10) + "." + imageFormat
		if ok := gocv.IMWrite(imgFileName, img); !ok {
			return errors.New("failed to write image to file")
		}

		// break if last frame was captured
		if timeElapsed >= duration {
			break
		}
		time.Sleep(interval)
		timeElapsed += interval
	}

	return nil
}

// stitchFrames stitches all image frames from the provided directory into a video file.
func stitchFrames(dirName string, fps, frameStep uint) error {
	img := gocv.NewMat()
	defer img.Close()

	// get all image names in frames directory
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return fmt.Errorf("error collecting files in \"%s\" directory: %s", dirName, err)
	}

	// read single image to determine video dimensions
	if len(files) < 2 {
		return errors.New("two or more files are required to perform stitching")
	}
	img = gocv.IMRead(dirName+files[0].Name(), gocv.IMReadAnyColor)

	// iterate over all images and stitch them together
	videoFileName := "fps-" + strconv.FormatUint(uint64(fps), 10)
	videoFileName += "_step-" + strconv.FormatUint(uint64(frameStep), 10) + "_timelapse." + videoFormat
	writer, err := gocv.VideoWriterFile(dirName+"../"+videoFileName, "MJPG", float64(fps), img.Cols(), img.Rows(), true)
	if err != nil {
		return fmt.Errorf("error creating video writer: %s", err)
	}
	defer writer.Close()

	log.Println("  Timelapse video will be named \"" + videoFileName + "\".")

	// read each image file and write to video file writer
	count := 0
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), imageFormat) {
			continue
		}
		// skip frames which are not divisible by the step
		count++
		if count%int(frameStep) != 0 {
			continue
		}

		img = gocv.IMRead(dirName+f.Name(), gocv.IMReadAnyColor)
		if err := writer.Write(img); err != nil {
			return fmt.Errorf("error reading image file (%s) from \"%s\" directory: %s", f.Name(), dirName, err)
		}
	}

	return nil
}
