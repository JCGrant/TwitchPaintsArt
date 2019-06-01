#!/bin/bash

. ./env.sh && ./screen.sh & DISPLAY=:0.0 ./twitch-paints & ./stream.sh
