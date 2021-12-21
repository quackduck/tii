#!/usr/bin/env fish

function __fish_command_not_found_handler --on-event fish_command_not_found
     __fish_default_command_not_found_handler $argv
     echo
     if status --is-interactive # make sure a human is there to agree or disagree
         echo -e "Searching for command with Tii"
         tii "$argv[1]"
     end
end