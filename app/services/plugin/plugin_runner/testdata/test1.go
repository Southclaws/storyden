package main

import (
	"encoding/json"
	"fmt"
)

type Manifest struct {
	ID          string `json:"id"`
	Author      string `json:"author"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version"`
}

var manifest = Manifest{
	Name:    "My First Plugin",
	Version: "v1.0.2",
	ID:      "test-plugin",
	Author:  "Southclaws",
}

func main() {
	b, _ := json.Marshal(manifest)
	fmt.Println(string(b))
}
