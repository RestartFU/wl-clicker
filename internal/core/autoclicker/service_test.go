package autoclicker

import (
	"sync"
	"testing"
)

type recordingInjector struct {
	mu     sync.Mutex
	events []Event
	closed bool
}

func (r *recordingInjector) WriteEvents(events ...Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, events...)
	return nil
}

func (r *recordingInjector) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.closed = true
	return nil
}

func (r *recordingInjector) snapshot() []Event {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Event, len(r.events))
	copy(out, r.events)
	return out
}

func (r *recordingInjector) isClosed() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.closed
}

type noopLogger struct{}

func (noopLogger) Debug(string, ...any) {}
func (noopLogger) Info(string, ...any)  {}
func (noopLogger) Warn(string, ...any)  {}
func (noopLogger) Error(string, ...any) {}

func testConfig(startEnabled bool) Config {
	return Config{
		TriggerCode:    LeftButtonCode,
		ToggleCode:     LeftButtonCode + 1,
		TriggerSources: map[string]struct{}{"device": {}},
		ToggleSources:  map[string]struct{}{"device": {}},
		GrabSources:    map[string]struct{}{},
		GrabEnabled:    false,
		CPS:            10,
		StartEnabled:   startEnabled,
	}
}

func assertReleaseSuffix(t *testing.T, events []Event) {
	t.Helper()
	if len(events) < 2 {
		t.Fatalf("expected at least 2 events, got %d", len(events))
	}
	up := events[len(events)-2]
	syn := events[len(events)-1]
	if up != (Event{Type: EventTypeKey, Code: LeftButtonCode, Value: 0}) {
		t.Fatalf("unexpected release event: %#v", up)
	}
	if syn != (Event{Type: EventTypeSyn, Code: SynReportCode, Value: 0}) {
		t.Fatalf("unexpected sync event: %#v", syn)
	}
}

func TestSetEnabledFalseReleasesLeftButton(t *testing.T) {
	injector := &recordingInjector{}
	service, err := NewService(testConfig(true), injector, noopLogger{})
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}

	if err := service.writeEvents(
		Event{Type: EventTypeKey, Code: LeftButtonCode, Value: 1},
		Event{Type: EventTypeSyn, Code: SynReportCode, Value: 0},
	); err != nil {
		t.Fatalf("writeEvents() error = %v", err)
	}
	if !service.leftButtonDown.Load() {
		t.Fatalf("expected left button to be tracked as down")
	}

	service.SetEnabled(false)

	if service.leftButtonDown.Load() {
		t.Fatalf("expected left button to be tracked as up after disabling")
	}
	assertReleaseSuffix(t, injector.snapshot())
}

func TestStopReleasesLeftButtonBeforeClosingInjector(t *testing.T) {
	injector := &recordingInjector{}
	service, err := NewService(testConfig(true), injector, noopLogger{})
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}

	if err := service.writeEvents(
		Event{Type: EventTypeKey, Code: LeftButtonCode, Value: 1},
		Event{Type: EventTypeSyn, Code: SynReportCode, Value: 0},
	); err != nil {
		t.Fatalf("writeEvents() error = %v", err)
	}

	service.Stop()

	if !injector.isClosed() {
		t.Fatalf("expected injector to be closed")
	}
	assertReleaseSuffix(t, injector.snapshot())
}
