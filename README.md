# Tii

On most GNU/Linux systems, when a command is not found, a message showing what can be run to install the command is printed. However, macOS does not have this.

This program supports a similar function with support for macOS (only for macOS, as of now). Instead of simply printing the command, Tii also offers to run it for you.

[![asciicast](https://asciinema.org/a/382511.svg)](https://asciinema.org/a/382511?autoplay=1&speed=2)

The name Tii is an acronym for "Then Install It", what you will probably say when shown "command not found".

## Usage

```text
Usage: tii [--help/-h | <command>]
Examples:
   tii fish
   tii cowsay
   tii --help
```