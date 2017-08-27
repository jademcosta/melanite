# Melanite
A server that helps converting images on the fly.

## Usage

Melanite is a very fast image converter and resizer. It is great to reduce bandwith of images for your mobile apps and website by converting images to WEBP (just an example). If you're sending PNG images, Melanite is a must.

On the [test](https://github.com/jademcosta/melanite/tree/master/test/images) folder you have examples of images that equal, but in different formats, and can have a taste of the difference in disk size between each format.

If your site/app has high traffic, you'd be better use melanite behind a CDN. With this, you can safely run Melanite on a simple machine, and enjoy all the speed and savings.

More about the capabilities below.

### Resizing images

Assuming you have Melanite running on localhost:8080, and the image_source of it is `https://www.google.com.br`, to reduce the size of an image you can access the address `localhost:8080/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png?res=60x0`. This will produce an image with width of 60 pixels, and proportional height.

If you want to enlarge the image, all you have to do is provide a higher size to the `res` query params, like `localhost:8080/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png?res=1000x0` .

You can mix resizing and converting on the same image query, like this:
localhost:8080/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png?res=60x0&out=webp

### Converting images

Assuming you have Melanite running on localhost:8080, and the image_source of it is `https://www.google.com.br`, to convert the format of an image you have go to address `localhost:8080/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png?out=webp` (it can be on your browser). This will produce a WEBP image.

Have in mind that when you convert an image that has transparency (WEBP and PNG) to JPG, Melanite will fill the transparent pixels with black, as JPG does not supports transparency.

Currently, Melanite knows how to convert to JPG, PNG and WEBP. You can mix resizing an converting on the same image query, like this:
localhost:8080/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png?res=60x0&out=webp

### Benchmarks
Melanite uses [libvips](https://github.com/jcupitt/libvips), and leverages its speed and low memory usage. You can check a libvips benchmark on their repository, [here](https://github.com/jcupitt/libvips/wiki/Speed-and-memory-use).

## Deploying it

Run `go build melanite`. Get the binary generated and puts it on a folder where you have a config file. You can get the `melanite_example_config.yaml` in this repository, *and replace* the `image_source` property with your images' server address (WITHOUT THE TRAILING SLASH. The address should be http://example.com and not http://example.com/). After that, you can visit `http://your_server_ip:8080/some_image_path.extension` and check if it works. If it's everything ok, check on the examples section what can you do with your images.

### Config
The config file can be given using the -c param, when running melanite, like `melanite -c /etc/melanite/melanite_config.yml`. A config file example can be found at the root of this repo (please note that in order to use it you'll need to replace some values!). Below you'll find the valid entries of config file.

* *image_source* [required]: The url of the image server. Should not end with a slash. The url http://example.com is valid, while http://example.com/ is invalid.
* *port* [not-required]: The port where melanite will run. If no value is given, it will default to port 8080.

### Android support for WEBP
Android doesn't suppport WEBP on all its versions. If you are using Melanite as your images proxy, don't convert to WEBP on older versions of Android. Read more about it on these links:
* https://developer.android.com/studio/write/convert-webp.html
* https://developer.android.com/guide/topics/media/media-formats.html

### iOS support WEBP
iOS seems to be starting to support WEBP, but Melanite was not tested against it yet. If you tested it and it worked on iOS, please send a PR editing this part. I'm basing this assumption on this link:
* https://stackoverflow.com/questions/8672393/webp-image-format-on-ios


## Development


### Libs used

We are relying heavly on the excelent [bimg lib](https://github.com/h2non/bimg) which, in turns, rely on [libvips](https://github.com/jcupitt/libvips). Libvips is a blazing fast image lib, and in order to build this project you'll need to have libvips installed on your machine (on Ubuntu it is just a `sudo apt install libvips-dev`).

### Testing

Run all tests with the command `go test $(go list ./... | grep -v /vendor/)`

### Dependency management
This project is using [govendor](https://github.com/kardianos/govendor) to manage dependencies.
To add a new dependency, just run `govendor fetch PATH_TO_DEPENDENCY`.

### Future work
* [Security] Add security header to config: If the header does not match with the request, it gets denied.
* [Perf] Pass the image as a pointer inside the controller. This will improve GC and memory.
* [UX-Perf] Add ETags on each image, and allow the server to respond 304 - Nothing changed.
* [UX] Add optional Prometheus endpoint to monitoring.
* [Monit] Add middleware that inserts a X-Request-Id header if not existent. This will help debugging and loggers.
* [UX] Allow the level setting of the logger through config file. This is to disable log for those who don't care about it.
* [UX-Perf] Allows to set a max-age cache number on config file.
* [Perf] If more than one request arrives for the same image, adds a chanel, so that only 1 GET is done to the upstream.
* [Feature] Adds to config the option to store generated files on disk.
* [Feature] Adds to config the option to store generated files on memory (LRU cache with limited size).
* [Feature] Adds to config the option to store generated files on external repo. But think twice before doing it on S3 (or anything like), because this app (thanks libvips) is so fast that it might take longer to download the image that it takes to get it from disk and do all the necessary processing.
* [Feature] Allows the user to crop images. Allows the option "crop" when resizing (downsizing) images.  
* [Feature] Allows color transformation to users (gray, sepia, etc).
* [Feature] Allows smartcrop images. Based on face detection and feature detection.
* [Code] Pass the resizer and converter with dependency injection to the controller. This will allow us to test it (controller) more.
* [UX] Return 400 when the user tries an invalid resize param.
* [UX] Return errors on a header, to help debugging.
* [Docs] Add a benchmark of melanite running on localhost (conversion speed benchmark).


## Thanks
A big thank you for all the people involved into [libvips](https://github.com/jcupitt/libvips), and also the [bimg](https://github.com/h2non/bimg).
