package structs

import "time"

type Proxy struct {
	Addr   string
	Time   time.Duration
	IsBusy bool
}
