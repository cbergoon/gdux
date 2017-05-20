## Gdux 
### A Redux Implementation Written in Golang

State container and management API; like Redux but in Go. Gdux's implementation of the state management concepts 
described below utilize channels to receive/trigger actions and also to notify subscribers. This is unique to Gdux and 
and allows it to be used in a wide variety of application types, most importantly concurrent and parallel applications. 
While channels provide great power and flexibility, Gdux can also be used without channels. 

#### Gdux Concept
Some applications have require a complex state management. In large applications this can become especially tedious to 
manitain; Gdux and the concepts implemented and established by [Redux](https://github.com/reactjs/redux) provide a 
single, immutable state container.
 
To paraphrase the Redux documentation, the state of an application using Gdux is stored in a single object. There should 
be only one store and thus one state representation or object in an application. All changes to the state are made by 
pure reducer functions which take as an argument an action. Actions are objects that describe how the Reducers should 
manipulate the state. Changes to the state are broadcast to listeners which the application can react to.

Below are explanations of Redux concepts and how they are implemented in Gdux. The following explanations and example 
implement a counter. 

##### State
State is the object representation of the state that the container is managing. This can be any struct that implements 
the State interface shown below. 
```go
type State interface {
}
```
This type should provide any values that you wish to track using Gdux as exported members. In the example below the state
is simply an integer called accumulator. 
```go
type ExampleState struct {
        Accumulator int
}
```
##### Action
An action contains the relevant information for the reducer to make a change to the state. Additionally, each action has a
'type' field to be used by the reducer as context of how to modify the state object. Actions may take different form but 
Gdux requires that they all implement the Action interface shown below. 
```go
type Action interface {
	GetType() string
}
```
The action below has only a 'type' as the state can only be modified by two actions in this example: increment or decrement. The 
action in this example is trivial but it is possible to include other data to describe the modification the action 
represents to the reducer. 
```go
type ExampleAction struct {
        T string
}

func (a *ExampleAction) GetType() string {
        return a.T
}
```

##### Reducer
Reducers are pure functions that take two arguments; the current state, and the action. The reducer operates on the state
and returns a new state object. The reducer is registered when the store is created and invoked when actions are triggered.  
In Gdux reducers must have the following function signature. 
```go
type Reducer func(state State, action Action) (State, error)
```
In this example the reducer has three possible deterministic paths. The first is 'INC' which increments the counter, the 
second is 'DEC' which decrements the counter, and the default makes no modification. Each path returns the new resulting 
state.
```go
func ExampleReducer(state State, action Action) (State, error) {
	if action.GetType() == "INC" {
		tmp := state.(ExampleState)
		tmp.Accumulator += 1
		return tmp, nil
	} else if action.GetType() == "DEC" {
		tmp := state.(ExampleState)
		tmp.Accumulator -= 1
		return tmp, nil
	} else {
		return state, nil
	}
}
```

#### Complete Example
```go
func BasicReducer(state State, action Action) (State, error) {
	if action.GetType() == "INC" {
		tmp := state.(ExampleState)
		tmp.Accumulator += 1
		return tmp, nil
	} else if action.GetType() == "DEC" {
		tmp := state.(ExampleState)
		tmp.Accumulator -= 1
		return tmp, nil
	} else {
		return state, nil
	}
}

func main() {

	st, err := NewStore(BasicState{Accumulator: 0}, BasicReducer) // Create a new store with an initial state and the reducer function. 
	if err != nil {
        errors.New("Something went wrong...")
	}

	ch := make(chan State) // Create a channel to receive updates to state on.
	quitChan := make(chan struct{}) 
	st.SubscribeChannel(ch) // Register the receive channel with the store. This can be done for as many channels as necessary.

	go func(s *Store) { // Start a go routine to listen for changes to state.
		for {
			select {
			case n := <-ch: // Receive on subscriber channel.
				fmt.Println("Reveived state on chan: ", n)
			case <-quitChan:
				return
			}
		}
	}(st)

	st.TriggerChan <- BasicAction{T: "INC"} // Send actions to TriggerChan to dispatch actions to store.
	st.TriggerChan <- BasicAction{T: "INC"} // ...
	st.TriggerChan <- BasicAction{T: "INC"} // ...
	st.TriggerChan <- BasicAction{T: "DEC"} // ...
	st.TriggerChan <- BasicAction{T: "INC"} // ...
	st.TriggerChan <- BasicAction{T: "DEC"} // ...
	st.TriggerChan <- BasicAction{T: "DEC"} // ...
	
	st.Trigger(BasicAction{T: "INC"}) // Dispatch another action via function call.
	fmt.Println(st.GetState()) // Print the state.

	quitChan <- struct{}{} // Tell our go routine above to exit.

}
```

#### License 

This project is licensed under the MIT License. See LICENSE file. 

#### Future Ideas
* History and ability to replay; great for testing/debug
* Combine reducers 
* Be able to register function listeners (need to understand what scope the function would run within)