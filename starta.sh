#!/bin/sh
# Startar Ormtornet: bygger spelet, öppnar spelledarens sida och håller
# datorn vaken så länge tornet är öppet.
#
#   ./starta.sh            vanliga inställningar
#   ./starta.sh -ai 2      med två datorstyrda vandrare
#
# Alla flaggor som spelet förstår går att skicka med (se README).
# Stäng tornet med Ctrl+C eller knappen på spelledarens sida.
set -e
cd "$(dirname "$0")"

if ! command -v go >/dev/null 2>&1; then
	echo "Go saknas. Installera det med: brew install go"
	exit 1
fi

echo "Bygger tornet ..."
go build -o ormtornet .
exec ./ormtornet -oppna -vaken "$@"
