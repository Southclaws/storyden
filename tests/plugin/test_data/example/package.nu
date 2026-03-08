#!/usr/bin/env nu

def main [] {
    print "Building example plugin..."

    go build -o example-plugin example.go

    print "Creating example.zip..."

    ^zip "example.zip" "manifest.json" "example-plugin"

    print "Done! Created example.zip"
}
