# termracer
Practise your typing skills from within your terminal. termracer is inspired by various online typing tutor websites.

The Goal is to type a given paragraph as fast and accurate as possible, termracer will calculate your typing speed with words per minute and accuracy % metrics. You can also view your progress by viewing the past race results.

For each race, You'll be presented a paragraph randomly picked from a predefined pool of paragraphs.

![](https://github.com/jan25/termracer/blob/master/assets/example.gif)

## Install

Without Go environment setup:
```
TODO
Upload binary to release
Add snippet to download and put it in $PATH
```

If you have Go environment setup locally:
```
# download and install
$ go get -u github.com/jan25/termracer

# run application
# if $GOPATH/bin is in $PATH
$ termracer

# OR
$ $GOPATH/bin/termracer
```

## Development
This application uses go modules. So, you could clone this repo under any
directory and build/test/run. As a helper we have Makefile in this repo, which will allow to build/test/run with single
command.
```
# Run available tests
$ make test

# Build and Run executable
$ make run

# Build and Run in debug mode
$ make debug

# Builds executable
$ make build
```

TODO Outdated:
The design/features are written in [NOTES.md](https://github.com/jan25/termracer/blob/master/NOTES.md).


TODO:
1. Update README
2. Update NOTES
3. Basic testing - clear data dirs and restart
4. Remove any dead code