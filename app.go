package main

import (
  "bufio"
  "sync"
  "time"
)

type App struct {
  wg       sync.WaitGroup
  scanner  *bufio.Scanner
  bundler  Bundler
  writer   Writer
}

func NewApp(scanner *bufio.Scanner, bundler Bundler, writer Writer) *App {
  return &App{
    scanner: scanner,
    bundler: bundler,
    writer:  writer,
  }
}

func (a *App) Run(autoTerminate bool) {
  defer a.wg.Wait()
  a.eventLoop(autoTerminate)
}

// eventLoop checks new messages in source and puts them in queue
func (a *App) eventLoop(autoTerminate bool) {
  for a.scanner.Scan() {
    text := a.scanner.Text()
    a.wg.Add(1)
    go a.processMessage(text)
  }
}

// processMessage delegates message processing to reader
func (a *App) processMessage(text string) {
  defer a.wg.Done()

  message, err := ParseMessage(text)
  if err != nil {
    return
  }

  complete := a.bundler.Bundle(message)
  if !complete {
    return
  }

  // Wait few milliseconds, just in case, if some logs are late
  time.Sleep(100 * time.Millisecond)
  tree, err := a.bundler.GetTraceTree(message.TraceID)
  a.writer.GetChan() <- tree
}
