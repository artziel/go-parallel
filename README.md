# Go Parallel
Artziel Narvaiza <artziel@gmail.com>

### Sample
```golang
package main

import GoParallel "github.com/artziel/go-parallel"

func main() {
    runner := GoParallel.Processor{}

    values := []int{}

    task := func(t *Concurrency.Worker) error {
        data := t.Data.(int)
        fmt.Printf("Value: %v\n", data)
    }

    for i = 0; i < 100; i++ {
		runner.AddWorker(&GoParallel.Worker{Data: i, Fnc: task})
    }

	err := runner.Run(25, func(progress int) {
		fmt.Println("All workers Finish!")
	})
}
```
