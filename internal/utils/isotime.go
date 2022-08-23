package utils

import "time"

func FormatISO(t time.Time) string { return t.Format(time.RFC3339) }
