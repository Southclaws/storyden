#!/usr/bin/env nu

def main [] {
    print "Building crash_connect plugin..."

    go build -o crash_connect main.go

    print "Creating crash_connect.sdx..."

    ^zip "crash_connect.sdx" "manifest.json" "crash_connect"

    print "Done! Created crash_connect.sdx"
}
