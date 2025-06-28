package boot_time

import "time"

// StartedAt is written to in main() and read from once the HTTP server boots.
var StartedAt time.Time
