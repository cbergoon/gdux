package gdux

import (
	"testing"
	"time"
)

func Test(t *testing.T) {

	chanSends := []BasicAction{
		BasicAction{T: "INC"},
		BasicAction{T: "INC"},
		BasicAction{T: "INC"},
		BasicAction{T: "INC"},
		BasicAction{T: "DEC"},
		BasicAction{T: "DEC"},
	}

	triggerSends := []BasicAction{
		BasicAction{T: "INC"},
		BasicAction{T: "INC"},
		BasicAction{T: "INC"},
		BasicAction{T: "INC"},
		BasicAction{T: "DEC"},
		BasicAction{T: "DEC"},
	}

	st, err := NewStore(BasicState{Accumulator: 0}, BasicReducer)
	if err != nil {
		t.Errorf("failure: expected non nil error got %v", err)
	}

	st.WithOptions(Options{Record: true, RecordSize: 5})

	if !st.options.Record || st.options.RecordSize != 5 {
		t.Error("failure: invalid options")
	}

	ch := make(chan State)
	quitChan := make(chan struct{}, 10)
	st.SubscribeChannel(ch)

	testReceiveCount := 0
	go func(s *Store) {
		for {
			select {
			case <-ch:
				testReceiveCount = testReceiveCount + 1
			case <-quitChan:
				return
			}
		}
	}(st)

	for _, action := range chanSends {
		st.TriggerChan <- action
	} //testReceiveCount = 6

	for _, action := range triggerSends {
		st.Trigger(action)
	} //testReceiveCount = 12

	st.NotifySubscribers() //testReceiveCount = 13

	time.Sleep(time.Second) //cheating a little

	state := st.GetState().(BasicState)
	if state.LastAction != "DEC" || state.Accumulator != 4 {
		t.Error("failure: invalid state after triggers")
	}

	//close(st.TriggerChan)

	quitChan <- struct{}{}

	if testReceiveCount != 13 {
		t.Error("failure: received invalid number of updates")
	}

}
