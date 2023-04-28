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

## Installation

```
make && make install
```

This will install the binary in `~/.evm-runners`

Note that if you want to run the evm-runners binary from any directory, you need to make sure that `${HOME}/.evm-runners` is added to your PATH environment variable. You can do this by adding the following line to your shell configuration file (e.g., `~/.bashrc` or `~/.zshrc`):

```
export PATH="${HOME}/.evm-runners:${PATH}"
```

## Commands

**Display help**

```
evm-runners -h
```

**Initialize evm runners**
This command clones the [evm-runners-levels repo](https://github.com/ethernautdao/evm-runners-levels) into the current directory and creates a .env file in `~/.evm-runners/`

```
evm-runners init
```

**Authentication**
Authenticates the user. As of now only Discord authentication is available

```
evm-runners auth <platform>
```

**Start solving a level**

```
evm-runners start --level <level_name>
```

e.g. `evm-runners start Average`
Optional flags:

- `--lang` or `-l`, to directly choose the language of the solution file you want to work on, e.g. `evm-runners start Average -l sol`

**Validate a solution for a level**

```
evm-runners validate <level_name>
```

Optional flags:

- `--bytecode` or `-b`, to validate bytecode directly, e.g. `evm-runners validate Average --bytecode 0xabcd`
- `--lang` or `-l`, to choose the language of the solution file when more than one solution file is present, e.g. `evm-runners validate Average -l sol`

**Submit a solution**

```
./evm-runners submit <level_name>
```

Optional flags:

- `--bytecode` or `-b`, to submit bytecode directly, e.g. `evm-runners submit Average -b 0xabcd`
- `--lang` or `-l`, to choose the language of the solution file when more than one solution file is present, e.g. `evm-runners submit Average -l sol`

**Display a list of all levels**

```
evm-runners list
```

**Display the gas and codesize leaderboard of a level**

```
evm-runners leaderboard <level_name>
```
