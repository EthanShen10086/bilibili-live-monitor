package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/EthanShen10086/bilibili-live-monitor/go/internal/monitor"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	o, e := monitor.ParseOptions(os.Args[1:])
	if e == nil {
		e = monitor.CLI(ctx, o)
	}
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
