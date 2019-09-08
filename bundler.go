package main

import (
  "fmt"
  "sync"
)

type Bundler interface {
  Bundle(message *Message) (complete bool)
  GetTraceTree(traceId string) (tree *Tree, err error)
}

type bundlerImpl struct {
  // It's expected that there will be more reads than writes
  // so we use sync.map to avoid cache contention
  TraceServices *sync.Map
}

func NewBundler() *bundlerImpl{
  return &bundlerImpl{
    TraceServices: &sync.Map{},
  }
}

// Bundle bundles messages with same trace id
func (b *bundlerImpl) Bundle(message *Message) (complete bool) {
  complete = message.IsRoot()

  res, ok := b.TraceServices.Load(message.TraceID)
  if !ok {
    res, _ = b.TraceServices.LoadOrStore(message.TraceID, &Tree{
      Id:         message.TraceID,
      ServiceMap: &sync.Map{},
    })
  }

  tree := res.(*Tree)

  tree.UpdateServiceInfo(message)
  if !complete {
    tree.UpdateParentServiceInfo(message)
  }

  return
}

func (b *bundlerImpl) GetTraceTree(traceId string) (tree *Tree, err error) {
  res, ok := b.TraceServices.Load(traceId)
  if !ok {
    err = fmt.Errorf("Trace ID not found")
    return
  }

  tree = res.(*Tree)
  b.TraceServices.Delete(traceId)

  return
}
