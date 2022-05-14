#!/bin/bash
gitlog=`TZ=UTC git log -1 --date=format-local:'%Y%m%d%H%M%S' --pretty=format:v0.0.0-%ad-%H`;echo ${gitlog:0:34}
