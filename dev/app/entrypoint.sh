#!/usr/bin/env sh
set -o pipefail

#
# This script is for Alpine Linux
#

# default value
if [ "$UID" = "" ]; then
  echo "UID is empty: use 0 (root) as UID"
  UID=0
fi

# validate if input UID is number
expr "$UID" + 1 &>/dev/null
if [ $? -ne 0 ]; then
  echo "UID \"$UID\" is not number"
  exit 1
fi

USERNAME=$(getent passwd "$UID" | cut -d: -f1)
if [ $? -ne 0 ]; then
  USERNAME=dummy
  adduser -u $UID -D -S $USERNAME
fi
USERHOME=$(getent passwd "$UID" | cut -d: -f6)

sudo -u $USERNAME -E sh -c "HOME=$USERHOME $@"
