package main_test

import (
  "fmt"
  "testing"

  tracetree "github.com/GrigoriyMikhalkin/tracetree"
)

const (
  traceId = "2ovkwqzt"
  rootService = "rootService"

  firstMessage = "2013-10-23T10:12:35.021Z 2013-10-23T10:12:35.053Z 2ovkwqzt service1 k3zdpao7->5hun6fkq"
  secondMessage = "2013-10-23T10:12:36.021Z 2013-10-23T10:12:36.053Z 2ovkwqzt service2 k3zdpao7->fun76fkq"
  rootMessage = "2013-10-23T10:12:36.021Z 2013-10-23T10:12:36.053Z 2ovkwqzt rootService null->k3zdpao7"
)

func TestBundler(t *testing.T) {
  bundler := tracetree.NewBundler()

  // New trace id
  message, _ := tracetree.ParseMessage(firstMessage)
  complete := bundler.Bundle(message)
  if complete {
    t.Errorf("Bundle shouldn't be complete")
  }

  // Existing trace id
  message, _ = tracetree.ParseMessage(secondMessage)
  complete = bundler.Bundle(message)
  if complete {
    t.Errorf("Bundle shouldn't be complete")
  }

  // Root message
  message, _ = tracetree.ParseMessage(rootMessage)
  complete = bundler.Bundle(message)
  if !complete {
    t.Errorf("Bundle should be complete")
  }

  // Check trace tree
  tree, err := bundler.GetTraceTree(traceId)
  if err != nil {
    t.Errorf("Unexpected error, %s", err)
  }

  if tree.Id != traceId {
    t.Errorf("Unexpected traceId, %s", tree.Id)
  }

  if tree.Root.Service != rootService {
    t.Errorf("Unexpected root service, %s", tree.Root.Service)
  }

  expectedCalls := 2
  if len(tree.Root.Calls) != expectedCalls {
    t.Errorf(
      "Unexpected root service's number of calls, %d", len(tree.Root.Calls))
  }

  expectedTotalServiceCount := 3
  totalCount := 0
  tree.ServiceMap.Range(func(key, value interface{}) bool {
    totalCount++
    return true
  })
  if totalCount != expectedTotalServiceCount {
    t.Errorf("Unexpected total service count, %d", totalCount)
  }

  // Endure that trace id was deleted from TraceServices
  _, ok := bundler.TraceServices.Load(traceId)
  if ok {
    t.Errorf("%s entry should be deleted from  TraceServices", traceId)
  }

  fmt.Println("TestBundler finished!")
}
