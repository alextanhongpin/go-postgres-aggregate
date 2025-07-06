package main

import (
	"context"
	"encoding/json"
	"fmt"
)

func testCounterRoot2() {
	counter := NewCounterRoot2()
	panicIfErr(counter.Increment(5))
	panicIfErr(counter.Decrement(2))
	fmt.Println(counter.root.State.Value) // Should print 3
	fmt.Println(counter.root.Events)

	counter2 := NewCounterRoot2()
	counter2.root.ApplyEvents(counter.root.Events)
	fmt.Println(counter2.root.State.Value) // Should also print 3
	fmt.Println(counter2.root.Events)
}

type repository interface {
	Load(ctx context.Context, aggregateType, aggregateID string) ([]Event, error)
	Save(ctx context.Context, aggregateType, aggregateID string, events []Event) error
}

type CounterRoot2 struct {
	root *Root[*Counter]
}

func (r *CounterRoot2) Increment(n int) error {
	b, err := json.Marshal(n)
	if err != nil {
		return err
	}
	return r.root.RaiseEvent(Event{
		Type: "increment",
		Data: b,
	})
}

func (r *CounterRoot2) Decrement(n int) error {
	b, err := json.Marshal(n)
	if err != nil {
		return err
	}
	return r.root.RaiseEvent(Event{
		Type: "decrement",
		Data: b,
	})
}

func NewCounterRoot2() *CounterRoot2 {
	return &CounterRoot2{
		root: &Root[*Counter]{
			Events: make([]Event, 0),
			State:  &Counter{},
			Type:   "Counter",
		},
	}
}

type root interface {
	Apply(event Event) error
}

type Root[T root] struct {
	Events []Event
	ID     string
	State  T
	Type   string
}

func (r *Root[T]) ApplyEvents(events []Event) error {
	for _, event := range events {
		if err := r.State.Apply(event); err != nil {
			return err
		}
	}

	return nil
}

func (r *Root[T]) RaiseEvent(event Event) error {
	if err := r.State.Apply(event); err != nil {
		return err
	}

	r.Events = append(r.Events, event)
	return nil
}
