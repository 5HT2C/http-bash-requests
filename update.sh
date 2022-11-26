#!/bin/bash

# Ensure up to date
git pull || exit $?

# If bin folder does not exist, make it
if [[ ! -d ~/.local/bin ]]; then
    mkdir -p ~/.local/bin || exit $?
fi

# If systemd config folder does not exist, make it
if [[ ! -d ~/.config/systemd/user ]]; then
    mkdir -p ~/.config/systemd/user || exit $?
fi

# Stop services
systemctl --user stop http-bash-requests.timer
systemctl --user stop http-bash-requests.service

# If timer does not exist, copy it
if [[ ! -f ~/.config/systemd/user/http-bash-requests.timer ]]; then
    cp http-bash-requests.timer ~/.config/systemd/user/http-bash-requests.timer || exit $?
fi

# If service does not exist, copy it
if [[ ! -f ~/.config/systemd/user/http-bash-requests.service ]]; then
    cp http-bash-requests.service ~/.config/systemd/user/http-bash-requests.service || exit $?
fi

# Re-build
go build -o ~/.local/bin/http-bash-requests .

# Reload daemon and start service again
systemctl --user daemon-reload
systemctl --user enable --now http-bash-requests.timer

# Get status
systemctl --user status http-bash-requests.service
