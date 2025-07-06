package main

import (
	"encoding/json"
	"fmt"
)

func init() {
	fmt.Println("Counter Actions Example")
	actions := NewCounterActions()
	panicIfErr(actions.Increment(5))
	panicIfErr(actions.Decrement(2))
	fmt.Println(actions.Counter.Value) // Should print 3
	fmt.Println(actions.Events)

	actions2 := CounterActionsFromEvents(actions.Events)
	fmt.Println(actions2.Counter.Value) // Should also print 3
	fmt.Println("Counter Actions Example Completed")
}

type CounterAggregate struct {
	Value int
}

type CounterActions struct {
	Events  []Event
	Counter *CounterAggregate
}

func NewCounterActions() *CounterActions {
	return &CounterActions{
		Events:  make([]Event, 0),
		Counter: &CounterAggregate{},
	}
}

func CounterActionsFromEvents(events []Event) *CounterActions {
	actions := NewCounterActions()
	for _, event := range events {
		if err := actions.ApplyEvent(event); err != nil {
			panic(fmt.Sprintf("failed to apply event: %v", err))
		}
	}
	return actions
}

func (c *CounterActions) Increment(n int) error {
	b, err := json.Marshal(n)
	if err != nil {
		return err
	}

	return c.Raise(Event{
		Type: "increment",
		Data: b,
	})
}

func (c *CounterActions) Decrement(n int) error {
	b, err := json.Marshal(n)
	if err != nil {
		return err
	}

	return c.Raise(Event{
		Type: "decrement",
		Data: b,
	})
}

func (c *CounterActions) ApplyEvent(event Event) error {
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

func (c *CounterActions) Raise(event Event) error {
	if err := c.ApplyEvent(event); err != nil {
		return fmt.Errorf("failed to apply event: %w", err)
	}
	c.Events = append(c.Events, event)
	return nil
}
