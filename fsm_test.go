package fsm

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	State State
}

func IsTestStructValid(e *Event) (bool, error) {
	return true, nil
}

func IsTestStructInvalid(e *Event) (bool, error) {
	return false, nil
}

func TestSetState(t *testing.T) {
	testStruct := &TestStruct{
		State: State("started"),
	}

	fsm := NewFSM()

	fsm.Register(reflect.TypeOf((*TestStruct)(nil)), "State", Events{{
		Name:  "make",
		From:  []State{"started"},
		To:    State("finished"),
		Guard: IsTestStructValid,
	}})

	fsm.Fire(testStruct, "make")
	if testStruct.State != State("finished") {
		t.Error("expected state to be 'finished'")
	}
}

func TestInvalidTransition(t *testing.T) {
	testStruct := &TestStruct{
		State: State("started"),
	}

	fsm := NewFSM()

	fsm.Register(reflect.TypeOf((*TestStruct)(nil)), "State", Events{{
		Name:  "make",
		From:  []State{"started"},
		To:    State("finished"),
		Guard: IsTestStructInvalid,
	}})

	err := fsm.Fire(testStruct, "make")

	if e, ok := err.(InvalidTransitionError); !ok && e.Event != "make" && e.State != "started" {
		t.Error("expected 'InvalidTransitionError'")
	}
}

func TestInvalidEvent(t *testing.T) {
	testStruct := &TestStruct{
		State: State("started"),
	}

	fsm := NewFSM()
	fsm.Register(reflect.TypeOf((*TestStruct)(nil)), "State", Events{{
		Name:  "make",
		From:  []State{"started"},
		To:    State("finished"),
		Guard: IsTestStructInvalid,
	}})

	err := fsm.Fire(testStruct, "some_event_name")

	if e, ok := err.(UnknownEventError); !ok && e.Event != "some_event_name" {
		t.Error("expected 'UnknownEventError'")
	}
}
