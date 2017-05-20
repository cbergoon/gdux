package gdux

import (
	"sync"

	"github.com/pkg/errors"
)

//Store represents the Gdux state container.
type Store struct {
	currentState     State        //currentState represents the state
	currentStateLock sync.RWMutex //currentStateLock covers: 'currentState'

	reducerFunction             Reducer                   //reducerFunction executed for each action
	channelSubscriptionRegistry []RegistrableStateChannel //channelSubscriptionRegistry list of channels to send updated state on

	options Options //options represents settings and configuration for store

	actionRecord []Action //actionRecord history of actions if enabled

	TriggerChan chan Action //triggerChan channel store receives actions
}

//Options represents options for the state.
type Options struct {
	Record     bool //Record indicates whether action history is stored
	RecordSize int  //RecordSize sets limit on action history; 0 for no limit
}

//NewStore creates a new store with the initial state value and reducer function. Starts go routine listening for action.
//Returns a pointer to the store and a possible error.
func NewStore(state State, reducer Reducer) (*Store, error) {
	if state == nil || reducer == nil {
		return nil, errors.New("Error: Cannot create store with nil state or reducer.")
	}
	store := &Store{
		currentState:    state,
		reducerFunction: reducer,
		TriggerChan:     make(chan Action),
		options:         Options{},
	}
	go func(s *Store) {
		for n := range s.TriggerChan {
			s.Trigger(n)
		}
	}(store)
	return store, nil
}

//WithOptions sets options on the store.
func (store *Store) WithOptions(options Options) {
	store.options = options
}

//GetState returns the current state.
func (store *Store) GetState() State {
	store.currentStateLock.RLock()
	defer store.currentStateLock.RUnlock()
	return store.currentState
}

/*
func (store *Store) SetState(state State) {
	store.currentStateLock.Lock()
	defer store.currentStateLock.Unlock()
	store.currentState = state
}
*/

//SubscriberChannel adds a channel to the list of channels that any updates to the state will be sent on.
func (store *Store) SubscribeChannel(registree RegistrableStateChannel) {
	store.registerSubscriptionChannel(registree)
}

//registerSubscriptionChannel internal only function to add a RegisterableStateChannel to the subscriber slice.
func (store *Store) registerSubscriptionChannel(registree RegistrableStateChannel) {
	store.channelSubscriptionRegistry = append(store.channelSubscriptionRegistry, registree)
}

//Trigger invokes a change in state with the provided action.
func (store *Store) Trigger(action Action) {
	store.trigger(action)
}

//triggerinternal only function that manages state changes.
func (store *Store) trigger(action Action) {
	store.currentStateLock.Lock()
	newState, err := store.reducerFunction(store.currentState, action)
	if err != nil {
		//no-op on state
	}
	store.currentState = newState
	if store.options.Record {
		store.actionRecord = append(store.actionRecord, action)
	}
	if store.options.RecordSize != 0 {
		for len(store.actionRecord) > store.options.RecordSize {
			store.actionRecord = store.actionRecord[1:]
		}
	}
	store.triggerNotifySubscribors()
	//store.currentStateLock.Unlock()

}

//triggerNotifySubscribers called by trigger and sends changes to the state to the subscribers.
func (store *Store) triggerNotifySubscribors() {
	x := store.currentState
	store.currentStateLock.Unlock()
	send := func(s State, sendChannel RegistrableStateChannel) {
		//store.currentStateLock.RLock()
		sendChannel <- s
		//store.currentStateLock.RUnlock()
	}
	for _, c := range store.channelSubscriptionRegistry {
		send(x, c)
	}
}

//NotifySubscribors sends changes to the state to the subscribers.
func (store *Store) NotifySubscribers() {
	store.currentStateLock.Lock()
	x := store.currentState
	store.currentStateLock.Unlock()
	send := func(s State, sendChannel RegistrableStateChannel) {
		//store.currentStateLock.RLock()
		sendChannel <- s
		//store.currentStateLock.RUnlock()
	}
	for _, c := range store.channelSubscriptionRegistry {
		send(x, c)
	}
}

//State is the interface that must be implemented for the state type.
type State interface {
	//Clone() State
}

//Actions is the interface all actions shouls implement.
type Action interface {
	GetType() string
}

//Reducer represents the signature for all reducer functions.
type Reducer func(state State, action Action) (State, error)

type RegistrableStateChannel chan State
type RegistrableActionChannel <-chan Action
