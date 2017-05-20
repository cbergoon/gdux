package gdux

import (
	"fmt"
	"time"
)

func Example() {

	st, err := NewStore(BasicState{Accumulator: 0}, BasicReducer)
	if err != nil {

	}

	st.WithOptions(Options{Record: true, RecordSize: 5})

	ch := make(chan State)
	quitChan := make(chan struct{}, 10)
	st.SubscribeChannel(ch)

	go func(s *Store) {
		for {
			select {
			case n := <-ch:
				fmt.Println("From Subscribor: ", n)
			case <-quitChan:
				return
			}
		}
	}(st)

	st.TriggerChan <- BasicAction{T: "INC"}
	st.TriggerChan <- BasicAction{T: "INC"}
	st.TriggerChan <- BasicAction{T: "INC"}
	st.TriggerChan <- BasicAction{T: "INC"}
	st.TriggerChan <- BasicAction{T: "DEC"}
	st.TriggerChan <- BasicAction{T: "DEC"}

	st.Trigger(BasicAction{T: "INC"})

	time.Sleep(time.Second)

	fmt.Println(st.actionRecord)

	//close(st.TriggerChan)
	quitChan <- struct{}{}

	//Output:
	//From Subscribor:  {1 INC}
	//From Subscribor:  {2 INC}
	//From Subscribor:  {3 INC}
	//From Subscribor:  {4 INC}
	//From Subscribor:  {3 DEC}
	//From Subscribor:  {2 DEC}
	//From Subscribor:  {3 INC}
	//[{INC} {INC} {DEC} {DEC} {INC}]

}
