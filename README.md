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
Optional flags:
- `--lang` or `-l`, to directly choose the language of the solution file you want to work on, e.g. `./evm-runners start Average -l sol`

Validate a challenge
```
./evm-runners validate <level_name>
``` 
Optional flags:
- `--bytecode` or `-b`, to validate bytecode directly, e.g. `./evm-runners validate Average --bytecode 0xabcd`
- `--lang` or `-l`, to choose the language of the solution file when more than one solution file is present, e.g. `./evm-runners validate Average -l sol`

Submit a solution
```
./evm-runners submit <level_name> --user_id <userid>
```
Optional flags:
- `--bytecode` or `-b`, to submit bytecode directly, e.g. `./evm-runners submit Average -b 0xabcd`
- `--lang` or `-l`, to choose the language of the solution file when more than one solution file is present, e.g. `./evm-runners submit Average -l sol`

Display a list of all levels
```
./evm-runners list
```

Display the gas and codesize leaderboard of a level
```
./evm-runners leaderboard <level_name>
```