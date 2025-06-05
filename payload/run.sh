#!/bin/bash
cd /var/mobile/

# setup.sh #

ldid -Sentitlements.xml keychain_dumper

# updateEntitlements.sh #
KEYCHAIN_DUMPER_FOLDER=./
ENTITLEMENT_PATH=./ent.xml

dbKeychainArray=()
declare -a invalidKeychainArray=("com.apple.bluetooth"
        "com.apple.cfnetwork"
        "com.apple.cloudd"
        "com.apple.continuity.encryption"
        "com.apple.continuity.unlock"
        "com.apple.icloud.searchpartyd"
        "com.apple.ind"
        "com.apple.mobilesafari"
        "com.apple.rapport"
        "com.apple.sbd"
        "com.apple.security.sos"
        "com.apple.siri.osprey"
        "com.apple.telephonyutilities.callservicesd"
        "ichat"
        "wifianalyticsd"
      )

echo "<?xml version=\"1.0\" encoding=\"UTF-8\"?>" > $ENTITLEMENT_PATH
echo "<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">" >> $ENTITLEMENT_PATH
echo "<plist version=\"1.0\">" >> $ENTITLEMENT_PATH
echo "  <dict>" >> $ENTITLEMENT_PATH
echo "    <key>keychain-access-groups</key>" >> $ENTITLEMENT_PATH
echo "    <array>" >> $ENTITLEMENT_PATH

sqlite3 /var/Keychains/keychain-2.db "SELECT DISTINCT agrp FROM genp" > ./allgroups.txt
sqlite3 /var/Keychains/keychain-2.db "SELECT DISTINCT agrp FROM cert" >> ./allgroups.txt
sqlite3 /var/Keychains/keychain-2.db "SELECT DISTINCT agrp FROM inet" >> ./allgroups.txt
sqlite3 /var/Keychains/keychain-2.db "SELECT DISTINCT agrp FROM keys" >> ./allgroups.txt

su mobile

while IFS= read -r line; do
  dbKeychainArray+=("$line")
  if [[ ! " ${invalidKeychainArray[@]} " =~ " ${line} " ]]; then
      echo "      <string>${line}</string>">> $ENTITLEMENT_PATH
  fi
done < ./allgroups.txt

rm ./allgroups.txt

echo "    </array>">> $ENTITLEMENT_PATH
echo "    <key>platform-application</key> <true/>">> $ENTITLEMENT_PATH
echo "    <key>com.apple.private.security.no-container</key>  <true/>">> $ENTITLEMENT_PATH
echo "  </dict>">> $ENTITLEMENT_PATH
echo "</plist>">> $ENTITLEMENT_PATH

ldid -Sent.xml keychain_dumper
rm ent.xml

keychain_dumper > output.txt

# Find the passcode #

in_parental_controls=false

while IFS= read -r line; do
  if [[ "$line" =~ ^Service:\ ParentalControls ]]; then
    in_parental_controls=true
  fi

  if $in_parental_controls && [[ "$line" =~ ^Keychain\ Data:\ ([0-9]{4})$ ]]; then
    echo "${BASH_REMATCH[1]}"
    in_parental_controls=false
  fi
done < ./output.txt

rm output.txt
rm entitlements.xml
rm keychain_dumper