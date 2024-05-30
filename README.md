# Aporia

Aporia is a login manager that displays ascii art. It supports x11 and wayland sessions, and the [XDG base directory specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html).

<img src="https://github.com/Lunarmagpie/aporia/assets/65521138/98e321a7-9bd7-4ef8-9b1d-da0be0a3c7dc" alt="drawing" width="700"/>

Only systemd is supported.

## Compilation & Installtion
Install the dependencies:
- Go compiler
- gcc
- pam (libpam-dev on debian based distros)

Run the install.sh script as root.
```sh
$ sudo -E ./install.sh
```

You have to disable whatever display manager is running as well.

## Usage
To have ascii art you must create a file called `/etc/aporia/NAME.ascii`. Name should be the name of your ascii art.

Aporia will detect desktop environments automatically.

### Configuration
To configure Aporia, create a file called `/etc/aporia/config`. Currently the only conig option is `ascii = name`. See [`examples/config`](/examples/config) for more
information.

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
- <kbd>F11</kbd>: Shutdown
- <kbd>F12</kbd>: Reboot 

## FAQ
<details>
<summary>How do I make the font size in my TTY smaller?</summary><br>

You can edit the `FONTSIZE` variable in `/etc/default/console-setup`. Alternitively,
run the command `sudo dpkg-reconfigure console-setup`. This will allow you to safely change the
font and all font sizes available on your computer.

</details>

<details>
<summary>Why does my ascii art show up as diamonds?</summary><br>
Some terminal fonts do not support braille characters. Using an ascii art generator that does not output braille characters will fix your problem.
</details>

<details>
<summary>How do I use a fingerprint sensor to login.</summary><br>

Follow the arch wiki instructions for SDDM. More info coming in future.
NOTE: I do not reccomend because fingerprints can not decrypt drives or unlock your keyring.

</details>

<details>
<summary>Is Aporia production ready?</summary><br>

I daily drive Aporia with no issues.

</details>

<details>
<summary>How do I setup gnome-keyring or kde-wallet?</summary><br>

Both gnome-keyring and kde-wallet will automatically open with no modifications to Aporia. If you experience issues, please verify `/etc/pam.d/aporia` is correct for your setup.

</details>


## Thank You
- Thank you to the creator of [ly](https://github.com/FairyGlade/ly) for making your project under WTFPL.
I used this project to help me figure out PAM.

- Thank you to gsgx for your [display manager guide](https://gsgx.me/posts/how-to-write-a-display-manager/).

In return, I encourage people to use this project's code however they want.
