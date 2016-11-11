// Package instrumentation defines behaviors to instrument the deepshare service

package instrumentation

import "time"

// We keep them separate for now, just in case they evolve in different direction.

type Shorturl interface {
	HTTPGetDuration(d time.Duration)     // overall time performing GET HTTP Request
	HTTPPostDuration(d time.Duration)    // overall time performing POST HTTP Request
	StorageGetDuration(d time.Duration)  // time spent getting data from storage
	StorageSaveDuration(d time.Duration) // time spent saving data to the storage
}

type InappData interface {
	HTTPPostDuration(d time.Duration) // overall time performing POST HTTP Request
}

type Sharelink interface {
	HTTPGetDuration(d time.Duration) // overall time performing GET HTTP Request
}

type Match interface {
	HTTPGetDuration(d time.Duration)       // overall time performing GET HTTP Request
	HTTPPostDuration(d time.Duration)      // overall time performing POST HTTP Request
	StorageGetDuration(d time.Duration)    // time spent getting data from storage
	StorageSaveDuration(d time.Duration)   // time spent saving data to the storage
	StorageDeleteDuration(d time.Duration) // time spent deleting data from the storage
	StorageHSetDuration(d time.Duration)   // time spent on storage HSet
	StorageHGetDuration(d time.Duration)   // time spent on storage HGet
	StorageHDelDuration(d time.Duration)   // time spent on storage HDel
}
