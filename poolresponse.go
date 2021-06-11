package goncu

import "time"

type PoolResponse struct {
	Duration time.Duration
	Errors   []error
	Hits     []interface{}
}
