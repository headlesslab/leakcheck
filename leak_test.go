package leakcheck_test

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/headlesslab/leakcheck"
)

func TestAncestors(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	wait := make(chan int)

	go func() {
		<-wait
	}()

	go func() {
		go func() {
			go func() {
				<-wait
			}()
		}()

		time.Sleep(100 * time.Millisecond)

		id := leakcheck.Get(false)[0].GoroutineID

		ts := leakcheck.Get(true).Filter(leakcheck.IgnoreNonChildren())
		if len(ts) != 1 && ts[0].GoroutineID != id {
			t.Fail()
		}

		wg.Done()
	}()

	for _, trace := range leakcheck.Get(true) {
		if trace.GoroutineID != 1 && !trace.HasParent(1) {
			t.Fail()
		}
	}

	wg.Wait()
	close(wait)
}

func TestCheckLeak(t *testing.T) {
	t.Parallel()

	m := &mockT{out: bytes.NewBuffer(nil)}

	wait := make(chan int)
	go func() {
		<-wait
	}()

	leakcheck.CheckLeak(m, 100*time.Millisecond)

	if !m.Failed() {
		t.Log("m should fail")
		t.Fail()
	}

	if !strings.Contains(m.out.String(), "leaking goroutines") {
		t.Log("should find the leak")
		t.Fail()
	}

	close(wait)
}

func TestCheckLeakAlreadyFailed(_ *testing.T) {
	m := &mockT{failed: true}

	wait := make(chan int)
	go func() {
		<-wait
	}()

	leakcheck.CheckLeak(m, time.Millisecond)

	close(wait)
}

type mockT struct {
	failed bool
	out    *bytes.Buffer
}

func (m *mockT) Helper()      {}
func (m *mockT) Fail()        { m.failed = true }
func (m *mockT) Failed() bool { return m.failed }
func (m *mockT) Cleanup(f func()) {
	f()
}
func (m *mockT) Logf(format string, args ...interface{}) {
	fmt.Fprintf(m.out, format, args...)
}

func TestCheckMainLeak(t *testing.T) {
	old := leakcheck.Exit
	t.Cleanup(func() {
		leakcheck.Exit = old
		log.SetOutput(os.Stderr)
	})
	leakcheck.Exit = func(_ int) {}
	buf := bytes.NewBuffer(nil)
	log.SetOutput(buf)

	wait := make(chan int)
	go func() {
		<-wait
	}()

	leakcheck.CheckMainLeak(&mockM{}, time.Millisecond)

	if !strings.Contains(buf.String(), "leaking goroutines") {
		t.Fail()
	}

	close(wait)
}

func TestCheckMainLeakAlreadyFailed(t *testing.T) {
	old := leakcheck.Exit
	t.Cleanup(func() {
		leakcheck.Exit = old
		log.SetOutput(os.Stderr)
	})
	leakcheck.Exit = func(_ int) {}
	buf := bytes.NewBuffer(nil)
	log.SetOutput(buf)

	wait := make(chan int)
	go func() {
		<-wait
	}()

	leakcheck.CheckMainLeak(&mockM{code: 1}, time.Millisecond)

	t.Log(buf.String())

	if buf.String() != "" {
		t.Fail()
	}

	close(wait)
}

type mockM struct {
	code int
}

func (m *mockM) Run() int {
	return m.code
}
