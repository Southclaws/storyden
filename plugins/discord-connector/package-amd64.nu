#!/usr/bin/env nu

def main [] {
    print "Rendering manifest.json from manifest.yaml..."
    (open manifest.yaml | to json) | save --force manifest.json

    print "Building discord-connector binary for linux/amd64..."
    with-env {
        CGO_ENABLED: "0"
        GOOS: "linux"
        GOARCH: "amd64"
    } {
        go build -trimpath -ldflags "-s -w" -o discord-connector main.go
    }

    print "Creating discord-connector.zip..."
    if ("discord-connector.zip" | path exists) {
        rm --force discord-connector.zip
    }
    ^zip "discord-connector.zip" "manifest.json" "discord-connector"

    print "Done! Created discord-connector.zip"
}
