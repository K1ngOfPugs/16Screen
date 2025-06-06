#!/bin/bash

BACKUP_FOLDER="ID_BACKUP_($(date '+%Y-%m-%d_%H-%M-%S'))"

cd /var/mobile
mkdir $BACKUP_FOLDER
cd $BACKUP_FOLDER
mv /private/var/mobile/Library/Accounts/ .
mv /private/var/mobile/Library/Prefrences/com.apple.accountsd.plist .
mkdir /private/var/mobile/Library/Accounts/

# Soft reboot to remove Apple ID #
sbreload