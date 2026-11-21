#!/bin/sh
# Startar Ormtornet: bygger spelet, öppnar introt (som leder till spelledarens
# sida) och håller datorn vaken så länge tornet är öppet.
#
#   ./starta.sh            vanliga inställningar
#   ./starta.sh -ai 2      med två datorstyrda vandrare
#
# Alla flaggor som spelet förstår går att skicka med (se README).
# Stäng tornet med Ctrl+C eller knappen på värdsidan.
set -e
cd "$(dirname "$0")"

if ! command -v go >/dev/null 2>&1; then
	echo "Go saknas. Installera det med: brew install go"
	exit 1
fi

# Kartan snett ovanifrån (klient/) byggs om när källkoden är nyare än det
# incheckade skriptet web/iso/iso.js. Det kräver Node; annars används det
# incheckade skriptet som det är.
if [ -n "$(find klient/src klient/package.json klient/vite.config.ts -newer web/iso/iso.js 2>/dev/null)" ]; then
	if command -v npm >/dev/null 2>&1 && [ -d klient/node_modules ]; then
		echo "Bygger kartan ..."
		npm run --silent --prefix klient bygg
	else
		echo "klient/ har ändrats men byggs inte om: kör npm install i klient/ (kräver Node)."
	fi
fi

echo "Bygger tornet ..."
go build -o ormtornet .
exec ./ormtornet -oppna -vaken "$@"
