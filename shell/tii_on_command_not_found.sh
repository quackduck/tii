#!/bin/bash
case "$(ps -cp "$$" -o command="")" in
zsh)
  command_not_found_handler() {
    if [ -t 0 ]; then
      echo -e "Searching for command with Tii"
      tii $1
    fi
  }
  ;;
bash)
  command_not_found_handle() {
    if [ -t 0 ]; then
      echo -e "Searching for command with Tii"
      tii $1
    fi
  }
  ;;
esac
