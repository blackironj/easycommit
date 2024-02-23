# easycommit

easycommit is a command-line interface (CLI) tool designed to streamline the process of writing Git commit messages. It uses an [Ollama](https://github.com/ollama/ollama) as backend.

🚧 I'm currently trying to find appropriate and accurate prompts. and I think `Mistral` is the best fit on my CLI now.

## Installation

### Binary

- Download a binary file at the [release page](https://github.com/blackironj/easycommit/releases)

### Manual

- build it from code

``` bash
git clone https://github.com/blackironj/easycommit.git
cd easycommit
go build .
```

## Usage

- You have to run the Ollama before running *easycommit*.
- After installing Ollama. you should pull models what you want

    ``` bash
    ollama pull mistral
    ```

- Please see the [Ollama page](https://github.com/ollama/ollama)

``` bash
simple cli tool for generating commit using Ollama

Usage:
  easycommit [flags]

Flags:
  -d, --debug               for debugging log
  -e, --endpoint string     ollama host url (default "http://127.0.0.1:11434")
  -h, --help                help for easycommit
  -m, --model string        llama model (default "mistral")
  -n, --num-predict int     num predict (default 200)
  -t, --temperature float   temperature (default 0.7)
```
