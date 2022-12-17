//go:build unit

package local_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/monadicstack/abide/eventsource"
	"github.com/monadicstack/abide/eventsource/local"
	"github.com/monadicstack/abide/internal/testext"
	"github.com/stretchr/testify/suite"
)

func TestLocalBroker(t *testing.T) {
	suite.Run(t, new(LocalBrokerSuite))
}

type LocalBrokerSuite struct {
	suite.Suite
}

func (suite *LocalBrokerSuite) TestPublish_canceledContext() {
	broker := local.Broker()

	// Canceled explicitly
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	suite.Error(broker.Publish(ctx, "Foo", []byte("Hello")))

	// Canceled due to deadline
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	time.Sleep(2 * time.Millisecond)
	suite.Error(broker.Publish(ctx, "Foo", []byte("Hello")))
}

func (suite *LocalBrokerSuite) TestPublish_noSubscribers() {
	broker := local.Broker()
	suite.NoError(broker.Publish(context.Background(), "Foo", []byte("Hello")))
	suite.NoError(broker.Publish(context.Background(), "Bar", []byte("Goodbye")))
	suite.NoError(broker.Publish(context.Background(), "Baz", []byte("Seriously, go home.")))
}

func (suite *LocalBrokerSuite) publish(broker eventsource.Broker, key string, value string) {
	msg := "Publishing with a valid context, should always succeed"
	suite.Require().NoError(broker.Publish(context.Background(), key, []byte(value)), msg)
}

func (suite *LocalBrokerSuite) subscribe(broker eventsource.Broker, sequence *testext.Sequence, key string) eventsource.Subscription {
	subs, err := broker.Subscribe(key, func(ctx context.Context, evt *eventsource.EventMessage) error {
		if string(evt.Payload) == "error" {
			return fmt.Errorf("nope")
		}

		sequence.Append(fmt.Sprintf("%s:%s", key, string(evt.Payload)))
		return nil
	})
	suite.Require().NoError(err, "There shouldn't be any issues subscribing locally... ever.")
	return subs
}

func (suite *LocalBrokerSuite) subscribeGroup(broker eventsource.Broker, sequence *testext.Sequence, key string, group string, which string) eventsource.Subscription {
	subs, err := broker.SubscribeGroup(key, group, func(ctx context.Context, evt *eventsource.EventMessage) error {
		if string(evt.Payload) == "error" {
			return fmt.Errorf("nope")
		}

		sequence.Append(fmt.Sprintf("%s:%s:%s:%s", key, group, which, string(evt.Payload)))
		return nil
	})
	suite.Require().NoError(err, "There shouldn't be any issues subscribing locally... ever.")
	return subs
}

func (suite *LocalBrokerSuite) assertFired(sequence *testext.Sequence, expected []string) {
	// This sucks, but there's no other way for us to determine if all of the handlers have
	// finished their work. It's all small, in-memory lists, so this should be more than
	// enough time to be sure that the sequence contains the handler values.
	time.Sleep(25 * time.Millisecond)
	suite.ElementsMatch(expected, sequence.Values())
}

func (suite *LocalBrokerSuite) TestPublish_noMatching() {
	results := &testext.Sequence{}
	broker := local.Broker()
	suite.subscribe(broker, results, "Foo")
	suite.subscribe(broker, results, "Foo.Bar")

	suite.NoError(broker.Publish(context.Background(), "Foo.Foo", []byte("Hello")))
	suite.NoError(broker.Publish(context.Background(), "Bar", []byte("Goodbye")))
	suite.NoError(broker.Publish(context.Background(), "Baz", []byte("Seriously, go home.")))

	time.Sleep(10 * time.Millisecond)
	suite.Len(results.Values(), 0, "None of the event handlers should have fired")
}

func (suite *LocalBrokerSuite) TestPublish_matching() {
	results := &testext.Sequence{}
	broker := local.Broker()
	suite.subscribe(broker, results, "Foo")
	suite.subscribe(broker, results, "Bar")
	suite.subscribe(broker, results, "Foo.Bar")
	suite.subscribe(broker, results, "Foo.*")
	suite.subscribe(broker, results, "*")
	suite.subscribe(broker, results, "*.*")
	suite.subscribe(broker, results, "Foo.Bar.Goo")
	suite.subscribe(broker, results, "Foo.*.Goo")

	results.Reset()
	suite.publish(broker, "Foo", "A")
	suite.assertFired(results, []string{
		"Foo:A",
		"*:A",
	})

	results.Reset()
	suite.publish(broker, "Foo.Bar", "A")
	suite.publish(broker, "Foo.Bar", "B")
	suite.assertFired(results, []string{
		"Foo.Bar:A",
		"Foo.Bar:B",
		"Foo.*:A",
		"Foo.*:B",
		"*.*:A",
		"*.*:B",
	})

	results.Reset()

	suite.publish(broker, "Bar", "A")
	suite.publish(broker, "Bar.Baz", "B")
	suite.publish(broker, "Hello.World", "C")
	suite.publish(broker, "Foo.Bar.Goo", "D")
	suite.publish(broker, "Foo.Baz.Goo", "E")
	suite.publish(broker, "Nope.Nope.Nope", "F")
	suite.assertFired(results, []string{
		"Bar:A",
		"*:A",
		"*.*:B",
		"*.*:C",
		"Foo.Bar.Goo:D",
		"Foo.*.Goo:D",
		"Foo.*.Goo:E",
	})

	// Multiple subscribers to the same event should ALL get the event.
	results.Reset()
	suite.subscribe(broker, results, "Foo")
	suite.subscribe(broker, results, "Foo")
	suite.publish(broker, "Foo", "A")
	suite.assertFired(results, []string{
		"Foo:A",
		"Foo:A",
		"Foo:A",
		"*:A", // there's three explicit Foo handlers, but only one * handler
	})
}

// Ensure that the correct subscribers and groups fire when mixed together.
func (suite *LocalBrokerSuite) TestPublish_mixedGroups() {
	results := &testext.Sequence{}
	broker := local.Broker()

	// One group handling Foo.
	suite.subscribeGroup(broker, results, "Foo", "1", "")
	suite.subscribeGroup(broker, results, "Foo", "1", "")

	// Another group handling Foo.
	suite.subscribeGroup(broker, results, "Foo", "2", "")
	suite.subscribeGroup(broker, results, "Foo", "2", "")
	suite.subscribeGroup(broker, results, "Foo", "2", "")
	suite.subscribeGroup(broker, results, "Foo", "2", "")
	suite.subscribeGroup(broker, results, "Foo", "2", "")
	suite.subscribeGroup(broker, results, "Foo", "2", "")

	// Another generic group that should get every 1-token key.
	suite.subscribeGroup(broker, results, "*", "3", "")

	// Some one-off handlers that should get everything.
	suite.subscribe(broker, results, "Foo")
	suite.subscribe(broker, results, "Foo")

	suite.publish(broker, "Foo", "A")
	suite.publish(broker, "Foo", "B")
	suite.publish(broker, "Bar", "C")
	suite.publish(broker, "Foo.Bar", "D") // nothing matches this

	suite.assertFired(results, []string{
		// Group matches
		"Foo:1::A",
		"Foo:2::A",
		"*:3::A",
		"Foo:1::B",
		"Foo:2::B",
		"*:3::B",
		"*:3::C",

		// Individual listener matches
		"Foo:A",
		"Foo:A",
		"Foo:B",
		"Foo:B",
	})
}

func (suite *LocalBrokerSuite) TestPublish_groupRoundRobin() {
	results := &testext.Sequence{}
	broker := local.Broker()
	suite.subscribeGroup(broker, results, "Foo", "1", "0")
	suite.subscribeGroup(broker, results, "Foo", "1", "1")
	suite.subscribeGroup(broker, results, "*", "1", "2")

	results.Reset()
	suite.publish(broker, "Foo", "A")
	suite.publish(broker, "Foo", "B")
	suite.publish(broker, "Foo", "C")
	suite.publish(broker, "Foo", "D")
	suite.publish(broker, "Foo", "E")
	suite.publish(broker, "Foo", "F")
	suite.publish(broker, "Foo", "G")
	suite.assertFired(results, []string{
		"Foo:1:0:A",
		"Foo:1:1:B",
		"*:1:2:C",
		"Foo:1:0:D",
		"Foo:1:1:E",
		"*:1:2:F",
		"Foo:1:0:G",
	})
}

// Publishing should still work even if subscribers fail.
func (suite *LocalBrokerSuite) TestPublish_subscriberErrors() {
	results := &testext.Sequence{}
	broker := local.Broker(local.WithErrorHandler(func(err error) {
		results.Append("oops")
	}))

	suite.subscribeGroup(broker, results, "Foo", "1", "")
	suite.subscribeGroup(broker, results, "Foo", "1", "")
	suite.subscribeGroup(broker, results, "Foo", "2", "")
	suite.subscribe(broker, results, "Foo")
	suite.subscribe(broker, results, "Foo")
	suite.subscribe(broker, results, "*")

	// The value "error" is not actually added to the sequence and the handlers return a non-nil error.
	suite.publish(broker, "Foo", "error")
	suite.publish(broker, "Foo", "error")
	suite.publish(broker, "Bar", "error")
	suite.assertFired(results, []string{
		// Errors from publish #1
		"oops",
		"oops",
		"oops",
		"oops",
		"oops",
		// Errors from publish #2
		"oops",
		"oops",
		"oops",
		"oops",
		"oops",
		// Errors from publish #3 (only * matched)
		"oops",
	})
}

func (suite *LocalBrokerSuite) TestUnsubscribe() {
	results := &testext.Sequence{}
	broker := local.Broker(local.WithErrorHandler(func(err error) {
		results.Append("oops")
	}))

	s1 := suite.subscribe(broker, results, "Foo")
	s2 := suite.subscribeGroup(broker, results, "Foo", "1", "")
	s3 := suite.subscribe(broker, results, "*")

	suite.publish(broker, "Foo", "A")
	suite.assertFired(results, []string{
		"*:A",
		"Foo:A",
		"Foo:1::A",
	})

	results.Reset()
	suite.NoError(s1.Unsubscribe())
	suite.publish(broker, "Foo", "B")
	suite.assertFired(results, []string{
		"*:B",
		"Foo:1::B",
	})

	results.Reset()
	suite.NoError(s2.Unsubscribe())
	suite.publish(broker, "Foo", "C")
	suite.assertFired(results, []string{
		"*:C",
	})

	results.Reset()
	suite.NoError(s3.Unsubscribe())
	suite.publish(broker, "Foo", "D")
	suite.assertFired(results, []string{})

	// Should be able to add more back in after the fact
	suite.subscribe(broker, results, "Foo")
	suite.publish(broker, "Foo", "I'm Back!")
	suite.assertFired(results, []string{
		"Foo:I'm Back!",
	})
}
