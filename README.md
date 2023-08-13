# Aporia

Aporia is a login manager that displays ascii art.

<img src="https://github.com/Lunarmagpie/aporia/assets/65521138/7c5ab59a-0aa4-45ac-983c-d7002501bfdf" alt="drawing" width="700"/>

Only systemd is supported.

## Installtion
Run the install.sh script as root
```sh
# ./install.sh
```

You have to disable whatever display manager is running as well.

## Usage
To have ascii art you must put a file in `/etc/aporia/NAME.ascii`. Name can be whatever you want. It doesn't matter.
The file must follow the format of the example file `examples/luna.ascii`. Be careful not to make an error!

## Keybinds
Aporia supports basics keybinds.

- <kbd>Enter</kbd>: Confirm
- <kbd>Tab</kbd>: Next Line
