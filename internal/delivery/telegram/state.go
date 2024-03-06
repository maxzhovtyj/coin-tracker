package telegram

import "sync"

type FSM struct {
	usersState   map[int64]State
	usersStateMx sync.RWMutex
}

func NewFSM() *FSM {
	return &FSM{
		usersState: make(map[int64]State),
	}
}

type State struct {
	Caller string
	Step   string
	Data   any
}

func (fsm *FSM) Get(id int64) State {
	fsm.usersStateMx.RLock()
	defer fsm.usersStateMx.RUnlock()

	return fsm.usersState[id]
}

func (fsm *FSM) Update(id int64, s State) {
	fsm.usersStateMx.Lock()
	defer fsm.usersStateMx.Unlock()

	fsm.usersState[id] = s
}

func (fsm *FSM) Remove(id int64) {
	fsm.usersStateMx.Lock()
	defer fsm.usersStateMx.Unlock()

	delete(fsm.usersState, id)
}
