package subpub

import (
	"context"
	"errors"
	"sync"
)

type MessageHandler func(msg interface{})

type Subscription interface {
	Unsubscribe()
}

type SubPub interface {
	Subscripe(subject string, cb MessageHandler) (Subscription, error)

	Publish(subject string, msg interface{}) error

	Close(ctx context.Context) error
}

func NewSubPub() SubPub {
	return &subPub{
		subs: make(map[string][]*sub),
	}
}

type sub struct {
	handler MessageHandler
	channel chan interface{}
	stop    chan struct{}
	once    sync.Once
}

func (s *sub) start(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case msg := <-s.channel:
			s.handler(msg)
		case <-s.stop:
			return
		}
	}
}

func (s *sub) close() {
	s.once.Do(func() {
		close(s.stop)
	})
}

type subscription struct {
	unsubOnce sync.Once
	unsubFn   func()
}

func (s *subscription) Unsubscribe() {
	s.unsubOnce.Do(s.unsubFn)
}

type subPub struct {
	mu        sync.Mutex
	subs      map[string][]*sub
	wg        sync.WaitGroup
	is_closed bool
}

func (sp *subPub) Subscripe(subject string, cb MessageHandler) (Subscription, error) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if sp.is_closed {
		return nil, errors.New("Supbup closed")
	}

	sub := &sub{
		handler: cb,
		channel: make(chan interface{}, 16),
		stop:    make(chan struct{}),
	}

	sp.subs[subject] = append(sp.subs[subject], sub)

	sp.wg.Add(1)
	go sub.start(&sp.wg)

	unsub := func() {
		sp.mu.Lock()
		defer sp.mu.Unlock()
		subs := sp.subs[subject]
		for i, s := range subs {
			if s == sub {
				sp.subs[subject] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
	}

	return &subscription{
		unsubFn: unsub,
	}, nil
}

func (sp subPub) Publish(subject string, msg interface{}) error {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if sp.is_closed {
		return errors.New("Subpub is closed")
	}

	for _, s := range sp.subs[subject] {
		select {
		case s.channel <- msg:
		default:
			go func(ch chan interface{}) {
				ch <- msg
			}(s.channel)

		}
	}

	return nil
}

func (sp *subPub) Close(ctx context.Context) error {
	sp.mu.Lock()
	if sp.is_closed {
		sp.mu.Unlock()
		return nil
	}
	sp.is_closed = true
	subList := make(map[string][]*sub)
	for subject, subs := range sp.subs {
		subList[subject] = append([]*sub{}, subs...)
	}

	sp.mu.Unlock()

	done := make(chan struct{})

	go func() {
		for _, subs := range subList {
			for _, s := range subs {
				s.close()
			}
		}
		sp.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
