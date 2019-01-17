#!/bin/bash

nohup ./timelapse -duration=12h -interval=1m &> timelapse.log &
