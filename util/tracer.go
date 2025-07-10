package util

import (
	"fmt"
	"strings"
	"time"
)

type Entry struct {
	TimeMillis int64
	ID         int
}

type Tracer struct {
	enable  bool
	stack   []*Entry
	printer func(message string)
}

func NewTracer() *Tracer {
	return NewTracerWithPrinter(func(message string) {
		if message != "" {
			fmt.Println(message)
		}
	})
}

func NewTracerWithPrinter(printer func(message string)) *Tracer {
	t := &Tracer{
		enable:  false,
		printer: printer,
		stack:   []*Entry{},
	}
	t.stack = append(t.stack, &Entry{ID: 0, TimeMillis: time.Now().UnixMilli()})
	return t
}

func (t *Tracer) IsEnable() bool {
	return t.enable
}

func (t *Tracer) SetEnable(isTrace bool) {
	t.enable = isTrace
}

func (t *Tracer) Println(message string, args ...interface{}) {
	message = "[trace] " + message
	t.printTrace(message)
}

func (t *Tracer) StartTimer() {
	t.StartTimerWithMsg("")
}

func (t *Tracer) StartTimerWithMsg(message string, args ...interface{}) {
	if !t.enable {
		return
	}
	entry := t.newTimeEntry()
	t.stack = append(t.stack, entry)
	msg := t.makeBlank(entry.ID-1) + fmt.Sprintf("[trace%d]start ", entry.ID)
	if message != "" {
		msg += fmt.Sprintf(message, args...)
	}
	t.printTrace(msg)
}

func (t *Tracer) EndTimer(message string, args ...interface{}) {
	if !t.enable {
		return
	}
	if len(t.stack) == 0 {
		return
	}
	entry := t.stack[len(t.stack)-1]
	t.stack = t.stack[:len(t.stack)-1]
	timeUsed := time.Now().UnixMilli() - entry.TimeMillis
	msg := t.makeBlank(entry.ID-1) + fmt.Sprintf("[trace%d]end:%dms ", entry.ID, timeUsed)
	if message != "" {
		msg += fmt.Sprintf(message, args...)
	}
	t.printTrace(msg)

	if len(t.stack) == 0 {
		t.stack = append(t.stack, &Entry{ID: 0, TimeMillis: time.Now().UnixMilli()})
	}
}

func (t *Tracer) newTimeEntry() *Entry {
	lastId := t.stack[len(t.stack)-1].ID
	return &Entry{ID: lastId + 1, TimeMillis: time.Now().UnixMilli()}
}

func (t *Tracer) printTrace(message string) {
	if t.printer != nil {
		t.printer(message)
	}
}

func (t *Tracer) makeBlank(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(" ", n)
}
