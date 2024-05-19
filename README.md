# ASOAI: Another Stupid OpenAI client

Simple but efficient OpenAI client.

It features:
- basic cli chat, with optional streaming mode
- a REPL/interactive mode
- configurable sessions management
- history database

## Building & running

```sh
$ go build
$ ./asoai help
asoai is another stupid OpenAI client

Usage:
  asoai [flags]
  asoai [command]

Available Commands:
  chat        interact with chatgpt
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  models      list models
  session     handle sessions

Flags:
  -h, --help   help for asoai

Use "asoai [command] --help" for more information about a command.
```

## Usage

### Basic chat

Set the `OPENAI_API_KEY` env var with your own api key. Then, launch a chat using `asoai chat <your input>`:

```sh
$ ./asoai "hello gpt!"
Hello! How can I assist you today?
```

### Using sessions

Create a session giving a model and a system prompt:

```sh
$ ./asoai --db-path ./data.db session create --name testaroo --model gpt-4o --system-prompt "You're a GPT4 model that only respond with an hello world golang program without anythi
ng but code"
testaroo

$ ./asoai --db-path ./data.db session set-current testaroo

$ ./asoai --db-path ./data.db chat "print me something"
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}

$ ./asoai --db-path ./data.db session dump
Current session: testaroo
Model: gpt-4o

system> You're a GPT4 model that only respond with an hello world golang program without anything but code
user> print me something
assistant> ```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}
```


Have fun!

