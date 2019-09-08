package main

import (
  "sync"
)

type Tree struct {
  Id         string       `json:"id"`
  Root       *ServiceInfo `json:"root"`
  ServiceMap *sync.Map    `json:"-"`
}

type ServiceInfo struct {
  mx sync.Mutex

  Service string         `json:"service"`
  Start   string         `json:"start"`
  End     string         `json:"end"`
  Span    string         `json:"span"`
  Calls   []*ServiceInfo `json:"calls"`
}

type serviceMap struct {
  mx sync.Mutex

  services map[string]*ServiceInfo
}

func (t *Tree) UpdateServiceInfo(message *Message) {
  res, ok := t.ServiceMap.Load(message.Span)
  if !ok {
    res, _ = t.ServiceMap.LoadOrStore(
      message.Span,
      &ServiceInfo{Calls: make([]*ServiceInfo, 0)},
    )
  }

  serviceInfo := res.(*ServiceInfo)

  serviceInfo.mx.Lock()
  serviceInfo.Service = message.Service
  serviceInfo.Span = message.Span
  serviceInfo.Start = message.Start
  serviceInfo.End = message.End
  serviceInfo.mx.Unlock()

  // Check if root service
  if message.IsRoot() {
    t.Root = serviceInfo
  }
}

// UpdateParentServiceInfo should be always called after UpdateServiceInfo
func (t *Tree) UpdateParentServiceInfo(message *Message) {
  res, ok := t.ServiceMap.Load(message.ParentSpan)
  if !ok {
    res, _ = t.ServiceMap.
      LoadOrStore(
        message.ParentSpan,
        &ServiceInfo{Calls: make([]*ServiceInfo, 0)},
      )
  }
  parentServiceInfo := res.(*ServiceInfo)

  res, ok = t.ServiceMap.Load(message.Span)
  if !ok {
      // Means that UpdateServiceInfo wasn't called before
      panic("UpdateServiceInfo wasn't called")
  }
  serviceInfo := res.(*ServiceInfo)

  parentServiceInfo.mx.Lock()
  parentServiceInfo.Calls = append(parentServiceInfo.Calls, serviceInfo)
  parentServiceInfo.mx.Unlock()
}
