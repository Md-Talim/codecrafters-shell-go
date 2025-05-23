[![progress-banner](https://backend.codecrafters.io/progress/shell/3296d15e-9c84-4c21-9af8-f8ab677604b7)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

# Build Your Own Shell (Go Edition)

This project is a custom implementation of a POSIX-compliant shell built in Go. It interprets shell commands, runs external programs, supports built-in commands like `cd`, `pwd`, `echo`, `history`, and handles features like I/O redirection, multi-stage pipelines, command autocompletion, and command history navigation.

This project is part of the ["Build Your Own Shell" Challenge](https://app.codecrafters.io/courses/shell/overview) by Codecrafters. It's an excellent way to learn about the internal workings of shells and practice Go programming concepts.

## 📑 Table of Contents

- [Build Your Own Shell (Go Edition)](#build-your-own-shell-go-edition)
  - [📑 Table of Contents](#-table-of-contents)
  - [📌 What This Project Does](#-what-this-project-does)
  - [⚙️ How to Set Up and Run](#️-how-to-set-up-and-run)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
    - [Building](#building)
    - [Running](#running)
    - [Usage](#usage)
  - [✨ Key Features](#-key-features)
  - [🔍 How It Works Internally](#-how-it-works-internally)
  - [🧰 Technologies Used](#-technologies-used)
  - [📁 Folder and File Structure](#-folder-and-file-structure)
  - [💡 Challenges \& Lessons Learned](#-challenges--lessons-learned)
    - [Challenges](#challenges)
    - [Lessons Learned](#lessons-learned)
  - [🛠️ Why I Built This Project](#️-why-i-built-this-project)
  - [🚀 Future Features](#-future-features)

## 📌 What This Project Does

This project implements a shell in Go that:

-   Accepts user input via a REPL (Read-Eval-Print Loop).
-   Parses complex command lines, including quoted arguments and escape sequences.
-   Executes built-in commands like `cd`, `pwd`, `echo`, `exit`, `type`, and `history`.
-   Runs external programs by searching for executables in the system's `PATH`.
-   Supports multi-stage command pipelines (e.g., `ls | grep .go | wc -l`).
-   Handles input/output redirection (e.g., `>`, `>>`, `2>`).
-   Provides autocompletion for commands (both built-in and external) and file paths.
-   Allows navigation and recall of command history using arrow keys.
-   Manages terminal I/O in raw mode for advanced input handling.

The purpose of this project is to gain a deep understanding of shell mechanics and to build a non-trivial system using Go.

## ⚙️ How to Set Up and Run

### Prerequisites

-   **Go**: Version 1.24 or later
-   **Linux Environment**: The shell utilizes POSIX features and terminal control mechanisms best suited for a Linux-like environment.

### Installation

1.  Clone the repository:
    ```bash
    git clone https://github.com/md-talim/codecrafters-shell-go.git
    cd codecrafters-shell-go
    ```

### Building

The project is built using the standard Go toolchain. The `your_program.sh` script compiles the application:
```bash
go build -o /tmp/codecrafters-build-shell-go app/*.go
```
You can run this command manually in the project root if needed.

### Running

Execute the shell using the provided script:
```bash
./your_program.sh
```

### Usage

-   You'll be greeted with a `$` prompt.
-   Enter commands like `pwd`, `echo Hello World`, `ls -l`, or `cat file.txt | grep keyword`.
-   Use `Ctrl+D` or the `exit` command to terminate the shell.
-   Press `Tab` for command autocompletion.
-   Use Up/Down arrow keys to navigate through command history.
-   Type `history` to see a list of previous commands, or `history <n>` to see the last `n` commands.

## ✨ Key Features

-   **Built-in Commands**:
    -   `cd`: Change current directory.
    -   `pwd`: Print working directory.
    -   `echo`: Display a line of text.
    -   `exit`: Terminate the shell.
    -   `type`: Display information about command type (builtin or external).
    -   `history [n]`: Display command history, optionally limited to the last `n` entries.
-   **External Command Execution**:
    -   Locates and runs external programs using the system `PATH`.
    -   Utilizes Go's `os/exec` package for process management.
-   **Pipeline Support**:
    -   Allows chaining multiple commands, where the output of one command becomes the input of the next (e.g., `cmd1 | cmd2 | cmd3`).
    -   Manages inter-process communication using OS pipes.
-   **I/O Redirection**:
    -   Redirects standard output (`>`), appends standard output (`>>`).
    -   Redirects standard error (`2>`).
-   **Autocompletion**:
    -   Press `Tab` to autocomplete command names (built-ins and executables from `PATH`).
    -   Suggests multiple completions if ambiguous.
-   **Command History Navigation**:
    -   Recall previous commands using the Up arrow key.
    -   Navigate to newer recalled commands or an empty line using the Down arrow key.
-   **Error Handling**:
    -   Provides informative error messages for issues like command not found, incorrect arguments, or file permission errors.
    -   Designed to prevent crashes from unexpected input or runtime issues.
-   **Raw Terminal Mode**:
    -   Captures input character-by-character for features like autocompletion, history navigation, and immediate feedback, without waiting for Enter.
    -   Uses `github.com/pkg/term/termios` and `golang.org/x/sys/unix` for low-level terminal control.

## 🔍 How It Works Internally

The shell is structured into several key components:

| Component                               | Role                                                                                                            |
| :-------------------------------------- | :-------------------------------------------------------------------------------------------------------------- |
| `app/main.go`                           | Entry point, contains the main REPL loop.                                                                       |
| `app/console.go`                        | Handles raw terminal input, character processing, prompt display, arrow key navigation for history.             |
| `app/autocomplete.go`                   | Implements tab-completion logic for commands and files.                                                         |
| `internal/parser/parser.go`             | Parses the raw input string into commands, arguments, and redirection configurations.                           |
| `internal/executor/command_executor.go` | Orchestrates command execution, deciding between single commands and pipelines.                                 |
| `internal/executor/pipeline_runner.go`  | Manages the setup and execution of multi-stage command pipelines, including pipe creation and I/O.              |
| `internal/executor/builtins.go`         | Defines and implements built-in shell commands, stores command history, provides history access.                |
| `internal/executor/helpers.go`          | Contains helper functions for the executor, including adding commands to history and parsing history arguments. |
| `internal/shellio/shellio.go`           | Manages I/O streams, including handling **file** redirections.                                                  |
| `internal/utils/path.go`                | Provides utility functions, such as searching for executables in the system `PATH`.                             |

## 🧰 Technologies Used

-   **Go**: The primary programming language.
-   **Go Standard Library**: Extensive use of packages like `os`, `os/exec`, `fmt`, `strings`.
-   **`github.com/pkg/term/termios` & `golang.org/x/sys/unix`**: For low-level terminal control (raw mode, non-canonical input).
-   **POSIX Standards**: Adherence to POSIX conventions for shell behavior and command execution.

## 📁 Folder and File Structure

```
├── app/                      # Main application package (REPL, console, autocomplete)
│   ├── main.go
│   ├── console.go
│   └── autocomplete.go
├── internal/                 # Internal packages not intended for external use
│   ├── executor/             # Command execution logic and built-ins
│   │   ├── builtins.go         # Built-in commands, history storage and access
│   │   ├── command_executor.go # Main execution dispatcher
│   │   ├── helpers.go          # Executor helpers, history utilities
│   │   └── pipeline_runner.go  # Pipeline execution logic
│   ├── parser/               # Input parsing logic
│   │   ├── parser.go
│   │   └── utils.go
│   ├── shellio/              # I/O stream management and redirection
│   │   └── shellio.go
│   └── utils/                # General utility functions
│       └── path.go
├── go.mod                    # Go module definition
├── your_program.sh           # Script to build and run the shell locally
└── README.md                 # This file
```

**Key files:**
-   `app/main.go`: The entry point and main loop of the shell.
-   `app/console.go`: Handles raw terminal input, including arrow key processing for history.
-   `internal/parser/parser.go`: Handles the parsing of user input.
-   `internal/executor/command_executor.go`: Manages the execution of parsed commands.
-   `internal/executor/pipeline_runner.go`: Specifically handles the execution of command pipelines.
-   `internal/executor/builtins.go`: Contains implementations for built-in commands and history access functions.
-   `internal/executor/helpers.go`: Includes functions for adding commands to history and processing `history` command arguments.

## 💡 Challenges & Lessons Learned

### Challenges

1.  **Command Parsing**:
    -   Correctly handling various quoting mechanisms (single, double quotes), escape characters, and tokenizing input into distinct arguments and operators.
2.  **Pipeline Implementation**:
    -   Managing multiple OS pipes for inter-process communication between commands in a pipeline.
    -   Ensuring correct setup of `stdin` and `stdout` for each command in the chain.
    -   Properly closing file descriptors to avoid deadlocks or resource leaks.
3.  **I/O Redirection**:
    -   Integrating redirection with single commands and pipelines, ensuring correct file opening, and stream management.
4.  **Raw Terminal Mode, Autocompletion & History Navigation**:
    -   Interacting with low-level terminal settings to capture individual keystrokes (including escape sequences for arrow keys).
    -   Implementing a responsive and context-aware autocompletion system.
    -   Managing state for history navigation (`historyNavigationIndex`) and correctly updating the display.
5.  **Process Management**:
    -   Correctly launching, managing, and waiting for external processes using `os/exec`.

### Lessons Learned

-   Improved proficiency in Go, process management, and low-level I/O.
-   Gained a much deeper understanding of how operating system shells function internally.
-   Learned the intricacies of POSIX terminal control, escape sequence parsing, and inter-process communication.
-   Developed skills in designing and implementing modular systems with clear separation of concerns, especially in managing state for interactive features.

## 🛠️ Why I Built This Project

I built this project to:

-   Deepen my understanding of system-level programming concepts.
-   Explore the mechanics of command-line interfaces and shell environments.
-   Practice Go for building a complex, interactive application.
-   Challenge myself by solving a real-world problem from the ground up.
-   Learn how fundamental tools like Bash or Zsh operate "under the hood."

This project has been a significant learning experience, enhancing my skills as a software engineer.

## 🚀 Future Features

-   **Job Control**: Add support for backgrounding processes (`&`) and managing jobs (`fg`, `bg`, `jobs`).
-   **Shell Variables & Expansion**: Introduce support for setting and using shell variables.
-   **Globbing**: Implement wildcard expansion for filenames (e.g., `ls *.txt`).
-   **Scripting**: Basic shell script execution capabilities.
-   **More Advanced Editing**: Cursor movement within the line (left/right arrows), insert/delete characters.
