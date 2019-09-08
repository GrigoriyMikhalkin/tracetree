package main

import (
  "bufio"
  "flag"
  "os"
  "sync"
)

func main() {
  var wg sync.WaitGroup

  inputFilePtr := flag.String("infile", "", "Input file")
  outputFilePtr := flag.String("outfile", "", "Output file")
  flag.Parse()

  infile, scanner := getScanner(*inputFilePtr)
  if infile != nil {
    defer infile.Close()
  }

  outfile, writer := getWriter(*outputFilePtr, &wg)
  if outfile != nil {
    defer outfile.Close()
  }

  bundler := NewBundler()
  app := NewApp(scanner, bundler, writer)

  wg.Add(1)
  defer wg.Wait()
  go writer.Run()
  app.Run()
  close(writer.GetChan())
}

func getScanner(path string) (file *os.File, scanner *bufio.Scanner) {
  if path == "" {
    scanner = bufio.NewScanner(os.Stdin)
    return
  }

  file, err := os.Open(path)
  if err != nil {
    panic(err)
  }

  scanner = bufio.NewScanner(file)
  return
}

func getWriter(path string, wg *sync.WaitGroup) (file *os.File, writer Writer) {
  if path == "" {
    writer = NewWriter(os.Stdout, wg)
    return
  }

  file, err := os.Create(path)
  if err != nil {
    panic(err)
  }

  writer = NewWriter(file, wg)
  return
}
