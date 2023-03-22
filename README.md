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
