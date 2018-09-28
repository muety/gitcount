# gitcount

[![Say thanks](https://img.shields.io/badge/SayThanks.io-%E2%98%BC-1EAEDB.svg)](https://saythanks.io/to/n1try)

A command-line tool to estimate the time spent on a git project, based on a very simple heuristic, inspired by [kimmobrunfeldt/git-hours](https://github.com/kimmobrunfeldt/git-hours).

### Assumptions: 
* Commits with a time difference less than 2 hours are considered to be in one coding session
* A multiple (x3) of the average time between commits in all sessions is added to the very first commit of every session

## Example
```
$ gitcount -dir .
Project root: /home/ferdinand/dev/mininote
mail@ferdinand-muetsch.de: 13.06 hours
exorcismo@gmail.com: 0.95 hours
noreply@github.com: 3.80 hours
btbtravis@gmail.com: 1.11 hours
kiantrue@gmail.com: 0.95 hours
fmuetsch@inovex.de: 0.00 hours
---------
Total: 19.86 hours
```

## Requirement
* Go to be installed

## How to use?
1. `go get github.com/n1try/gitcount`
2. `gitcount` or `gitcount -dir /some/project/path`

## License
MIT @ [Ferdinand Mütsch](https://ferdinand-muetsch.de)
