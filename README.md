# Go-Pomodoro

A simple pomodoro program intended for use with Waybar.

## Installation

### Prerequisite

- pkg-config
- alsa-lib

### Pomodoro

Clone this repository, navigate into it, and run `go build`.

### Waybar

Put this in waybar config and add `custom/pomodoro` to one of the modules.
Remember to replace <path-to-go-pomodoro> with the path to your clone of this repo.

```
"custom/pomodoro" = {
    format = "{}";
    signal = 10;
    return-type = "json";
    exec = "cat <path-to-go-pomodoro>/log";
    on-click = "bash -c 'pgrep pomodoro && pkill pomodoro || <path-to-go-pomodoro>/pomodoro &'";
};
```

## TODO

- [x] Notifications
- [x] Waybar integration
- [x] Notification sound
- [ ] Better output file location
- [ ] Proper install
- [ ] Flake.nix
- [ ] Documentation
- [ ] Longer break everly few sessions
