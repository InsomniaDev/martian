package pubsub

import "sync"

type Pubsub struct {
	mu   sync.RWMutex
	subs map[string][]chan string
}

var Service *Pubsub

func init() {
	Service = &Pubsub{}
	Service.subs = make(map[string][]chan string)
}

func (ps *Pubsub) Subscribe(topic string, ch chan string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.subs[topic] = append(ps.subs[topic], ch)
}

func (ps *Pubsub) Publish(topic string, msg string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, ch := range ps.subs[topic] {
		ch <- msg
	}
}
