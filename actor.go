package main

import (
	"encoding/json"
	"fmt"
)

func init() {
	fmt.Println("Counter Actor Example")
	actions := NewCounterActor()
	panicIfErr(actions.Increment(5))
	panicIfErr(actions.Decrement(2))
	fmt.Println(actions.Counter.Value) // Should print 3
	fmt.Println(actions.Events)

	actions2 := CounterActorFromEvents(actions.Events)
	fmt.Println(actions2.Counter.Value) // Should also print 3
	fmt.Println("Counter Actor Example Completed")
}

type CounterAggregate struct {
	Value int
}

type CounterActor struct {
	Events  []Event
	Counter *CounterAggregate
}

func NewCounterActor() *CounterActor {
	return &CounterActor{
		Events:  make([]Event, 0),
		Counter: &CounterAggregate{},
	}
}

func CounterActorFromEvents(events []Event) *CounterActor {
	actions := NewCounterActor()
	for _, event := range events {
		if err := actions.ApplyEvent(event); err != nil {
			panic(fmt.Sprintf("failed to apply event: %v", err))
		}
	}
	return actions
}

func (c *CounterActor) Increment(n int) error {
	return c.Raise("increment", n)
}

func (c *CounterActor) Decrement(n int) error {
	return c.Raise("decrement", n)
}

func (c *CounterActor) ApplyEvent(event Event) error {
	switch event.Type {
	case "increment":
		var n int
		if err := json.Unmarshal(event.Data, &n); err != nil {
			return err
		}
		c.Counter.Value += n
	case "decrement":
		var n int
		if err := json.Unmarshal(event.Data, &n); err != nil {
			return err
		}
		c.Counter.Value -= n
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
	return nil
}

func (c *CounterActor) Raise(name string, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	event := Event{
		Type: name,
		Data: b,
	}
	if err := c.ApplyEvent(event); err != nil {
		return fmt.Errorf("failed to apply event: %w", err)
	}
	c.Events = append(c.Events, event)
	return nil
}
