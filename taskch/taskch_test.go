package taskch

import (
	"errors"
	"reflect"
	"testing"
)

func TestTaskch_New(t *testing.T) {
	if tch := New(); tch == nil {
		t.Error("unexpected nil")
	}
}

func TestTaskch_noError(t *testing.T) {
	tch := New()

	tch.Go(func() error { return nil })
	tch.Go(func() error { return nil })

	if err := tch.Wait(); err != nil {
		t.Error("unexpected error")
	}
}

func TestTaskch_singleError(t *testing.T) {
	tch := New()

	expErr := errors.New("expected error")

	tch.Go(func() error { return nil })
	tch.Go(func() error { return expErr })

	switch err := tch.Wait(); err {
	case expErr:
		// OK
	case nil:
		t.Error("unexpected success")
	default:
		t.Error("unexpected error", err)
	}
}

func TestTaskch_multipleErrors(t *testing.T) {
	tch := New()
	expErr := errors.New("expected error")

	tch.Go(func() error { return expErr })
	tch.Go(func() error { return expErr })

	switch err := tch.Wait(); err {
	case expErr:
		// OK
	case nil:
		t.Error("unexpected success")
	default:
		t.Error("unexpected error", err)
	}
}

func TestTaskch_processAllErrors(t *testing.T) {
	tch := New()

	err1 := errors.New("expected error 1")
	err2 := errors.New("expected error 2")
	err3 := errors.New("expected error 3")

	tch.Go(func() error { return err1 })
	tch.Go(func() error { return err2 })
	tch.Go(func() error { return err3 })

	errs := map[error]bool{}
	for err := tch.Wait(); err != nil; err = tch.Wait() {
		errs[err] = true
	}

	expected := map[error]bool{err1: true, err2: true, err3: true}
	if !reflect.DeepEqual(errs, expected) {
		t.Error("unexpected errors", errs)
	}
}
