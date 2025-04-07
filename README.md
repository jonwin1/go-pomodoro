# Go-Pomodoro

A simple pomodoro program intended for use with Waybar.

## Try go-pomodoro

```
nix run github:jonwin1/go-pomodoro
# or with flags
nix run github:jonwin1/go-pomodoro -- -w 5 -b 25
```

## Installation

### Nix

Add go-pomodoro to your system flake inputs and to your system packages.

```
inputs.pomodoro.url = "github:jonwin1/go-pomodoro"
```

```
environment.systemPackages = [ inputs.pomodoro.packages.${system}.default ];
```

### Non-Nix

#### Prerequisite

- pkg-config
- alsa-lib

#### Build

Clone this repository, navigate into it, and run `go build`.

## Configuring Waybar

Put this in waybar config and add `custom/pomodoro` to one of the modules.
Make sure that pomodoro in in your path or replace the last pomodoro in on-click with the full path to the executable.

```
"custom/pomodoro" = {
    format = "{}";
    signal = 10;
    return-type = "json";
    exec = "cat $HOME/.local/share/pomodoro/output.txt";
    on-click = "bash -c 'pgrep pomodoro && pkill pomodoro || pomodoro &'";
};
```

## TODO

- [x] Notifications
- [x] Waybar integration
- [x] Notification sound
- [x] Better output file location
- [x] Flake.nix
- [ ] Documentation
- [ ] Longer break everly few sessions
- [ ] Mouse hover label
- [ ] Waybar adjustable duration
    - [ ] Scroll to adjust time
    - [ ] Middle/Right click to change session type
    - [ ] Left click to start session
