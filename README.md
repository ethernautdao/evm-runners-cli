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
e.g. `./evmrunners start --level S01E01-Average`

Validate a challenge
```
./evmrunners validate --level <level_name>
``` 
Optional flag `--bytecode`, to validate bytecode directly, e.g. `./evmrunners validate --level S01E01-Average --bytecode 0xabcd`

Submit a solution
```
./evmrunners submit --level <level_name>
```
Optional flag `--bytecode`, to submit bytecode directly (WIP)

Display a list of all levels
```
./evmrunners list
```