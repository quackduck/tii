# Tii

On most GNU/Linux systems, when a command is not found, a message showing what can be run to install the command is printed. However, macOS does not have this.

This program adds a similar function for macOS (only for macOS, as of now). Instead of simply printing the command, Tii also offers to run it for you.

[comment]: <> ([![asciicast]&#40;https://asciinema.org/a/382511.svg&#41;]&#40;https://asciinema.org/a/382511?autoplay=1&speed=2&#41;)

<a href="https://asciinema.org/a/382511?autoplay=1&speed=2" target="_blank">
<img src="https://asciinema.org/a/382511.svg" alt="demo">
</a>

The name Tii is an acronym for "Then Install It", which is what you will probably say when shown "Command not found".

## Installing
As of now, only macOS is supported
```shell
brew install quackduck/tap/tii
```

## Usage

Tii will be automatically triggered if a command is not found and so you usually do not need to directly interact with it. 

The Tii binary has the following usage:

```text
Usage: tii [--help/-h | --version/-v | <command>]

Examples:
   tii fish
   tii cowsay
   tii --help
```

## Environment
These are the environment variables that can affect Tii:

* If `TII_DISABLE_INTERACTIVE` is set to "true", Tii will
    disable interactive output (prompting for confirmation) and not install
    any packages.

* If `TII_AUTO_INSTALL_EXACT_MATCHES` is set to "true", Tii will
    automatically install exact matches without prompting for confirmation. **This variable overrides `TII_DISABLE_INTERACTIVE`.**
  
## Uninstalling
If you have issues with Tii, head over to [issues](https://github.com/quackduck/tii/issues).

You can uninstall with:
```shell
brew uninstall tii
```

Here's a list of all the files Tii uses:
```text
/usr/local/bin/tii
/usr/local/share/fish/vendor_functions.d/tii_on_command_not_found.fish
/etc/profile.d/tii_on_command_not_found.sh
```

## Any other business
Have a question, idea or just want to share something? Head over to [Discussions](https://github.com/quackduck/uniclip/discussions)
