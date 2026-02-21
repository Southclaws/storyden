#!/usr/bin/env nu

def main [] {
    print "Building crash_boot plugin..."

    go build -o crash_boot main.go

    print "Creating crash_boot.sdx..."

    ^zip "crash_boot.sdx" "manifest.json" "crash_boot"

    print "Done! Created crash_boot.sdx"
}
