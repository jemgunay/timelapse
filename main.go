package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	// parse flags
	var (
		cmd       = flag.String("cmd", "both", "The cmd can be either \"capture\", \"stitch\" or \"both\".")
		deviceID  = flag.Uint("cam_device_id", 0, "[capture] Recording camera device ID.")
		duration  = flag.Duration("duration", time.Hour, "[capture] Total duration of time to record a timelapse for.")
		interval  = flag.Duration("interval", time.Minute, "[capture] Time interval between the capture of each frame.")
		fps       = flag.Uint("fps", 25, "[stitch] Frame rate of the output timelapse video.")
		frameStep = flag.Uint("frame_step", 1, "[stitch] The step by which to iterate over frames to be stitched/excluded.")
		stitchDir = flag.String("stitch_dir", "", "[stitch] The path to the directory containing frames to stitch. Ignored if cmd is not set to \"stitch\".")
	)
	flag.Parse()

	if *frameStep == 0 {
		log.Println("frame_step must be greater than 0")
		return
	}
	if *fps == 0 {
		log.Println("fps must be greater than 0")
		return
	}

	// timestamp timelapse subdirectory
	dirName := time.Now().Format("Mon-02-Jan-2006_15-04-05")

	// validate cmd
	var performCapture, performStitch bool
	switch *cmd {
	case "both":
		performCapture = true
		performStitch = true
	case "capture":
		performCapture = true
	case "stitch":
		performStitch = true
		dirName = *stitchDir
	default:
		log.Println("cmd must be either \"both\", \"capture\" or \"stitch\"")
		return
	}

	// capture timelapse frames
	if performCapture {
		dirName = "./timelapses/" + dirName + "/frames/"
		log.Println("> Starting timelapse frame capture.")
		if err := captureFrames(dirName, *deviceID, *duration, *interval); err != nil {
			log.Printf("failed to capture frames: %s", err)
			return
		}
		log.Println("  Timelapse frame capture complete.")
	}

	// stitch captured frames from specified directory
	if performStitch {
		log.Println("> Generating timelapse video from frames.")
		if err := stitchFrames(dirName, *fps, *frameStep); err != nil {
			log.Printf("failed to stitch frames: %s", err)
			return
		}
		log.Println("  Timelapse video stitching complete.")
	}

	log.Println("> Timelapse processing complete.")
}
