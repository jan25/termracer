# termracer
Practise your typing skills from within your terminal. termracer is inspired by various online typing tutor websites.

Your goal is to type a given paragraph as fast and accurate as possible, termracer will calculate your typing speed with words per minute and accuracy % metrics. You can also view your progress by viewing the past race results.

![](https://github.com/jan25/termracer/blob/master/example.gif)

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

# Run in debug mode
$ make debug
```

## Note: The current master is in one single main package with global shared variables. So in-order to seperate concerns the project is in rewrite period, so at end of it we'll have components seperated into nicer modules/packages. Expect the finished rewrite by milestone 0.2.0-alpha scheduled on 01-12-2019 CET

The design/features are written in [NOTES.md](https://github.com/jan25/termracer/blob/master/NOTES.md).