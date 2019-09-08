package main

import (
  "fmt"
  "strings"
  "time"
)

const (
  spansCount = 2
  messageTextElements = 5
  referenceDateTimeFormat = "2006-01-02T15:04:05.000Z"

  RootParent = "null"
)

type Message struct {
  Start      string
  End        string
  TraceID    string
  Service    string
  ParentSpan string
  Span       string
}

func (m *Message) IsRoot() (root bool) {
  root = m.ParentSpan == RootParent

  return
}

/* Parses text into message. Expected text format:
 *   start_date end_date trace_id service_name parent_span->span
 *
 *  Example:
 *    2013-10-23T10:12:35.021Z 2013-10-23T10:12:35.053Z 2ovkwqzt service7 k3zdpao7->5hun6fkq
 *
 * If format is different, returns error
*/
func ParseMessage(text string) (message *Message, err error) {
  splittedText := strings.Split(text, " ")
  valid := validateMessage(splittedText)
  if !valid {
    err = fmt.Errorf("Message is invalid")
    return
  }

  spans := strings.Split(splittedText[4], "->")
  message = &Message{
    Start:      splittedText[0],
    End:        splittedText[1],
    TraceID:    splittedText[2],
    Service:    splittedText[3],
    ParentSpan: spans[0],
    Span:       spans[1],
  }

  return
}

func validateMessage(splittedText []string) (valid bool) {
  if len(splittedText) != messageTextElements {
    return
  }

  if _, err := time.Parse(referenceDateTimeFormat, splittedText[0]); err != nil {
    return
  }

  if _, err := time.Parse(referenceDateTimeFormat, splittedText[1]); err != nil {
    return
  }

  // Check that there 2 spans, and they're not the same
  spans := strings.Split(splittedText[4], "->")
  if len(spans) != spansCount {
    return
  } else if spans[0] == spans[1] {
    return
  }

  valid = true
  return
}
