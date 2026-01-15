package models

import "time"

type CmdConfig struct {
	BaseURL          string
	RecordsPerTable  int
	BatchSize        int
	ConcurrentReqs   int           // How many concurrent requests per ID type
	RequestTimeout   time.Duration // Timeout per request
	DelayBetweenReqs time.Duration
}
