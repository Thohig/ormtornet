#!/bin/sh
# Bygger Ormtornet som en sida för GitHub Pages (docs/github-pages.md):
# sidorna i web/, sidornas texter, tornet som WebAssembly (torn.wasm) och
# Go:s wasm_exec.js, i en mapp.
#
#   verktyg/bygg-sida.sh            till mappen sida/
#   verktyg/bygg-sida.sh ut         till mappen ut/
#
# Prova lokalt med en enkel webbserver i mappen, till exempel
#   python3 -m http.server 8000 -d sida
# och öppna http://localhost:8000/.
set -e
cd "$(dirname "$0")/.."
dir="${1:-sida}"

rm -rf "$dir"
go run . -sida "$dir"
GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o "$dir/torn.wasm" .
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" "$dir/"
echo "Tornet är byggt i $dir/ ($(du -h "$dir/torn.wasm" | cut -f1) för torn.wasm)."
