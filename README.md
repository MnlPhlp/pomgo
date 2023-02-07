# pomgo
A simple and adjustable pmodoro timer for the commandline providing a progress bar and notifications.

## Installation
If you have go installed just run: `go install github.com/MnlPhlp/pomgo@latest`

Otherwise you can download one of the binaries from the releases section of this repo.

## Usage
Just running pomgo starts a default plan. To change this you can give the program any combination of different modes as a 'plan'.
This plan can also be saved in a file. `pomgo -h` gives you a help where different modes are explained

```
pomgo -h
Usage: pomTimer [OPTIONS] [PLAN]

Plan:
    list some tasks to form your plan (eg. wswlwsw)
    whitespaces are ignored so you can seperate tasks as you want
    or
    specify a file to read tasks from

Tasks:
    w[TIME]           work                 default time: 25min
    s[TIME]           take a short break   default time:  5min
    l[TIME]           take a long break    default time: 15min
    cTIME[TEXT]:      custom task          default time:  none

Options:
    -h        show this help
    -p        show the parsed plan and exit (includes -t)
    -t        show the time your plan will take and exit
    -n        no notifications
    -r N, -repeat N  repeat the plan N times
```
