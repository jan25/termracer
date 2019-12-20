# termracer
Practise your typing skills from within your terminal. termracer is inspired by various online typing tutor websites.

The Goal is to type a given paragraph as fast and accurate as possible, termracer will calculate your typing speed with words per minute and accuracy % metrics. You can also view your progress by viewing the past race results.

For each race, You'll be presented a paragraph randomly picked from a predefined pool of paragraphs.

![](https://github.com/jan25/termracer/blob/master/assets/example.gif)

## Install

```
# download and install latest release
$ go build -o $GOPATH/bin/termracer github.com/jan25/termracer/cmd

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

The design/features are written in [NOTES.md](https://github.com/jan25/termracer/blob/master/NOTES.md).


TODO:
4. Remove any dead code
5. Upload new demo asset