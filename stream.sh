#!/bin/bash

ffmpeg -draw_mouse 0 -async 1 -vsync 1 -f x11grab -s "${WIDTH}x${HEIGHT}" -framerate 15 -i :0.0 -c:v libx264 -g 30 -pix_fmt yuv420p -s "${WIDTH}x${HEIGHT}" -threads 0 -f flv "rtmp://live.twitch.tv/app/$STREAM_KEY"
