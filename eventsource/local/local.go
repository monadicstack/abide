package local

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/monadicstack/abide/eventsource"
	"github.com/monadicstack/abide/fail"
	"github.com/monadicstack/abide/internal/slices"
)

// Broker creates a new local/in-memory broker that dispatches events to subscribers
// running just within this Go process.
func Broker(options ...BrokerOption) eventsource.Broker {
	b := broker{
		groups: map[string]*subscriptionGroup{},
		mutex:  &sync.Mutex{},
		now:    time.Now,
		errorHandler: func(err error) {
			log.Printf("[WARN] Local broker publish error: %v", err)
		},
	}
	for _, option := range options {
		option(&b)
	}
	return &b
}

type broker struct {
	mutex        *sync.Mutex
	groups       map[string]*subscriptionGroup
	now          func() time.Time
	errorHandler fail.ErrorHandler
}

func (b *broker) Publish(ctx context.Context, key string, payload []byte) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("local broker publish: %w", err)
	}

	keyTokens := b.tokenizeKey(key)

	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Yes, I realize this isn't the most efficient way to do this. It would be better to
	// use something like a radix tree - similar to how HTTP routers typically figure out
	// which handler should fire based on a path.
	//
	// That's a lot more complexity than we really need at this moment. The local broker is
	// really only a simple reference implementation since it only works if you have a
	// single instance monolith. It's not useful for much beyond playing around with the
	// framework. You're probably going to swap in NATS or something like that, so that's
	// how you make this better.
	for _, group := range b.groups {
		if !group.matches(keyTokens) {
			continue
		}

		sub := group.subscriptions.next()
		if sub == nil {
			continue
		}

		go b.publishMessage(ctx, sub, eventsource.EventMessage{
			Timestamp: b.now(),
			Key:       key,
			Payload:   payload,
		})
	}
	return nil
}

func (b *broker) publishMessage(ctx context.Context, sub *subscription, msg eventsource.EventMessage) {
	defer func() {
		if recovery := recover(); recovery != nil {
			err, _ := recovery.(error)
			b.errorHandler(fmt.Errorf("local broker publish: %s: %w", sub.group.key, err))
		}
	}()

	if err := sub.handlerFunc(ctx, &msg); err != nil {
		b.errorHandler(fmt.Errorf("local broker publish: %s: %w", sub.group.key, err))
	}
}

func (b *broker) Subscribe(key string, handlerFunc eventsource.EventHandlerFunc) (eventsource.Subscription, error) {
	// We want this handler to absolutely fire no matter what other subscribers there are,
	// so create a unique group id, making this a consumer group of 1.
	group := strconv.FormatInt(time.Now().UnixNano(), 10)
	return b.SubscribeGroup(key, group, handlerFunc)
}

func (b *broker) SubscribeGroup(key string, groupKey string, handlerFunc eventsource.EventHandlerFunc) (eventsource.Subscription, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	group := b.loadGroup(key, groupKey)
	sub := subscription{
		broker:      b,
		group:       group,
		handlerFunc: handlerFunc,
	}
	return group.subscriptions.append(&sub), nil
}

func (b *broker) unsubscribe(sub *subscription) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	sub.group.subscriptions.remove(sub)
}

func (b *broker) loadGroup(key string, groupKey string) *subscriptionGroup {
	lookupKey := key + "-" + groupKey

	if group, ok := b.groups[lookupKey]; ok {
		return group
	}

	group := &subscriptionGroup{
		broker:        b,
		key:           key,
		keyTokens:     b.tokenizeKey(key),
		groupKey:      groupKey,
		subscriptions: &subscriptionRoundRobin{},
	}
	b.groups[lookupKey] = group
	return group
}

func (b *broker) tokenizeKey(key string) []string {
	return strings.Split(key, ".")
}

// ---------------------------------
// SUBSCRIPTION MANAGEMENT
// ---------------------------------

type subscriptionGroup struct {
	broker        *broker
	key           string
	keyTokens     []string
	groupKey      string
	subscriptions *subscriptionRoundRobin
}

// matches determines if an incoming keyTokens should be handled by this subscription. This
// compares the individual segments, allowing "*" to match any segment. Here are
// some examples:
//
//	subs = subscription{keyTokens: []string{"foo"}
//	subs.matches("foo")        // <-- true
//	subs.matches("bar")        // <-- false
//	subs.matches("foo", "bar") // <-- false
//
//	subs = subscription{keyTokens: []string{"foo", "bar"}
//	subs.matches("foo")               // <-- false
//	subs.matches("foo", "bar")        // <-- true
//	subs.matches("foo", "baz")        // <-- false
//	subs.matches("foo", "bar", "baz") // <-- false
//
//	subs = subscription{keyTokens: []string{"foo", "*", "baz"}
//	subs.matches("foo")               // <-- false
//	subs.matches("foo", "bar", "baz") // <-- true
//	subs.matches("foo", "*", "baz")   // <-- true
//	subs.matches("foo", "*", "*")     // <-- true
//	subs.matches("foo", "bar", "*")   // <-- true
//	subs.matches("foo", "baz", "*")   // <-- false
func (group *subscriptionGroup) matches(incomingKey []string) bool {
	// Even wildcards only match one token, so the number of tokens must be the same.
	if len(incomingKey) != len(group.keyTokens) {
		return false
	}

	for i, token := range group.keyTokens {
		if token == "*" {
			continue
		}
		if incomingKey[i] == "*" {
			continue
		}
		if token != incomingKey[i] {
			return false
		}
	}
	return true
}

type subscription struct {
	broker      *broker
	group       *subscriptionGroup
	handlerFunc eventsource.EventHandlerFunc
}

func (sub *subscription) Unsubscribe() error {
	sub.broker.unsubscribe(sub)
	return nil
}

// ---------------------------------
// ROUND ROBIN MANAGEMENT
// ---------------------------------

type subscriptionRoundRobin struct {
	// index is our cursor that indicates the next subscription to return in the round robin.
	index int
	// subscriptions is the raw slice/ring of handlers we're rotating through.
	subscriptions []*subscription
}

func (robin *subscriptionRoundRobin) next() *subscription {
	if len(robin.subscriptions) == 0 {
		return nil
	}

	next := robin.subscriptions[robin.index]
	robin.index = (robin.index + 1) % len(robin.subscriptions)
	return next
}

func (robin *subscriptionRoundRobin) append(sub *subscription) *subscription {
	robin.subscriptions = append(robin.subscriptions, sub)
	return sub
}

func (robin *subscriptionRoundRobin) remove(sub *subscription) {
	robin.subscriptions = slices.Remove(robin.subscriptions, sub)

	// The cursor was at the end of the list, so now that there's one less we need to just go
	// back to the beginning.
	if robin.index >= len(robin.subscriptions) {
		robin.index = 0
	}
}

// ---------------------------------
// OPTIONS
// ---------------------------------

// BrokerOption allows you to tweak the local broker's behavior in some way.
type BrokerOption func(*broker)

// WithErrorHandler swaps the default error handler for this one.
func WithErrorHandler(handler fail.ErrorHandler) BrokerOption {
	return func(broker *broker) {
		broker.errorHandler = handler
	}
}
