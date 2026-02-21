#!/usr/bin/env nu

def main [] {
    print "Building crash_test plugin..."

    go build -o crash_test main.go

    print "Creating crash_test.sdx..."

    ^zip "crash_test.sdx" "manifest.json" "crash_test"

    print "Done! Created crash_test.sdx"
}
