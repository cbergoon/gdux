package gdux

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
