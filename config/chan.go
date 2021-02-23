package config

import "sync"

var notifier chan struct{}
var once sync.Once

func GetNotifier() chan struct{} {
	once.Do(func() {
		notifier = make(chan struct{})
	})
	return notifier
}
