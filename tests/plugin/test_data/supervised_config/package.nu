#!/usr/bin/env nu

def main [] {
    print "Building supervised_config plugin..."

    go build -o supervised_config main.go

    print "Creating supervised_config.sdx..."

    ^zip "supervised_config.sdx" "manifest.json" "supervised_config"

    print "Done! Created supervised_config.sdx"
}
