# termracer
Practise your typing skills from within your terminal. This app is inspired by various online typing tutor websites.

This application is very much in development. I'll update this README according to when this app is ready to be used.

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