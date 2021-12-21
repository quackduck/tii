#!/usr/bin/env fish

function __fish_command_not_found_handler --on-event fish_command_not_found
     __fish_default_command_not_found_handler
     if status --is-interactive # make sure a human is there to agree or disagree
         echo -e "Searching for command with Tii"
         tii "$argv[1]"
     end
end

#function tii_on_command_not_found --on-event fish_command_not_found
#    functions --erase __fish_command_not_found_handler
#    functions --erase __fish_command_not_found_setup
#    if status --is-interactive # make sure a human is there to agree or disagree
#        echo -e "Searching for command with Tii"
#        tii "$argv[1]"
#    end
#end