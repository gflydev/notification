package notification

import (
	"io"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
)

// resetFactory clears the global handler registry so each test starts clean.
func resetFactory() {
	mu.Lock()
	handlerFactory = nil
	mu.Unlock()
}

func TestMain(m *testing.M) {
	// Silence library logging during tests.
	log.SetOutput(io.Discard)
	m.Run()
}

// ---- Test fixtures -----------------------------------------------------------

// IGreet is a sample notification interface.
type IGreet interface {
	Greet() string
}

// greetNotification is a sample notification payload implementing IGreet.
type greetNotification struct{ msg string }

func (g greetNotification) Greet() string { return g.msg }

// countingHandler records how many times Notify was called.
type countingHandler struct {
	notifiable any
	counter    *int32
}

func (h *countingHandler) Notify() { atomic.AddInt32(h.counter, 1) }

// panicHandler always panics when notified.
type panicHandler struct{}

func (h *panicHandler) Notify() { panic("boom") }

// ---- Tests -------------------------------------------------------------------

func TestSend_DisabledReturnsNil(t *testing.T) {
	resetFactory()
	// NOTIFICATION_ENABLE unset => disabled, Send is a no-op returning nil.
	t.Setenv("NOTIFICATION_ENABLE", "false")

	if err := Send(greetNotification{msg: "hi"}); err != nil {
		t.Fatalf("expected nil when disabled, got %v", err)
	}
}

func TestSend_NilNotification(t *testing.T) {
	resetFactory()
	t.Setenv("NOTIFICATION_ENABLE", "true")

	if err := Send(nil); !errors.Is(err, errors.InvalidParameter) {
		t.Fatalf("expected InvalidParameter for nil notification, got %v", err)
	}
}

func TestSend_NoMatchingHandler(t *testing.T) {
	resetFactory()
	t.Setenv("NOTIFICATION_ENABLE", "true")

	if err := Send(greetNotification{msg: "hi"}); !errors.Is(err, errors.NotImplemented) {
		t.Fatalf("expected NotImplemented when no handler matches, got %v", err)
	}
}

func TestSend_InvokesMatchingHandler(t *testing.T) {
	resetFactory()
	t.Setenv("NOTIFICATION_ENABLE", "true")

	var counter int32
	Register(func(n any) INotifiable {
		return &countingHandler{notifiable: n, counter: &counter}
	}, (*IGreet)(nil))

	if err := Send(greetNotification{msg: "hi"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := atomic.LoadInt32(&counter); got != 1 {
		t.Fatalf("expected handler invoked once, got %d", got)
	}
}

func TestSend_SkipsNonMatchingInterface(t *testing.T) {
	resetFactory()
	t.Setenv("NOTIFICATION_ENABLE", "true")

	type IOther interface{ Other() }

	var counter int32
	Register(func(n any) INotifiable {
		return &countingHandler{notifiable: n, counter: &counter}
	}, (*IOther)(nil))

	// greetNotification does not implement IOther, so no handler matches.
	if err := Send(greetNotification{msg: "hi"}); !errors.Is(err, errors.NotImplemented) {
		t.Fatalf("expected NotImplemented, got %v", err)
	}
	if got := atomic.LoadInt32(&counter); got != 0 {
		t.Fatalf("expected handler not invoked, got %d", got)
	}
}

func TestSend_RecoversFromHandlerPanic(t *testing.T) {
	resetFactory()
	t.Setenv("NOTIFICATION_ENABLE", "true")

	var counter int32
	// A panicking handler must not crash the process nor block other handlers.
	Register(func(n any) INotifiable { return &panicHandler{} }, (*IGreet)(nil))
	Register(func(n any) INotifiable {
		return &countingHandler{notifiable: n, counter: &counter}
	}, (*IGreet)(nil))

	if err := Send(greetNotification{msg: "hi"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := atomic.LoadInt32(&counter); got != 1 {
		t.Fatalf("expected surviving handler invoked once, got %d", got)
	}
}

func TestRegister_PanicsOnNilHandler(t *testing.T) {
	resetFactory()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on nil handler")
		}
	}()
	Register(nil, (*IGreet)(nil))
}

func TestRegister_PanicsOnNonInterfacePtr(t *testing.T) {
	resetFactory()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on non-interface pointer")
		}
	}()
	notAnInterface := 0
	Register(func(n any) INotifiable { return nil }, &notAnInterface)
}

func TestRegister_ConcurrentSafe(t *testing.T) {
	resetFactory()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Register(func(n any) INotifiable { return nil }, (*IGreet)(nil))
		}()
	}
	wg.Wait()

	mu.RLock()
	got := len(handlerFactory)
	mu.RUnlock()
	if got != 50 {
		t.Fatalf("expected 50 registered handlers, got %d", got)
	}
}
