#! /bin/bash
# Xsession - run as user
# Copyright (C) 2016 Pier Luigi Fiorini <pierluigi.fiorini@gmail.com>

# This file is extracted from kde-workspace (kdm/kfrontend/genkdmconf.c)
# Copyright (C) 2001-2005 Oswald Buddenhagen <ossi@kde.org>

# Note that the respective logout scripts are not sourced.

# Aporia uses a modified version of this file to launch the user's startx script.

# Load Xsession scripts
# OPTIONFILE, USERXSESSION, USERXSESSIONRC and ALTUSERXSESSION are required
# by the scripts to work
xsessionddir="/etc/X11/Xsession.d"
OPTIONFILE=/etc/X11/Xsession.options
USERXSESSION=$HOME/.xsession
USERXSESSIONRC=$HOME/.xsessionrc
ALTUSERXSESSION=$HOME/.Xsession

if [ -d "$xsessionddir" ]; then
    for i in `ls $xsessionddir`; do
        if [ "$i" = "99x11-common_start" ]; then
          continue
        fi
        script="$xsessionddir/$i"
        echo "Loading X session script $script"
        if [ -r "$script"  -a -f "$script" ] && expr "$i" : '^[[:alnum:]_-]\+$' > /dev/null; then
            . "$script"
        fi
    done
fi

if [ -d /etc/X11/Xresources ]; then
  for i in /etc/X11/Xresources/*; do
    [ -f $i ] && xrdb -merge $i
  done
elif [ -f /etc/X11/Xresources ]; then
  xrdb -merge /etc/X11/Xresources
fi
[ -f $HOME/.Xresources ] && xrdb -merge $HOME/.Xresources
[ -f $XDG_CONFIG_HOME/X11/Xresources ] && xrdb -merge $XDG_CONFIG_HOME/X11/Xresources

if [ -f "$USERXSESSION" ]; then
  . "$USERXSESSION"
fi

echo "Aporia startx path: $APORIA_STARTX_PATH"

exec $APORIA_STARTX_PATH
