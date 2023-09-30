# evm-runners-cli

A command line interface for evm-runners, a terminal-based game with EVM based levels.

## Installation

### Installation script

```
curl -L get.evmr.sh | bash
```

This command will install the binary in `~/.evm-runners/bin` and updates PATH in your shell configuration file (e.g. `~/.bashrc`, `~/.zshrc`, ...).

After successfull installation you can initialize evm-runners with `evmr init`, or alternatively, `evm-runners init`.

To update the CLI to the latest version, run `evmrup`.

**Alternatively**

Install from source by running:

```
make && make install
```

This will install the binary in `~/.evm-runners/bin`

Note: If you wish to run the evm-runners binary from any directory, ensure that `${HOME}/.evm-runners` is added to your PATH environment variable. You can do this by adding the following line to your shell configuration file:

```
export PATH="${PATH}:{HOME}/.evm-runners/bin"
```

Make sure you have [Go 1.20](https://go.dev/doc/install) or a later version installed to compile the source code.

## How to play

After successful installation, run `evmr init` to initialize evm-runners. This will clone the [evm-runners-levels](https://github.com/ethernautdao/evm-runners-levels) repository into the current directory. You can then start solving a level by running `evmr start`.

To validate a solution, run `evmr validate <level>`. If it is valid, you can submit it by running `evmr submit <level>`. Before submitting a solution you have to authenticate your account by running `evmr auth discord`.

## Available commands

Note: You can invoke all commands with `evm-runners <cmd>` as well.

**Display info about evm-runners**

```
evmr about
```

**Authentication**

```
evmr auth <platform>
```

Authenticates your account. As of now only Discord authentication is available: `evmr auth discord`.
Additionally, `evmr auth wallet` (or `evmr auth address`) links an Ethereum address to your account, allowing you to submit solutions from the website.

**Display help**

```
evmr -h
```

Use `evmr <cmd> -h` to display help for a specific command.

**Initialize evm-runners**

```
evmr init
```

This command clones the [evm-runners-levels](https://github.com/ethernautdao/evm-runners-levels) repository into the current directory and updates the .env file in `~/.evm-runners/`

**Show the leaderboard of a level**

```
evmr leaderboard <level>
```

**Display a list of all levels**

```
evmr levels
```

**Start solving a level**

```
evmr start
```

Opens a list of levels to choose from. Alternatively, you can also start solving a level by providing the level name as an argument, e.g. `evmr start average`

Optional flags:

- `--lang` or `-l`, to directly choose the language of the solution file you want to work on, e.g. `evmr start average -l sol`

**Submit a solution**

```
evmr submit <level>
```

Optional flags:

- `--bytecode` or `-b`, to submit bytecode directly, e.g. `evmr submit average -b 0xabcd`
- `--lang` or `-l`, to choose the language of the solution file when more than one solution file is present, e.g. `evmr submit average -l sol`

**Update levels directory**

```
evmr update
```

This command runs `git pull` inside the levels directory, updating the levels to the latest version.

**Validate a solution for a level**

```
evmr validate <level>
```

Optional flags:

- `--bytecode` or `-b`, to validate bytecode directly, e.g. `evmr validate average --bytecode 0xabcd`
- `--lang` or `-l`, to choose the language of the solution file when more than one solution file is present, e.g. `evmr validate average -l sol`

**Show the current version of evm-runners**

```
evmr version
```
