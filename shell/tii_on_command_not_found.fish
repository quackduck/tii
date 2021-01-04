function tii_on_command_not_found --on-event fish_command_not_found
    if status --is-interactive # make sure a human is there to agree or disagree
        echo -e "Searching for command with Tii"
        tii $argv
    end
end