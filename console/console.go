package console

import (
	"fmt"
	"time"
)

var ConsoleInstance *Console = NewConsole()

type Console struct {
	Buffer     string
	BufferChan chan string
}

func NewConsole() *Console {
	return &Console{
		Buffer:     "",
		BufferChan: make(chan string, 100),
	}
}

func (c *Console) Start() {
	go func() {
		for {
			buf := <-c.BufferChan
			c.write(buf)
		}
	}()
}

func (c *Console) write(s string) {
	c.Buffer = fmt.Sprintf("%s: %s%s", time.Now().Format("2006-01-02 15:04:05"), s, c.Buffer)
}

func (c *Console) Write(s string) {
	s += "\n"
	c.BufferChan <- s
}

func (c *Console) Read() string {
	return c.Buffer
}

func (c *Console) Clear() {
	c.Buffer = ""
}
