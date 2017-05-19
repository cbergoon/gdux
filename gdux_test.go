package gdux

import (
	"fmt"
	"testing"
)

func TestNewStore(t *testing.T) {

	st, err := NewStore(BasicState{Accumulator: 0}, BasicReducer)
	if err != nil {

	}

	st.SetState(BasicState{Accumulator: 0})

	ch := make(chan State)
	quitChan := make(chan struct{})
	st.SubscribeChannel(ch)

	go func(s *Store) {
		for {
			select {
			case n := <-ch:
				fmt.Println("from chan ", n)
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

	//time.Sleep(time.Second)

	//close(st.TriggerChan)
	quitChan <- struct{}{}

}
