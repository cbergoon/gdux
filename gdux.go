package gdux

import (
	"sync"

	"github.com/pkg/errors"
)

type Store struct {
	currentState     State
	currentStateLock sync.RWMutex //covers: 'currentState'

	reducerFunction             Reducer
	channelSubscriptionRegistry []RegistrableStateChannel

	TriggerChan chan Action
	//quitChan    chan struct{}
}

func NewStore(state State, reducer Reducer) (*Store, error) {
	if state == nil || reducer == nil {
		return nil, errors.New("Error: Cannot create store with nil state or reducer.")
	}
	store := &Store{
		currentState:    state,
		reducerFunction: reducer,
		TriggerChan:     make(chan Action),
		//quitChan:        make(chan struct{}),
	}
	go func(s *Store) {
		for n := range s.TriggerChan {
			s.Trigger(n)
		}
	}(store)
	return store, nil
}

func (store *Store) GetState() State {
	store.currentStateLock.RLock()
	defer store.currentStateLock.RUnlock()
	return store.currentState
}

func (store *Store) SetState(state State) {
	store.currentStateLock.Lock()
	defer store.currentStateLock.Unlock()
	store.currentState = state
}

//Register Subscription Channel
func (store *Store) SubscribeChannel(registree RegistrableStateChannel) {
	store.registerSubscriptionChannel(registree)
}

func (store *Store) registerSubscriptionChannel(registree RegistrableStateChannel) {
	store.channelSubscriptionRegistry = append(store.channelSubscriptionRegistry, registree)
}

//todo (cbergoon): need something to listen on these listener channels. fan-in mechanism. needs to be in order received.

//Triggers
func (store *Store) Trigger(action Action) {
	store.trigger(action)
}

func (store *Store) trigger(action Action) {
	store.currentStateLock.Lock()
	newState, err := store.reducerFunction(store.currentState, action)
	if err != nil {
		//no-op on state
	}
	store.currentState = newState

	store.currentStateLock.Unlock()
	store.NotifySubscribers()
}

func (store *Store) NotifySubscribers() {
	send := func(s *Store, sendChannel RegistrableStateChannel) {
		store.currentStateLock.RLock()
		sendChannel <- s.currentState
		store.currentStateLock.RUnlock()
	}
	for _, c := range store.channelSubscriptionRegistry {
		send(store, c)
	}
}

type State interface {
	Clone() State
}

type Action interface {
	GetType() string
}

type Reducer func(state State, action Action) (State, error)

//type RegistrableFunction func(state State)//thought (might not want or need these)
type RegistrableStateChannel chan State
type RegistrableActionChannel <-chan Action

func BasicReducer(state State, action Action) (State, error) {
	if action.GetType() == "INC" {
		tmp := state.(BasicState)
		tmp.Accumulator += 1
		tmp.LastAction = "INC"
		return tmp, nil
	} else if action.GetType() == "DEC" {
		tmp := state.(BasicState)
		tmp.Accumulator -= 1
		tmp.LastAction = "DEC"
		return tmp, nil
	} else {
		return state, nil
	}
}

type BasicAction struct {
	T string
}

func (action BasicAction) GetType() string {
	return action.T
}

type BasicState struct {
	Accumulator int
	LastAction  string
}

func (state BasicState) Clone() State {
	return state
}
