#!/bin/bash
if [ "$ZSH_VERSION" ] && [ -t 0 ]; then # is zsh and is interactive
  command_not_found_handler() {
    echo "Searching for command with Tii"
    tii "$1"
  }
elif [ "$BASH_VERSION" ] && [ -t 0 ]; then # is bash and is interactive
  command_not_found_handle() {
    echo "Searching for command with Tii"
    tii "$1"
  }
else
  echo -e "This is not bash or zsh. Please run $0 only in these shells."
fi
