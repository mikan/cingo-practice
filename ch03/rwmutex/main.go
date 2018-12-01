package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

func main() {
	producer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		for i := 5; i > 0; i-- {
			l.Lock()
			l.Unlock()
			time.Sleep(1)
		}
	}
	observer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		l.Lock()
		defer l.Unlock()
	}
	test := func(count int, mutex, rwMutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count + 1)
		beginTestTime := time.Now()
		go producer(&wg, mutex)
		for i := count; i > 0; i-- {
			go observer(&wg, rwMutex)
		}
		wg.Wait()
		return time.Since(beginTestTime)
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer func() {
		err := tw.Flush()
		if err != nil {
			panic(err)
		}
	}()

	var m sync.RWMutex
	_, err := fmt.Fprintf(tw, "Reader\tRWMutex\tMutex\n")
	if err != nil {
		panic(err)
	}
	for i := 0; i < 22; i++ {
		count := int(math.Pow(2, float64(i)))
		_, err := fmt.Fprintf(
			tw,
			"%d\t%v\t%v\n",
			count,
			test(count, &m, m.RLocker()),
			test(count, &m, &m),
		)
		if err != nil {
			panic(err)
		}
	}
}
