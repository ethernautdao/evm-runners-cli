# evm-runners-cli

Command line interface for evm-runners

## Installation

### Shell script

```
curl -L https://raw.githubusercontent.com/ethernautdao/evm-runners-cli/main/install.sh | bash
```

This will install the binary in `~/.evm-runners` and updates your shell configuration file (e.g., `~/.bashrc` or `~/.zshrc`).

After successfull installation, restart your terminal to use `evm-runners`.

### From source

```
make && make install
```

This will install the binary in `~/.evm-runners`

Note that if you want to run the evm-runners binary from any directory, you need to make sure that `${HOME}/.evm-runners` is added to your PATH environment variable. You can do this by adding the following line to your shell configuration file (e.g., `~/.bashrc` or `~/.zshrc`):

```
export PATH="${HOME}/.evm-runners:${PATH}"
```

To build from source you need to have [Go 1.20](https://go.dev/doc/install) installed.

## Commands

**Display help**

```
evm-runners -h
```

**Initialize evm-runners**

```
evm-runners init
```

This command clones the [evm-runners-levels](https://github.com/ethernautdao/evm-runners-levels) repository into the current directory and creates a .env file in `~/.evm-runners/`

**Authentication**

```
evm-runners auth <platform>
```

Authenticates the user. As of now only Discord authentication is available, e.g. `evm-runners auth discord`

**Start solving a level**

```
evm-runners start
```

Opens a list of levels to choose from. Alternatively, you can also start solving a level by providing the level name as an argument, e.g. `evm-runners start average`

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
evm-runners submit <level_name>
```

Optional flags:

- `--bytecode` or `-b`, to submit bytecode directly, e.g. `evm-runners submit Average -b 0xabcd`
- `--lang` or `-l`, to choose the language of the solution file when more than one solution file is present, e.g. `evm-runners submit Average -l sol`

**Display a list of all levels**

```
evm-runners levels
```

**Display the gas and codesize leaderboard of a level**

```
evm-runners leaderboard <level_name>
```

**Display info about evm-runners**

```
evm-runners about
```
