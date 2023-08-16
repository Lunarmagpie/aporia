# Aporia

Aporia is a login manager that displays ascii art. It supports x11 and wayland sessions.

<img src="https://github.com/Lunarmagpie/aporia/assets/65521138/36184c2a-9fc8-45ce-8dd2-afa6b721fc24" alt="drawing" width="700"/>

Only systemd is supported.

## Installtion
Install the dependencies:
- Go compiler
- gcc
- pam (libpam-dev on ubuntu)

Run the install.sh script as root.
```sh
$ sudo -E ./install.sh
```

You have to disable whatever display manager is running as well.

## Usage
To have ascii art you must put a file in `/etc/aporia/NAME.ascii`. Name can be whatever you want. It doesn't matter.
The file must follow the format of the example file `examples/luna.ascii`. Be careful not to make an error!

Aporia will detect desktop environments automatically.

### Adding Custom Desktop Environments
Custom desktop environments are added as scripts.

#### Adding a bspwm environment (X11)
Create a file called `bspwm.x11` and put it in the `/etc/aporia` directory.
The file is used as your xinitrc.

```sh
#!/bin/bash
exec bspwm
```

#### Adding a hyprland environment (wayland)
Create a file called `hyprland.wayland` and put it in the `/etc/aporia` directory.
This file is run to start your window manager.

```sh
#!/bin/bash
exec Hyprland
```

These files MUST be set to executable with `sudo chmod +x filename`.

## Keybinds
Aporia supports basics keybinds.

- <kbd>Enter</kbd>: Confirm
- <kbd>Tab</kbd>: Next Line



## FAQ
<details>
<summary>How do I make the font size in my TTY smaller?</summary><br>

Edit the `FONTSIZE` variable in `/etc/default/console-setup`. I use fontzie
`16x12` on my 1920x1080 computer.
</details>

<details>
<summary>Why does my ascii art show up as diamonds?</summary><br>

The `terminus` font does not support certain braille characters. Using a different ascii art generator will likely
fix your problem.
</details>

## Thank You
- Thank you to the creator of [ly](https://github.com/FairyGlade/ly) for making your project under WTFPL.
I used this project to help me figure out PAM.

- Thank you to gsgx for your [display manager guide](https://gsgx.me/posts/how-to-write-a-display-manager/).

In return, I encourage people to use this project's code however they want.
