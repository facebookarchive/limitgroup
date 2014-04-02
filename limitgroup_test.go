package limitgroup_test

import (
	"testing"
	"time"

	"github.com/facebookgo/limitgroup"
)

func TestLimitNotHit(t *testing.T) {
	lg := limitgroup.NewLimitGroup(2)
	lg.Add(1)
	lg.Add(1)
	fin := make(chan struct{})
	go func() {
		defer close(fin)
		lg.Wait()
	}()
	lg.Done()
	lg.Done()
	<-fin
}

func TestLimitHit(t *testing.T) {
	lg := limitgroup.NewLimitGroup(1)
	lg.Add(1)
	fin := make(chan struct{})
	go func() {
		defer close(fin)
		lg.Add(1)
	}()
	select {
	case <-fin:
		t.Fatal("should not get here")
	case <-time.After(time.Second):
	}
}

func TestNegativeAdd(t *testing.T) {
	lg := limitgroup.NewLimitGroup(1)
	lg.Add(1)
	fin := make(chan struct{})
	go func() {
		defer close(fin)
		lg.Wait()
	}()
	lg.Add(-1)
	<-fin
}

func expectPanic(t *testing.T, v interface{}) {
	r := recover()
	if r != v {
		t.Fatalf("expected panic %v got %v", v, r)
	}
}

func TestZeroLimit(t *testing.T) {
	defer expectPanic(t, "zero is not a valid limit")
	limitgroup.NewLimitGroup(0)
}

func TestDoneMoreThanPossible(t *testing.T) {
	defer expectPanic(t, "trying to return more slots than acquired")
	limitgroup.NewLimitGroup(1).Done()
}

func TestAddNegativeMoreThanExpected(t *testing.T) {
	defer expectPanic(t, "trying to return more slots than acquired")
	limitgroup.NewLimitGroup(1).Add(-1)
}

func TestAddMoreThanLimit(t *testing.T) {
	defer expectPanic(t, "delta greater than limit")
	limitgroup.NewLimitGroup(1).Add(2)
}

func TestAddZeroIsIgnored(t *testing.T) {
	limitgroup.NewLimitGroup(1).Add(0)
}
