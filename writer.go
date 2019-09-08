package main

import (
  "bufio"
  "encoding/json"
  "io"
  "sync"
)

const (
  chanSize = 1000
)

type Writer interface {
  GetChan() chan<- *Tree
  Run()
}

type writerImpl struct {
  wg        *sync.WaitGroup
  readChan  <-chan *Tree
  writeChan chan<- *Tree
  writer    *bufio.Writer
}

func NewWriter(output io.Writer, wg *sync.WaitGroup) (writer *writerImpl) {
  in, out := createInfiniteChan()

  writer = &writerImpl{
    wg:        wg,
    writeChan: in,
    readChan:  out,
    writer:    bufio.NewWriter(output),
  }

  return
}

func (w *writerImpl) write(tree *Tree) {
  b, _ := json.Marshal(tree)
  w.writer.Write(b)
  w.writer.WriteString("\n")
  w.writer.Flush()
}

func (w *writerImpl) GetChan() chan<- *Tree {
  return w.writeChan
}

func (w *writerImpl) Run() {
  defer w.wg.Done()

  for {
    select {
    case tree, ok := <- w.readChan:
      if !ok {
        return
      }
      w.write(tree)
    }
  }
}

// Based on https://medium.com/capital-one-tech/building-an-unbounded-channel-in-go-789e175cd2cd
func createInfiniteChan() (in chan *Tree, out chan *Tree) {
  in = make(chan *Tree)
  out = make(chan *Tree)

  go func() {
    var inQueue []*Tree
    outCh := func() chan *Tree {
      if len(inQueue) == 0 {
        return nil
      }
      return out
    }
    curVal := func() *Tree {
      if len(inQueue) == 0 {
        return nil
      }
      return inQueue[0]
    }

    for len(inQueue) > 0 || in != nil {
      select {
      case v, ok := <-in:
        if !ok {
          in = nil
        } else {
          inQueue = append(inQueue, v)
        }
      case outCh() <- curVal():
        inQueue = inQueue[1:]
      }
    }

    close(out)
  }()

  return in, out
}
