#!/bin/sh

GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -o close-the-wardrobe-please .
ssh dorne 'sudo systemctl stop close-the-wardrobe-please.service'
scp close-the-wardrobe-please dorne:/opt/sensors/bin/close-the-wardrobe-please
scp close-the-wardrobe-please.service dorne:/tmp/close-the-wardrobe-please.service
ssh dorne 'sudo cp /tmp/close-the-wardrobe-please.service /etc/systemd/system/close-the-wardrobe-please.service'

ssh dorne 'sudo systemctl daemon-reload'
ssh dorne 'sudo systemctl enable close-the-wardrobe-please.service'
ssh dorne 'sudo systemctl start close-the-wardrobe-please.service'
