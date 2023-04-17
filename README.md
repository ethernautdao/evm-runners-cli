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
go build -o evm-runners main.go
```

## Formatting
```
go fmt /path/to/package
```

## Commands

Display help
```
./evm-runners -h
```

Initialize evm runners
```
./evm-runners init
```

Authentication (right now only Discord is available)
```
./evm-runners auth discord
```

Start a challenge
```
./evm-runners start --level <level_name>
```
e.g. `./evm-runners start Average`

Validate a challenge
```
./evm-runners validate <level_name>
``` 
Optional flag `--bytecode` or `-b`, to validate bytecode directly, e.g. `./evm-runners validate Average --bytecode 0xabcd`

Submit a solution
```
./evm-runners submit <level_name> --user_id <userid>
```
Optional flag `--bytecode` or `-b`, to submit bytecode directly, e.g. `./evm-runners submit Average -b 0xabcd`

Display a list of all levels
```
./evm-runners list
```

Display the gas and codesize leaderboard of a level
```
./evm-runners leaderboard <level_name>
```