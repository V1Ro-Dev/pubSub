package tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pubsub/internal/subpub"
)

func TestPubSubSubscribe(t *testing.T) {
	var tests = []struct {
		name      string
		subAmount int
		topics    []string
	}{
		{
			"1 Sub",
			1,
			[]string{"Random Topic"},
		},
		{
			"Multiple Subs",
			3,
			[]string{"Random Topic Amogus", "Random Topic Abobus", "Random Topic Aboba"},
		},
		{
			"Multiple Subs with same topic",
			3,
			[]string{"Random Topic Amogus"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ps := subpub.NewPubSub()
			var subs []subpub.Subscription

			for i := 0; i < test.subAmount; i++ {
				topic := test.topics[i%len(test.topics)]
				sub, err := ps.Subscribe(topic, func(msg interface{}) {
					time.Sleep(500 * time.Millisecond)
				})

				require.NoError(t, err)
				require.NotNil(t, sub)
				subs = append(subs, sub)
			}

			require.Equal(t, test.subAmount, len(subs))
		})
	}

}

func TestPubSubPublish(t *testing.T) {
	var tests = []struct {
		name        string
		topic       string
		subAmount   int
		message     string
		expectCount int
	}{
		{
			"1 sub receives message",
			"vk",
			1,
			"dadada tol'ko ne fake",
			1,
		},
		{
			"multiple subs receive message",
			"vk",
			3,
			"ping pong",
			3,
		},
		{
			"no subs",
			"empty",
			0,
			"foo bar",
			0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ps := subpub.NewPubSub()

			var (
				mu       sync.Mutex
				messages []string
			)

			for i := 0; i < test.subAmount; i++ {
				sub, err := ps.Subscribe(test.topic, func(msg interface{}) {
					time.Sleep(500 * time.Millisecond)
				})
				messages = append(messages, test.message)
				require.NoError(t, err)
				require.NotNil(t, sub)
			}

			err := ps.Publish(test.topic, test.message)
			require.NoError(t, err)

			time.Sleep(100 * time.Millisecond)

			mu.Lock()
			defer mu.Unlock()
			assert.Len(t, messages, test.expectCount)
			for _, m := range messages {
				assert.Equal(t, test.message, m)
			}
		})
	}
}

func TestPubSubUnsubscribe(t *testing.T) {
	var tests = []struct {
		name string
	}{
		{
			"default unsub",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ps := subpub.NewPubSub()

			sub, err := ps.Subscribe("events", func(msg interface{}) {
				time.Sleep(500 * time.Millisecond)
			})
			require.NoError(t, err)
			require.NotNil(t, sub)

			sub.Unsubscribe()
		})
	}
}

func TestPubSubClose(t *testing.T) {
	type testCase struct {
		name        string
		withTimeout bool
		topic       string
		message     string
	}

	tests := []testCase{
		{
			"close with timeout",
			true,
			"test test test",
			"work work work",
		},
		{
			"close with cancel",
			false,
			"",
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ps := subpub.NewPubSub()

			if test.withTimeout {
				sub, err := ps.Subscribe(test.topic, func(msg interface{}) {
					time.Sleep(500 * time.Millisecond)
				})
				require.NoError(t, err)
				require.NotNil(t, sub)

				err = ps.Publish(test.topic, test.message)
				require.NoError(t, err)

				ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
				defer cancel()

				err = ps.Close(ctx)
				require.ErrorIs(t, err, context.DeadlineExceeded)
			} else {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				err := ps.Close(ctx)
				require.NoError(t, err)
			}
		})
	}
}
