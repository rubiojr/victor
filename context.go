package victor

import (
    "math/rand"
    "time"
)

type User interface {
    Id() int
    Name() string
}

type Message interface {
    Id() int
    Type() string
    UserId() int
    RoomId() int
    Body() string
}

type Context struct {
    message Message

    Send  func(text string)
    Reply func(text string)
    Paste func(text string)
    Sound func(text string)

    matches []string
}

func (c *Context) SetMessage(msg Message) {
    c.message = msg
}

func (c *Context) Message() Message {
    return c.message
}

func (c *Context) SetMatches(matches []string) {
    c.matches = matches
}

func (c *Context) Matches() []string {
    return c.matches
}

func (c *Context) RandomString(strings []string) string {
    rand.Seed(time.Now().UTC().UnixNano())
    return strings[rand.Intn(len(strings))]
}
