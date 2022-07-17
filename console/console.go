package console

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var (
	ConsoleInstance *Console = NewConsole()
	now                      = time.Now()
	filePath                 = fmt.Sprintf("logs/log-%v-%v-%v-%v.txt", now.Month(), now.Day(), now.Hour(), now.Minute())
)

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
			buf = fmt.Sprintf("%s: %s", time.Now().Format("2006-01-02 15:04:05"), buf)
			c.write(buf)
			c.writeToFile(filePath, buf)
		}
	}()
}

func (c *Console) write(s string) {
	c.Buffer = fmt.Sprintf("%s%s", s, c.Buffer)
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
func (c *Console) writeToFile(filePath, content string) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	write.WriteString(content)

	write.Flush()
}
