package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/petaki/inertia-go"
	"github.com/petaki/support-go/cli"
	"github.com/petaki/support-go/mix"
)

func Serve(debug bool, addr, url, redisUrl, redisKeyPrefix string) {
	infoLog := log.New(os.Stdout, cli.Cyan("INFO\t"), log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, cli.Red("ERROR\t"), log.Ldate|log.Ltime|log.Lshortfile)

	redisPool := newRedisPool(redisUrl)

	mixManager, inertiaManager, err := newMixAndInertiaManager(url)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &App{
		debug:          debug,
		url:            url,
		errorLog:       errorLog,
		infoLog:        infoLog,
		redisPool:      redisPool,
		redisKeyPrefix: redisKeyPrefix,
		mixManager:     mixManager,
		inertiaManager: inertiaManager,
	}

	srv := &http.Server{
		Addr:         addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	infoLog.Printf("Starting server on "+cli.Green("%s"), addr)

	go func() {
		err = srv.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			errorLog.Fatal(err)
		}
	}()

	<-done
	infoLog.Print("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		redisPool.Close()
		cancel()
	}()

	err = srv.Shutdown(ctx)
	if err != nil {
		errorLog.Fatal(err)
	}

	infoLog.Print("Server exited properly")
}

func newMixAndInertiaManager(url string) (*mix.Mix, *inertia.Inertia, error) {
	mixManager := mix.New("")

	version, err := mixManager.Hash("")
	if err != nil {
		return nil, nil, err
	}

	inertiaManager := inertia.New(url, "./resources/views/app.gohtml", version)

	icons, err := mixManager.Mix("images/bootstrap-icons.svg", "")
	if err != nil {
		return nil, nil, err
	}

	inertiaManager.Share("title", "Satellite")
	inertiaManager.Share("icons", icons)
	inertiaManager.ShareFunc("mix", mixManager.Mix)

	return mixManager, inertiaManager, nil
}

func newRedisPool(url string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(url)
		},
	}
}
