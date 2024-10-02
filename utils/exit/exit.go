package exit

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	keepers []HouseKeeper
)

// HouseKeeper keeps the house clean.
type HouseKeeper func(os.Signal)

// Registry registry housekeeping functions.
func Registry(fn HouseKeeper) {
	keepers = append(keepers, fn)
}

// HouseKeeping delete resources before exit.
func HouseKeeping() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Printf("receiving sig %s\n", sig)
		for _, keeper := range keepers {
			keeper(sig)
		}
		os.Exit(0)
	}()
}
