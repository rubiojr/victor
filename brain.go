package victor

import (
	"github.com/brettbuddin/victor/adapter"
	"log"
	"regexp"
	"strings"
	"sync"
)

type Brain struct {
	*sync.RWMutex
	name      string
	identity  adapter.User
	listeners []ListenerFunc
	cache     *Cache
}

func NewBrain(name string) *Brain {
	return &Brain{
		RWMutex:   &sync.RWMutex{},
		name:      name,
		listeners: []ListenerFunc{},
		cache:     NewCache(),
	}
}

func (b *Brain) Cache() adapter.Cacher {
	return b.cache
}

func (b *Brain) Name() string {
	b.RLock()
	defer b.RUnlock()
	return b.name
}

func (b *Brain) Identity() adapter.User {
	b.RLock()
	defer b.RUnlock()
	return b.identity
}

func (b *Brain) SetIdentity(u adapter.User) {
	b.Lock()
	defer b.Unlock()
	b.identity = u
}

// Respond registers a listener that matches a direct message to
// the robot. This means "@botname command", "botname command" or
// "/command" forms
func (b *Brain) Respond(exp string, f func(adapter.Message)) (err error) {
	exp = strings.Join([]string{
		"(?i)", // flags
		"\\A",  // begin
		"(?:(?:@)?" + b.Name() + "[:,]?\\s*|/)", // bot name
		"(?:" + exp + ")",                       // expression
		"\\z",                                   // end
	}, "")

	return b.Hear(exp, f)
}

// Hear registers a listener that matches any instance of the
// phrase in the room. Excluding from itself.
func (b *Brain) Hear(exp string, f func(adapter.Message)) (err error) {
	pattern, err := regexp.Compile(exp)

	if err != nil {
		return err
	}

	log.Printf("LISTENER: %s\n", exp)
	b.register(listenerFunc(pattern, f))
	return
}

func (b *Brain) register(l ListenerFunc) {
	b.Lock()
	defer b.Unlock()
	b.listeners = append(b.listeners, l)
}

// Receive accepts an incoming message and applies
// it to all listeners.
func (b *Brain) Receive(m adapter.Message) {
	b.RLock()
	defer b.RUnlock()
	for _, l := range b.listeners {
		go l(m)
	}
}

type ListenerFunc func(adapter.Message)

func listenerFunc(pattern *regexp.Regexp, f ListenerFunc) ListenerFunc {
	return func(m adapter.Message) {
		results := pattern.FindAllStringSubmatch(m.Body(), -1)

		if len(results) > 0 {
			m.SetParams(results[0][1:])
			log.Printf("TRIGGER: %s\n", pattern)
			log.Printf("TRIGGER: %s\n", pattern)
			log.Printf("PARAMS: %s\n", m.Params())
			f(m)
		}
	}
}
