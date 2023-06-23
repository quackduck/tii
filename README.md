# Tii

On most GNU/Linux systems, when a command isn't found, a message showing what
to run to install the command is shown. However, macOS doesn't have
this. 

This tool adds a similar function with support for macOS using
the Homebrew package manager. Instead of simply printing the best matches, Tii shows package
descriptions and also offers to run an install command for you.

[comment]: <> ([![asciicast]&#40;https://asciinema.org/a/382511.svg&#41;]&#40;https://asciinema.org/a/382511?autoplay=1&speed=2&#41;)

<a href="https://asciinema.org/a/592995?autoplay=1&speed=2" target="_blank">
<img src="https://asciinema.org/a/592995.svg" alt="demo">
</a>

The name Tii is an acronym for "Then Install It", which is what you'll probably say when shown "Command not found".

## Installing
As of now, only macOS is supported
```shell
brew install quackduck/tap/tii
```

## Usage, environment and files

Tii will be automatically triggered if a command is not found and so you usually do not need to directly interact with it.

```text
Usage: tii [--help/-h | --version/-v | --refresh-cache/-r | <command>]

Examples:
   tii fish
   tii cowsay
   tii --help

Environment:
   TII_DISABLE_INTERACTIVE: If this variable is set to "true", Tii will
      disable interactive output (prompting for confirmation) and not install
      any packages.
   TII_AUTO_INSTALL_EXACT_MATCHES: If this variable is set to "true", Tii will
      automatically install exact matches without prompting for confirmation

Files:
   $XDG_DATA_HOME/tii: used to cache package list info. If $XDG_DATA_HOME is
      not set, ~/.local/share is used instead. Refresh the cache using the
      --refresh-cache option.
```

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
$XDG_DATA_HOME/tii or ~/.local/share/tii
```

## Any other business
Have a question, idea or just want to share something? Head over to [Discussions](https://github.com/quackduck/uniclip/discussions)
