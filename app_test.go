package main_test

import (
  "bufio"
  "fmt"
  "io"
  "testing"
  "time"

  tracetree "github.com/GrigoriyMikhalkin/tracetree"
)

const (
  validMessage = "2013-10-23T10:12:35.021Z 2013-10-23T10:12:35.053Z 2ovkwqzt service7 k3zdpao7->5hun6fkq\n"
  invalidMessage = "2013-10-23T10:12:35.021Z 2013-10-23T10:12:35.053Z 2ovkwqzt k3zdpao7->5hun6fkq\n"
)

func TestRun(t *testing.T) {
  r, w := io.Pipe()
  scanner := bufio.NewScanner(r)
  bundler := &fakeBundler{}
  writer := &fakeWriter{
    trees: make(chan *tracetree.Tree, 10),
  }
  app := tracetree.NewApp(scanner, bundler, writer)

  go app.Run()

  // Valid message, incomplete trace tree
  w.Write([]byte(validMessage))
  time.Sleep(1 * time.Second)

  // Invalid message
  w.Write([]byte(invalidMessage))
  time.Sleep(1 * time.Second)

  // Valid message, complete trace tree
  bundler.complete = true
  w.Write([]byte(validMessage))
  time.Sleep(1 * time.Second)

  expectedMessages := 2
  if len(bundler.messages) != expectedMessages {
    t.Errorf(
      "Expected %d messages, instead %d found",
      expectedMessages, len(bundler.messages))
  }

  expectedTrees := 1
  if len(writer.trees) != expectedTrees {
    t.Errorf(
      "Expected %d trees, instead %d found",
      expectedTrees, len(writer.trees))
  }

  fmt.Println("TestRun finished!")
}

type fakeBundler struct {
  complete bool
  messages []*tracetree.Message
}

func (b *fakeBundler) Bundle(message *tracetree.Message) (complete bool) {
  b.messages = append(b.messages, message)
  complete = b.complete
  return
}

func (b *fakeBundler) GetTraceTree(traceId string) (tree *tracetree.Tree, err error) {
  tree = &tracetree.Tree{
    Id: traceId,
  }

  return
}

type fakeWriter struct {
  trees chan *tracetree.Tree
}

func (w *fakeWriter) GetChan() chan<- *tracetree.Tree{
  return w.trees
}

func (w *fakeWriter) Run() {
  return
}
