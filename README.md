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
export PATH="${PATH}:/{HOME}/.evm-runners/bin"
```

Make sure you have [Go 1.20](https://go.dev/doc/install) or a later version installed to compile the source code.

## Commands

Note: You can invoke all commands with `evm-runners <cmd>` as well.

**Display help**

```
evmr -h
```

**Initialize evm-runners**

```
evmr init
```

This command clones the [evm-runners-levels](https://github.com/ethernautdao/evm-runners-levels) repository into the current directory and updates the .env file in `~/.evm-runners/`

**Start solving a level**

```
evmr start
```

Opens a list of levels to choose from. Alternatively, you can also start solving a level by providing the level name as an argument, e.g. `evmr start average`

Optional flags:

- `--lang` or `-l`, to directly choose the language of the solution file you want to work on, e.g. `evmr start average -l sol`

**Validate a solution for a level**

```
evmr validate <level_name>
```

Optional flags:

- `--bytecode` or `-b`, to validate bytecode directly, e.g. `evmr validate average --bytecode 0xabcd`
- `--lang` or `-l`, to choose the language of the solution file when more than one solution file is present, e.g. `evmr validate average -l sol`

**Submit a solution**

```
evmr submit <level_name>
```

Optional flags:

- `--bytecode` or `-b`, to submit bytecode directly, e.g. `evmr submit average -b 0xabcd`
- `--lang` or `-l`, to choose the language of the solution file when more than one solution file is present, e.g. `evmr submit average -l sol`

**Authentication**

```
evmr auth <platform>
```

Authenticates the user. As of now only Discord authentication is available, e.g. `evmr auth discord`.
Additionally, `evmr auth wallet` (or `evmr auth address`) links an Ethereum address to the user, enabling the user to submit solutions from the website.

**Display a list of all levels**

```
evmr levels
```

**Display the gas and codesize leaderboard of a level**

```
evmr leaderboard <level_name>
```

**Display info about evm-runners**

```
evmr about
```

**Show the current version of evm-runners**

```
evmr version
```
