# Pingo ![pass](https://github.com/andro404-MC/quigo-gui/actions/workflows/test.yml/badge.svg) ![GitHub License](https://img.shields.io/github/license/andro404-MC/quigo)

bspc rule -a Quigo state=floating

A dead simple app that let you know If you are connected . (with polybar support)

> [!WARNING]
> the notification mode is probably only usable with dunst

#### Why does this exist :

sooo there was a time when my internet was so good at disapearing when needed so I made this lil app to just tell me when to go back the moment its back.

and yea I just kept messing with it once a while and its public now. :)

> my internet still sucks so this thing is actually usefull.

## Requirement :

Nothing unless :

`libnotify` : for notification support.

`polybar` : duh, of course if you want to use it in polybar you will need polybar.

`go` : if you are going to build from source.

## Build :

> You need a to have `GOPATH` added to `PATH`

```
$ git clone https://github/untemi/pingo
$ cd pingo

// Run
$ go run .

// Install
$ go install .
```

## Usage :

To run :

```
$ pingo
```

Arguments :

- `-m`,`--mode` : Mode (term, termMin, ico, notify) (default "term")

- `--no-trail` : no trail (no replacing last line)

- `-n`,`--nonestop` : turn on noneStop (use '$ killall pingo' for stop)

- `-p`,`--polybar` : polybar colors

- `--recheck-delay` : delay between rechecks in Seconds (default 8)

- `--retry-delay` : delay between retrys in Seconds (default 1)
- `--timeout` : ping timeout in Milliseconds (default 200)

The default ping target is google.com but it can be changed using the environment variable `PINGOIP` :

```
$ PINGOIP=gnu.org pingo
```

## Todo :

None for now :)
