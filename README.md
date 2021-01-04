# Tii

On most GNU/Linux systems, when a command is not found, a message showing what can be run to install the command is printed. However, macOS does not have this.

This program supports a similar function with support for macOS (only for macOS, as of now). Instead of simply printing the command, Tii also offers to run it for you.

[comment]: <> ([![asciicast]&#40;https://asciinema.org/a/5eqCdcG6V8PC1nkekxbJ6gTXW.svg&#41;]&#40;https://asciinema.org/a/5eqCdcG6V8PC1nkekxbJ6gTXW?autoplay=1&&#41;)

![demo](https://cloud-5no1v5e34.vercel.app/0ezgif-6-80cd802cef12.gif)

The name Tii is an acronym for "Then Install It", what you will probably say when shown "command not found".

## Usage

```text
Usage: tii [--help/-h | <command>]
Examples:
   tii fish
   tii cowsay
   tii --help
```