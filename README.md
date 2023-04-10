# evm-runners-cli

Command line interface for evm-runners

## Prerequisites

- [Go 1.20](https://go.dev/doc/install)
- [git](https://github.com/git-guides/install-git)

## Building

```
make
```
or
```
go build -o evmrunners main.go
```

## Formatting
```
go fmt /path/to/package
```

## Commands

Display help
```
./evmrunners -h
```

Initialize evm runners
```
./evmrunners init
```

Start a challenge
```
./evmrunners start --level <level_name>
```
e.g. `./evmrunners start Average`

Validate a challenge
```
./evmrunners validate <level_name>
``` 
Optional flag `--bytecode`, to validate bytecode directly, e.g. `./evmrunners validate Average --bytecode 0xabcd`

Submit a solution
```
./evmrunners submit --level <level_name> --user_id <userid>
```
Optional flag `--bytecode`, to submit bytecode directly

Display a list of all levels (WIP)
```
./evmrunners list
```

Display the gas and codesize leaderboard of a level
```
./evmrunners leaderboard <level_name>
```