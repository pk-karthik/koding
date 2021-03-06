package client

import (
	"context"
	"reflect"
	"sync"
	"time"

	"koding/klient/fs"
	"koding/klient/machine/index"
)

// Cached allows user to cache Client method calls results. It is limited to
// unique argument calls so it is only useful for multiple identical method
// invocations. This type is thread safe but will not cache its results until
// first ones are available.
type Cached struct {
	currentUser    func() (string, error)                   // CurrentUser.
	sshAddKeys     func(string, ...string) error            // SSHAddKeys.
	mountHeadIndex func(string) (string, int, int64, error) // MountHeadIndex.
	mountGetIndex  func(string) (*index.Index, error)       // MountGetIndex.
	diskInfo       func(string) (fs.DiskInfo, error)        // DiskBlocks

	c Client // Client used by Context method.
}

var _ Client = (*Cached)(nil)

// NewCached creates a new Cached client instance.
func NewCached(c Client, interval time.Duration) *Cached {
	return &Cached{
		currentUser:    currentUser(c, interval),
		sshAddKeys:     sshAddKeys(c, interval),
		mountHeadIndex: mountHeadIndex(c, interval),
		mountGetIndex:  mountGetIndex(c, interval),
		diskInfo:       diskInfo(c, interval),
		c:              c,
	}
}

// CurrentUser calls registered Client's CurrentUser method and caches its
// result for the specified interval. It doesn't cache results from disconnected
// client.
func (c *Cached) CurrentUser() (string, error) {
	return c.currentUser()
}

func currentUser(c Client, interval time.Duration) func() (string, error) {
	lastCall, mu := time.Now().Add(-interval-time.Second), sync.Mutex{}

	rUser, rErr := "", error(nil)

	return func() (user string, err error) {
		mu.Lock()
		if time.Since(lastCall) < interval && rErr != ErrDisconnected {
			user, err = rUser, rErr
			mu.Unlock()
			return
		}
		mu.Unlock()

		user, err = c.CurrentUser()

		mu.Lock()
		rUser, rErr = user, err
		lastCall = time.Now()
		mu.Unlock()

		return
	}
}

// SSHAddKeys calls registered Client's SSHAddKeys method and caches its result
// for the specified interval. It doesn't cache results from disconnected
// client. If call arguments change, the cache will be invalidated.
func (c *Cached) SSHAddKeys(username string, keys ...string) error {
	return c.sshAddKeys(username, keys...)
}

func sshAddKeys(c Client, interval time.Duration) func(string, ...string) error {
	lastCall, mu := time.Now().Add(-interval-time.Second), sync.Mutex{}

	aUsername, aKeys := "", []string(nil)
	rErr := error(nil)

	return func(username string, keys ...string) (err error) {
		mu.Lock()
		if time.Since(lastCall) < interval && aUsername == username &&
			reflect.DeepEqual(aKeys, keys) && rErr != ErrDisconnected {
			err = rErr
			mu.Unlock()
			return
		}
		mu.Unlock()

		err = c.SSHAddKeys(username, keys...)

		mu.Lock()
		aUsername, aKeys = username, keys
		rErr = err
		lastCall = time.Now()
		mu.Unlock()

		return
	}
}

// MountHeadIndex calls registered Client's MountHeadIndex method and caches its
// result for the specified interval. It doesn't cache results from disconnected
// client. If call arguments change, the cache will be invalidated.
func (c *Cached) MountHeadIndex(path string) (string, int, int64, error) {
	return c.mountHeadIndex(path)
}

func mountHeadIndex(c Client, interval time.Duration) func(string) (string, int, int64, error) {
	lastCall, mu := time.Now().Add(-interval-time.Second), sync.Mutex{}

	aPath := ""
	rAbsPath, rCount, rDiskSize, rErr := "", 0, int64(0), error(nil)

	return func(path string) (absPath string, count int, diskSize int64, err error) {
		mu.Lock()
		if time.Since(lastCall) < interval && aPath == path && rErr != ErrDisconnected {
			absPath, count, diskSize, err = rAbsPath, rCount, rDiskSize, rErr
			mu.Unlock()
			return
		}
		mu.Unlock()

		absPath, count, diskSize, err = c.MountHeadIndex(path)

		mu.Lock()
		aPath = path
		rAbsPath, rCount, rDiskSize, rErr = absPath, count, diskSize, err
		lastCall = time.Now()
		mu.Unlock()

		return
	}
}

// MountGetIndex calls registered Client's MountGetIndex method and caches its
// result for the specified interval. It doesn't cache results from disconnected
// client. If call arguments change, the cache will be invalidated.
func (c *Cached) MountGetIndex(path string) (*index.Index, error) {
	return c.mountGetIndex(path)
}

func mountGetIndex(c Client, interval time.Duration) func(string) (*index.Index, error) {
	lastCall, mu := time.Now().Add(-interval-time.Second), sync.Mutex{}

	aPath := ""
	rIdx, rErr := (*index.Index)(nil), error(nil)

	return func(path string) (idx *index.Index, err error) {
		mu.Lock()
		if time.Since(lastCall) < interval && aPath == path && rErr != ErrDisconnected {
			idx, err = rIdx, rErr
			mu.Unlock()
			return
		}
		mu.Unlock()

		idx, err = c.MountGetIndex(path)

		mu.Lock()
		aPath = path
		rIdx, rErr = idx, err
		lastCall = time.Now()
		mu.Unlock()

		return
	}
}

// DiskInfo calls registered Client's DiskInfo method and caches its result for
// the specified interval. It doesn't cache results from disconnected client. If
// call arguments change, the cache will be invalidated.
func (c *Cached) DiskInfo(path string) (fs.DiskInfo, error) {
	return c.diskInfo(path)
}

func diskInfo(c Client, interval time.Duration) func(string) (fs.DiskInfo, error) {
	lastCall, mu := time.Now().Add(-interval-time.Second), sync.Mutex{}

	aPath := ""
	rDi, rErr := fs.DiskInfo{}, error(nil)

	return func(path string) (di fs.DiskInfo, err error) {
		mu.Lock()
		if time.Since(lastCall) < interval && aPath == path && rErr != ErrDisconnected {
			di, err = rDi, rErr
			mu.Unlock()
			return
		}
		mu.Unlock()

		di, err = c.DiskInfo(path)

		mu.Lock()
		aPath = path
		rDi, rErr = di, err
		lastCall = time.Now()
		mu.Unlock()

		return
	}
}

// Context calls registered Client's Context without any cache.
func (c *Cached) Context() context.Context {
	return c.c.Context()
}
