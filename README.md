# termracer
Practise your typing skills from within your terminal. This app is inspired by various online typing tutor websites.

## Install

```
# download and install
$ go get -u github.com/jan25/termracer

# run application
# if $GOPATH/bin is in $PATH
$ termracer

# OR
$ $GOPATH/bin/termracer
```

> Current version of termracer can't generate paragraphs and pick an interesting paragraph for a race. We only have one default paragraph that is used for all races. In upcoming versions, termracer will be able to choose a random and interesting paragraphs for you.

## Development
This application uses go modules. So, you could clone this repo under any
directory and build/test/run. As a helper we have Makefile in this repo, which will allow to build/test/run with single
command.
```
# Builds executable
$ make build

# Runs available tests
$ make test

# Builds and runs executable
$ make run
```

The paragraph for a race is served by a backend server which is under `/server` directory. Server can be started indepedently using either docker cli or docker-compose.
```
# Build and run using Docker cli
$ docker build -t termracer-server
$ docker run --rm -d termracer-server

# OR

# Use docker compose to bootup server
$ docker-compose up -d
# To stop server
$ docker-compose down
```

The design/features are written in [NOTES.md](https://github.com/jan25/termracer/blob/master/NOTES.md).