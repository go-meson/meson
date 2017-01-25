#!/bin/bash

rsync -av ../framework/src/api/version.h ./internal/binding/include/
rsync -av ../framework/src/api/meson.h ./internal/binding/include/
