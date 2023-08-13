# Aporia

Aporia is a login manager that displays ascii art.

<img src="https://github.com/Lunarmagpie/aporia/assets/65521138/7c5ab59a-0aa4-45ac-983c-d7002501bfdf" alt="drawing" width="700"/>

# Installtion
Run the install.sh script as root
```sh
# ./install.sh
```

You have to disable whatever display manager is running as well.

# Openrc
Currently I don't support openrc. If know how to support openrc please make a PR!

# Usage
To have ascii art you must put a file in `/etc/aporia/NAME.ascii`. Name can be whatever you want. It doesn't matter.
The file must follow the format of the example file `examples/luna.ascii`. Be careful not to make an error!
