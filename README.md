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
  database    database management functions
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

### Shell completion

`asoai` is built using [cobra](https://cobra.dev/). This allows adding auto-completion for your favorite shell:

```sh
$ ./asoai completion fish | source
$ ./asoai completion fish | source
chat                                             (interact with chatgpt)  database  (database management functions)  models       (list models)
completion  (Generate the autocompletion script for the specified shell)  help             (Help about any command)  session  (handle sessions)
```

### REPL

It includes a basic REPL for easier conversations:

```sh
$ ./asoai chat --new-session --name "my-kubernetes-talk" --repl --stream --system-prompt "You are an LLM that only talks about kubernetes"
user> hello chatgpt
assistant> Hello! How can I help you with Kubernetes today?
user> what is a job
assistant> In Kubernetes, a Job is a resource object that creates one or more Pods to run a particular task to completion. Once the task is completed, the Job itself is considered complete. Jobs are used for batch processing, running tasks that are supposed to run once and then stop, such as data processing, backups, or utility tasks. Jobs can be used to run parallel tasks, but each task is expected to run to completion before the Job is considered done.
user>
```

### Using sessions

Create a session giving a model and a system prompt:

```sh
$ ./asoai --db-path ./data.db session create --name testaroo --model gpt-4o --system-prompt "You are a GPT4 model that only respond with an hello world golang program without anything but code"
testaroo

$ ./asoai --db-path ./data.db session set-current testaroo

$ ./asoai --db-path ./data.db chat "print me something"
package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}
```

```sh
$ ./asoai --db-path ./data.db session dump
Current session: testaroo
Model: gpt-4o

system> You are a GPT4 model that only respond with an hello world golang program without anything but code
user> print me something
assistant> 
package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}
```

Have fun!

