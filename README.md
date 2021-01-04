# Tii

On most GNU/Linux systems, when a command is not found, a message showing what can be run to install the command is printed. However, macOS does not have this. 

This program supports a similar function with support for macOS (only for macOS, as of now). Instead of simply printing the command, Tii also offers to run it for you.

[comment]: <> ([![asciicast]&#40;https://asciinema.org/a/5eqCdcG6V8PC1nkekxbJ6gTXW.svg&#41;]&#40;https://asciinema.org/a/5eqCdcG6V8PC1nkekxbJ6gTXW?autoplay=1&&#41;)

![demo](https://cloud-5no1v5e34.vercel.app/0ezgif-6-80cd802cef12.gif)

# Usage

```text
ishan@mac [~] |>  tii --help
Tii - Directly install command when not found

On most GNU/Linux systems, when a command is not found, a message showing what
can be run to install the command is printed. However, macOS does not
have this. This program supports a similar function with support for macOS
(only for macOS, as of now). Instead of simply printing the command, Tii also
offers to run it for you.

Usage: tii [--help/-h | <command>]
Examples:
   tii fish
   tii cowsay
   tii --help

If Tii was installed correctly, using commands which are not found will
automatically trigger it. The name Tii is an acronym for "Then Install It".
```