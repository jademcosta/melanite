# Melanite
A server that helps converting images. It will, in the future, cache and resize images.

## Dependency management
This project is using [govendor](https://github.com/kardianos/govendor) to manage dependencies.
To add a new dependency, just run `govendor fetch PATH_TO_DEPENDENCY`.

## Libs used

We are relying heavly on the excelent [bimg lib](https://github.com/h2non/bimg) which, in turns, rely on [libvips](https://github.com/jcupitt/libvips). Libvips is a blazing fast image lib, and in order to build this project you'll need to have libvips installed on your machine (on Ubuntu it is just a `sudo apt install libvips-dev`).

### Android support for WEBP
WIP
https://developer.android.com/studio/write/convert-webp.html
https://developer.android.com/guide/topics/media/media-formats.html

### iOS support WEBP
WIP
https://stackoverflow.com/questions/8672393/webp-image-format-on-ios
