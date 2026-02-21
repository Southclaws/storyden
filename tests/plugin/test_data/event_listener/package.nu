#!/usr/bin/env nu

def main [] {
    print "Building event_listener plugin..."

    go build -o event_listener main.go

    print "Creating event_listener.sdx..."

    ^zip "event_listener.sdx" "manifest.json" "event_listener"

    print "Done! Created event_listener.sdx"
}
