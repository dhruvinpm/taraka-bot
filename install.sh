#!/bin/bash
set -e
echo "Installing TarakaBot..."
go build -o taraka-bot .
echo "Done! Run ./taraka-bot start"
