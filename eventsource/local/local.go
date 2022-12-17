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
		subscriptions: map[string]subscriptionGroup{},
		mutex:         &sync.Mutex{},
		now:           time.Now,
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
	mutex         *sync.Mutex
	subscriptions map[string]subscriptionGroup
	now           func() time.Time
	errorHandler  fail.ErrorHandler
}

func (b broker) Publish(ctx context.Context, key string, payload []byte) error {
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
	for _, group := range b.subscriptions {
		sub := group.next(keyTokens)
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
	if err := sub.handlerFunc(ctx, &msg); err != nil {
		b.errorHandler(fmt.Errorf("local broker publish: %s: %w", sub.keyTokens, err))
	}
}

func (b *broker) Subscribe(key string, handlerFunc eventsource.EventHandlerFunc) (eventsource.Subscription, error) {
	// We want this handler to absolutely fire no matter what other subscribers there are,
	// so create a unique group id, making this a consumer group of 1.
	group := strconv.FormatInt(time.Now().UnixNano(), 10)
	return b.SubscribeGroup(key, group, handlerFunc)
}

func (b *broker) SubscribeGroup(key string, group string, handlerFunc eventsource.EventHandlerFunc) (eventsource.Subscription, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	sub := subscription{
		broker:      b,
		key:         key,
		keyTokens:   b.tokenizeKey(key),
		group:       group,
		handlerFunc: handlerFunc,
	}
	b.subscriptions[group] = append(b.subscriptions[group], &sub)
	return &sub, nil
}

func (b *broker) tokenizeKey(key string) []string {
	return strings.Split(key, ".")
}

type subscriptionGroup []*subscription

func (group subscriptionGroup) append(sub *subscription) subscriptionGroup {
	return append(group, sub)
}

func (group subscriptionGroup) remove(sub *subscription) subscriptionGroup {
	return slices.Remove(group, sub)
}

func (group subscriptionGroup) next(keyTokens []string) *subscription {
	// We do a half-assed round-robin by locating the first match in the
	// group. If we find one, we pull it out of the slice and move it to the
	// end of the list; shifting the elements after it left by one. The match
	// gets put at the end of the list and returned as our match.

	var match *subscription
	for i, sub := range group {
		// We're planning on putting the match at the end of the group, so it's
		// at the end of the round-robin. We have that match in hand, so just
		// keep shifting the rest of the group left so that we have a space open
		// at the end to put the match.
		if match != nil {
			group[i-1] = sub
			continue
		}
		if sub.matches(keyTokens...) {
			match = sub
		}
	}
	if match == nil {
		return nil
	}
	group[len(group)-1] = match
	return match
}

type subscription struct {
	broker      *broker
	key         string
	keyTokens   []string
	group       string
	handlerFunc eventsource.EventHandlerFunc
}

func (sub *subscription) Unsubscribe() error {
	sub.broker.mutex.Lock()
	defer sub.broker.mutex.Unlock()

	if group, ok := sub.broker.subscriptions[sub.group]; ok {
		sub.broker.subscriptions[sub.group] = group.remove(sub)
	}
	return nil
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
func (sub *subscription) matches(incomingKey ...string) bool {
	// Even wildcards only match one token, so the number of tokens must be the same.
	if len(incomingKey) != len(sub.keyTokens) {
		return false
	}

	for i, token := range sub.keyTokens {
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

// BrokerOption allows you to tweak the local broker's behavior in some way.
type BrokerOption func(*broker)

// WithErrorHandler swaps the default error handler for this one.
func WithErrorHandler(handler fail.ErrorHandler) BrokerOption {
	return func(broker *broker) {
		broker.errorHandler = handler
	}
}
