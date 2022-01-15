package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

func main() {
	op := flag.String("op", "sum", "operation to be executed")
	col := flag.Int("col", 1, "column to execute operation upon")
	flag.Parse()

	if err := run(flag.Args(), *op, *col, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filenames []string, op string, column int, out io.Writer) error {
	var opFunc statsFunc

	if len(filenames) == 0 {
		return ErrNoFiles
	}
	if column < 1 {
		return ErrInvalidColumn
	}

	switch op {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg
	default:
		return ErrInvalidOperation
	}

	consolidate := make([]float64, 0)
	resChan := make(chan []float64)
	errChan := make(chan error)
	doneChan := make(chan struct{})
	filesChan := make(chan string)

	go func() {
		defer close(filesChan)
		for _, file := range filenames {
			filesChan <- file
		}
	}()

	wg := sync.WaitGroup{}
	//for _, fname := range filenames {
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for fname := range filesChan {
				log.Printf("worker %d: doing file: %s", i, fname)
				f, err := os.Open(fname)
				if err != nil {
					errChan <- err
					return
				}
				data, err := csvToFloat(f, column)
				if err != nil {
					errChan <- err
					return
				}
				if err := f.Close(); err != nil {
					errChan <- err
					return
				}
				resChan <- data
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(doneChan)
	}()
	//consolidate = append(consolidate, data...)
	for {
		select {
		case err := <-errChan:
			return err
		case data := <-resChan:
			consolidate = append(consolidate, data...)
		case <-doneChan:
			_, err := fmt.Fprintln(out, opFunc(consolidate))
			return err
		}
	}
	//_, err := fmt.Fprintln(out, opFunc(consolidate))
	//return err
}
