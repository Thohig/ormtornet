#!/bin/zsh
# Ormtornet.app: bygger spelet, startar tornet i bakgrunden och öppnar
# introt, som sedan leder till värdsidan. Tornet stängs med knappen
# Stäng tornet, eller av sig självt (se autostopp i README): efter 10 minuter
# utan spelare och utan synliga sidor, och senast när ingen har gjort något
# på 12 timmar. Macen hålls vaken så länge tornet är öppet.
#
# Skapad av skapa-app.sh. Ändra inte här, utan skapa om appen därifrån.
GAME="__GAME__"
GO="__GO__"
URL=http://localhost:8080
LOG=/tmp/ormtornet.log

alert() { osascript -e "display alert \"Ormtornet\" message \"$1\"" >/dev/null }
up() { curl -sf -o /dev/null -m 1 $URL/spelledare/status }

# Redan igång: öppna bara introt.
if up; then
  open $URL/
  exit 0
fi

cd "$GAME" || { alert "Hittar inte spelet i $GAME."; exit 1 }
if ! "$GO" build -o ormtornet . > $LOG 2>&1; then
  alert "Spelet gick inte att bygga. Se $LOG."
  exit 1
fi
nohup ./ormtornet -vaken >> $LOG 2>&1 &

for i in {1..50}; do up && break; sleep 0.1; done
if ! up; then
  alert "Tornet startade inte. Kanske används port 8080 av ett annat program. Se $LOG."
  exit 1
fi
open $URL/
