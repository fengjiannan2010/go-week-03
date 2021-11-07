package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	cctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g, ctx := errgroup.WithContext(cctx)
	g.Go(func() error {
		fmt.Println("start app")
		return serverApp(ctx)
	})
	g.Go(func() error {
		fmt.Println("start debug")
		return serverDebug(ctx)
	})

	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		s := <-c
		fmt.Println("Got signal:", s)
		cancel()
		return nil
	})

	err := g.Wait()
	fmt.Println(err)
	fmt.Println(ctx.Err())
}

func serverApp(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "Hello QCon!")
	})
	s := http.Server{Addr: ":8080", Handler: mux}
	go func() {
		<-ctx.Done()
		fmt.Println("stop app")
		s.Shutdown(context.Background())
	}()
	return s.ListenAndServe()
}

func serverDebug(ctx context.Context) error {
	mux := http.DefaultServeMux
	s := http.Server{Addr: ":9090", Handler: mux}
	go func() {
		<-ctx.Done()
		fmt.Println("stop debug")
		s.Shutdown(context.Background())
	}()
	return s.ListenAndServe()
}
