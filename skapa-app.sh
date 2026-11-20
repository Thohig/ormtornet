#!/bin/zsh
# Skapar Ormtornet.app med tornets ikon. Dubbelklick på appen startar tornet
# i bakgrunden och öppnar värdsidan i webbläsaren.
#
#   ./skapa-app.sh             skapar ~/Projekt/Lokalt/Ormtornet.app
#   ./skapa-app.sh ~/Desktop   skapar appen i en annan mapp
#
# Appen bygger alltid spelet från den här mappen, så den behöver bara skapas
# om när mappen flyttas eller Go installeras på nytt ställe.
set -e
GAME=${0:A:h}
DEST=${1:-$HOME/Projekt/Lokalt}
APP=$DEST/Ormtornet.app
ID=se.ormtornet.start

GO=$(command -v go) || { echo "Go saknas. Installera det med: brew install go"; exit 1 }

# Ersätt bara en tidigare Ormtornet.app, aldrig något annat med samma namn.
if [[ -e $APP ]]; then
  if [[ $(/usr/libexec/PlistBuddy -c "Print :CFBundleIdentifier" "$APP/Contents/Info.plist" 2>/dev/null) != $ID ]]; then
    echo "$APP finns redan och är inte en Ormtornet-app som det här skriptet har skapat. Rör den inte."
    exit 1
  fi
  rm -rf "$APP"
fi

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources"
cd "$GAME"
echo "Ritar ikonen ..."
"$GO" run ./verktyg/appikon "$TMP/Ormtornet.iconset"
iconutil -c icns -o "$APP/Contents/Resources/Ormtornet.icns" "$TMP/Ormtornet.iconset"

cat > "$APP/Contents/Info.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleName</key><string>Ormtornet</string>
  <key>CFBundleDisplayName</key><string>Ormtornet</string>
  <key>CFBundleIdentifier</key><string>$ID</string>
  <key>CFBundleExecutable</key><string>Ormtornet</string>
  <key>CFBundleIconFile</key><string>Ormtornet</string>
  <key>CFBundlePackageType</key><string>APPL</string>
  <key>CFBundleShortVersionString</key><string>1.0</string>
  <key>LSUIElement</key><true/>
</dict>
</plist>
EOF

sed -e "s|__GAME__|$GAME|" -e "s|__GO__|$GO|" verktyg/app-start.zsh > "$APP/Contents/MacOS/Ormtornet"
chmod +x "$APP/Contents/MacOS/Ormtornet"
touch "$APP" # så att Finder visar den nya ikonen

echo "Klart: $APP"
echo "Dubbelklicka på den för att öppna tornet. Dra den till Dock om du vill ha den där."
