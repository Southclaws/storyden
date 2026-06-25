#!/usr/bin/env nu

def main [] {
    print "Building discord-robot-tools binary..."
    ^go build -trimpath -ldflags "-s -w" -o discord-robot-tools .

    print "Creating discord-robot-tools.zip..."
    if ("discord-robot-tools.zip" | path exists) {
        rm --force discord-robot-tools.zip
    }
    ^zip "discord-robot-tools.zip" "manifest.json" "discord-robot-tools"

    print "Done! Created discord-robot-tools.zip"
}
