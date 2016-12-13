#!/bin/bash

rsync -av ../framework/out/D/Meson.framework ./dist
rsync -av ../framework/out/D/Meson\ Helper.app ./dist
rsync -av ../framework/src/api/meson.h ./dist/include/
