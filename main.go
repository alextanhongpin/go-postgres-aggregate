package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	root := NewCounterRoot()
	root.Counter.ID = 1
	panicIfErr(root.Increment(5))
	fmt.Println(root.Counter.Value)
	panicIfErr(root.Increment(3))
	fmt.Println(root.Counter.Value)
	fmt.Println(root.Events)

	repo := NewInMemoryCounterRepository()

	panicIfErr(repo.Save(root))

	{
		root, err := repo.Load(1)
		panicIfErr(err)
		fmt.Println(root.Counter.Value)
		fmt.Println(root.Events)
		panicIfErr(root.Increment(2))
		fmt.Println(root.Events)
	}

	testCounterRoot2()
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

type CounterRepository interface {
	Load(id int) (*CounterRoot, error)
	Save(root *CounterRoot) error
}

type InMemoryCounterRepository struct {
	Counters map[int][]Event
}

func NewInMemoryCounterRepository() *InMemoryCounterRepository {
	return &InMemoryCounterRepository{
		Counters: make(map[int][]Event),
	}
}
func (r *InMemoryCounterRepository) Load(id int) (*CounterRoot, error) {
	if events, exists := r.Counters[id]; exists {
		counter := NewCounterRoot()
		if err := counter.ApplyEvents(events); err != nil {
			return nil, fmt.Errorf("failed to apply events for counter with id %d: %w", id, err)
		}
		counter.Counter.ID = id
		return counter, nil
	}

	return nil, fmt.Errorf("counter with id %d not found", id)
}

func (r *InMemoryCounterRepository) Save(root *CounterRoot) error {
	if root.Counter.ID == 0 {
		return fmt.Errorf("counter ID is not set")
	}
	r.Counters[root.Counter.ID] = append(r.Counters[root.Counter.ID], root.Events...)
	return nil
}

type CounterRoot struct {
	Events  []Event
	Counter *Counter
}

func NewCounterRoot() *CounterRoot {
	return &CounterRoot{
		Events:  []Event{},
		Counter: &Counter{Value: 0},
	}
}

func (r *CounterRoot) ApplyEvents(events []Event) error {
	for _, event := range events {
		if err := r.Counter.Apply(event); err != nil {
			return err
		}
	}
	return nil
}

func (r *CounterRoot) Increment(n int) error {
	b, err := json.Marshal(n)
	if err != nil {
		return err
	}
	return r.Raise(Event{
		Type: "increment",
		Data: b,
	})
}

func (r *CounterRoot) Decrement(n int) error {
	b, err := json.Marshal(n)
	if err != nil {
		return err
	}
	return r.Raise(Event{
		Type: "decrement",
		Data: b,
	})
}

func (r *CounterRoot) Raise(event Event) error {
	if err := r.Counter.Apply(event); err != nil {
		return err
	}
	r.Events = append(r.Events, event)
	return nil
}

type Event struct {
	Type string
	Data json.RawMessage
}

type Counter struct {
	ID    int
	Value int
}

func (c *Counter) Apply(event Event) error {
	switch event.Type {
	case "increment":
		var inc int
		if err := json.Unmarshal(event.Data, &inc); err != nil {
			return err
		}
		c.Value += inc
	case "decrement":
		var dec int
		if err := json.Unmarshal(event.Data, &dec); err != nil {
			return err
		}
		c.Value -= dec
	default:
		return nil // Unknown event type, do nothing
	}
	return nil
}
