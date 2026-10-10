package namegen

import (
	"math/rand/v2"
	"strings"
)

// Fixed arrays keep the vocabulary in static data without startup parsing.
var moods = [...]string{
	"brave", "bright", "brisk", "calm", "clever", "cosy", "curious", "daring",
	"dreamy", "eager", "fair", "fierce", "fond", "gentle", "glad", "happy",
	"hearty", "humble", "jolly", "keen", "kind", "lively", "lucky", "merry",
	"mighty", "nimble", "noble", "patient", "placid", "playful", "proud", "quiet",
	"quirky", "quaint", "quick", "ready", "regal", "rested", "robust", "rosy",
	"royal", "serene", "shy", "sleepy", "smart", "snug", "social", "spry",
	"steady", "subtle", "sunny", "sweet", "swift", "tender", "tidy", "tiny",
	"trusty", "upbeat", "valiant", "vivid", "warm", "wary", "wild", "wise",
}

var colours = [...]string{
	"amber", "aqua", "azure", "beige", "blue", "bronze", "brown", "buff",
	"cherry", "coral", "cream", "cyan", "dun", "ebony", "fawn", "ginger",
	"gold", "golden", "green", "grey", "hazel", "indigo", "ivory", "jade",
	"khaki", "lilac", "lime", "linen", "maroon", "mauve", "mint", "moss",
	"ochre", "olive", "orange", "peach", "pearl", "pink", "plum", "purple",
	"red", "rose", "ruby", "russet", "rust", "sable", "salmon", "sand",
	"scarlet", "sepia", "silver", "slate", "smoke", "snow", "steel", "tan",
	"tawny", "teal", "topaz", "umber", "violet", "wheat", "white", "yellow",
}

var places = [...]string{
	"alder", "alpine", "apple", "aspen", "autumn", "bamboo", "bay", "birch",
	"bloom", "brook", "cedar", "cherry", "cliff", "cloud", "coast", "coral",
	"cove", "creek", "dawn", "desert", "dune", "dusk", "elm", "ember",
	"fern", "field", "fir", "flower", "forest", "frost", "garden", "glade",
	"grove", "heath", "hill", "hollow", "island", "lake", "laurel", "leaf",
	"linden", "maple", "marsh", "meadow", "mist", "moon", "moss", "oak",
	"ocean", "pine", "pond", "rain", "reed", "river", "rock", "shore",
	"spring", "spruce", "star", "stone", "stream", "summer", "thorn", "willow",
}

var characters = [...]string{
	"badger", "bard", "bear", "bee", "beetle", "bird", "bunny", "cat",
	"clerk", "colt", "crane", "crow", "deer", "dove", "druid", "duck",
	"elf", "faerie", "falcon", "finch", "firefly", "fish", "fawn", "fox",
	"frog", "gecko", "gnome", "goose", "hare", "hawk", "heron", "jay",
	"keeper", "kitten", "knight", "koala", "lark", "lynx", "mage", "mole",
	"monk", "moth", "mouse", "newt", "otter", "owl", "panda", "poet",
	"pup", "quail", "rabbit", "raven", "robin", "sage", "scribe", "seal",
	"snail", "sparrow", "sprite", "stoat", "swan", "toad", "wolf", "wren",
}

// Generate returns a word-only handle from 16,777,216 combinations, at most 30 bytes.
func Generate() string {
	return strings.Join([]string{
		moods[rand.IntN(len(moods))],
		colours[rand.IntN(len(colours))],
		places[rand.IntN(len(places))],
		characters[rand.IntN(len(characters))],
	}, "-")
}
