#!/bin/bash

rsync -av ../framework/out/D/Meson.framework ./debug_dist
rsync -av ../framework/out/D/Meson\ Helper.app ./debug_dist
rsync -av ../framework/src/api/meson.h ./debug_dist/include/
#rsync -av ../framework/out/D/Meson\ Helper\ NP.app ./dist
#rsync -av ../framework/out/D/Meson\ Helper\ EH.app ./dist
rsync -av ../framework/out/R/Meson.framework ./dist
rsync -av ../framework/out/R/Meson\ Helper.app ./dist
rsync -av ../framework/src/api/meson.h ./dist/include/
