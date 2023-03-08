package console

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/monopolly/errors"
	"github.com/monopolly/numbers"
)

func New(notime ...bool) (res *Log) {

	res = new(Log)

	if len(notime) == 0 {
		res.minute = Black(time.Now().Format("2 Jan 15:04"))
		go res.minutes()
	}

	go res.worker()
	return
}

type Item struct {
	Time time.Time
	Type int
	Body interface{}
}

func (a *Item) HumanLog() string {
	var types string
	switch a.Type {
	case OK:
		types = "OK"
	case Info:
		types = "Info"
	}
	return fmt.Sprintf("[%s] %s: %v", a.Time.Format("15:04"), types, a.Body)
}

type Log struct {
	handler      []func(Item) //handler
	errorhandler []func(errors.E)
	minute       string
	mute         bool
	stream       chan Item
	queue        struct {
		sync.RWMutex
		list []Item
	}
}

func (a *Log) Mute(v bool) *Log {
	a.mute = v
	return a
}

func (a *Log) worker() {
	for {
		select {
		case item := <-a.stream:
			a.queue.list = append(a.queue.list, item)
		}
	}
}

func (a *Log) queueworker() {
	for {
		switch len(a.queue.list) > 0 {
		case true:
			a.queue.Lock()
			list := a.queue.list
			a.queue.list = nil
			for _, x := range list {
				switch x.Type {
				case OK:
				case Info:
				case Error:

				}
				//println(a.minute, Green(a.format(v...)))
			}
			a.queue.Unlock()
		case false:

		}

	}
}

// берет текущую минуту
func (a *Log) minutes() {
	for {
		a.minute = Black(time.Now().Format("2 Jan 15:04")) //Black(time.Now().Format("15:04"))
		time.Sleep(time.Minute)
	}
}

// ставит хендлер когда ловим ошибку
// нужен чтобы обработать глобально например отправить куда нить
// в хранилище логов
func (a *Log) AddHandler(handler func(Item)) {
	a.handler = append(a.handler, handler)
}

func (a *Log) AddErrorHandler(handler func(errors.E)) {
	a.errorhandler = append(a.errorhandler, handler)
}

func (a *Log) handleerror(err errors.E) {
	for _, h := range a.errorhandler {
		h(err)
	}
}

func (a *Log) handle(v Item) {
	for _, h := range a.handler {
		h(v)
	}
}

// форматирование строки [1, "nice", true] в текст "1, nice, true"
func (a *Log) format(v ...interface{}) string {
	var lines []string
	for _, s := range v {
		l := fmt.Sprint(s)
		l = strings.TrimSpace(l)
		if l != "" {
			lines = append(lines, l)
		}
	}
	return strings.Join(lines, ", ")
}

func (a *Log) OK(v ...interface{}) *Log {
	if len(v) == 0 {
		return a
	}
	if a.mute {
		return a
	}
	println(a.minute, Green(a.format(v...)))
	go a.handle(Item{
		Time: time.Now(),
		Type: OK,
		Body: a.format(v...),
	})
	return a
}

func (a *Log) Green(v ...interface{}) (res string) {
	return Green(a.format(v...))
}

func (a *Log) Status(name string, v ...interface{}) *Log {
	if a.mute {
		return a
	}

	lines := []string{
		fmt.Sprintf("%s %s", name, strings.Repeat(".", 30-len(name))),
	}
	for _, x := range v {
		lines = append(lines, fmt.Sprintf("%v", x))
	}

	println(a.minute, Green(strings.Join(lines, " ")))
	go a.handle(Item{
		Time: time.Now(),
		Type: OK,
		Body: a.format(v...),
	})
	return a
}

func (a *Log) StatusError(name string, v ...interface{}) *Log {
	if a.mute {
		return a
	}

	lines := []string{
		fmt.Sprintf("%s %s", name, strings.Repeat(".", 30-len(name))),
	}
	for _, x := range v {
		lines = append(lines, fmt.Sprintf("%v", x))
	}

	println(a.minute, Red(strings.Join(lines, " ")))
	go a.handle(Item{
		Time: time.Now(),
		Type: OK,
		Body: a.format(v...),
	})
	return a
}

func (a *Log) Size(size int, count int) *Log {
	a.Status(numbers.Formats(count), humanize.Bytes(uint64(size*count)))
	return a
}

func (a *Log) OKf(text string, v ...interface{}) *Log {
	if a.mute {
		return a
	}
	a.OK(fmt.Sprintf(text, v...))
	return a
}

func (a *Log) Println(v ...interface{}) {
	a.Info(v...)
}

func (a *Log) Info(v ...interface{}) *Log {
	if a.mute {
		return a
	}
	println(a.minute, a.format(v...))
	go a.handle(Item{time.Now(), Info, a.format(v...)})
	return a
}

func (a *Log) Bytes(v ...[]byte) *Log {
	if a.mute {
		return a
	}
	go println(a.minute, a.format(string(bytes.Join(v, []byte(", ")))))
	go a.handle(Item{time.Now(), Info, a.format(string(bytes.Join(v, []byte(", "))))})
	return a
}

func (a *Log) Since(t time.Time, title ...string) *Log {
	if a.mute {
		return a
	}
	switch len(title) > 0 {
	case true:
		go a.OK(title[0] + " " + time.Since(t).String())
	default:
		go a.OK(time.Since(t).String())
	}

	return a
}

func (a *Log) Now() time.Time {
	return time.Now()
}

func (a *Log) Time() {
	if a.mute {
		return
	}
	go println(a.minute, fmt.Sprint(time.Now().Unix()))
}

func (a *Log) Unix() int64 {
	return time.Now().Unix()
}
func (a *Log) UnixNano() int64 {
	return time.Now().UnixNano()
}
func (a *Log) TimeNano() {
	if a.mute {
		return
	}
	go println(a.minute, fmt.Sprint(time.Now().UnixNano()))
}

func (a *Log) Json(v interface{}) *Log {
	if a.mute {
		return a
	}
	b, _ := json.MarshalIndent(v, "", "    ")
	go fmt.Println(string(b))
	return a
}

func (a *Log) Printf(text string, v ...interface{}) {
	if a.mute {
		return
	}
	go println(a.minute, Red(fmt.Printf(text, v...)))
	go a.Infof(text, v...)
}

func (a *Log) Infof(text string, v ...interface{}) *Log {
	if a.mute {
		return a
	}
	go a.Info(fmt.Sprintf(text, v...))
	return a
}

func (a *Log) Play(text string, v ...interface{}) *Log {
	if a.mute {
		return a
	}
	go fmt.Printf(fmt.Sprintf("\r%s %s", a.minute, text), v...)
	return a
}

func (a *Log) PlayNum(v interface{}) *Log {
	if a.mute {
		return a
	}
	go fmt.Printf("\r%s %v", a.minute, v)
	return a
}

func Play(text string, v ...interface{}) {
	fmt.Printf(fmt.Sprintf("\r %s", text), v...)
}

func PlayNum(v interface{}) {
	fmt.Printf("\r%v", v)
}

func (a *Log) Err(v ...interface{}) *Log {
	if a.mute {
		return a
	}
	go println(a.minute, Red(v))
	go a.handleerror(errors.Unknown(v))
	return a
}

func (a *Log) Error(e errors.E) *Log {
	if e == nil {
		return a
	}
	if a.mute {
		return a
	}
	go println(a.minute, Red(e.Error()))
	go a.handleerror(e)
	return a
}

func (a *Log) ErrorE(err error, some ...interface{}) {
	if a.mute {
		return
	}
	some = append(some, err)
	if err != nil {
		println(a.minute, Red(some...))
	}
	a.handleerror(errors.Unknown(err))
}

func (a *Log) Errorf(s string, v ...interface{}) *Log {
	if a.mute {
		return a
	}
	println(a.minute, Red(fmt.Sprintf(s, v...)))
	return a
}

// return the source filename after the last slash
func chopPath(original string) string {
	i := strings.LastIndex(original, "/")
	if i == -1 {
		return original
	} else {
		return original[i+1:]
	}
}
