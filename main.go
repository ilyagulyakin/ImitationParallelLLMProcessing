package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"ImitationParallelLLMProcessing/agent"
	"ImitationParallelLLMProcessing/server"
)

const addr = ":9091"

type config struct {
	workers     int
	softTimeout time.Duration
	hardTimeout time.Duration
}

func main() {
	workers := flag.Int("workers", 2, "число воркеров")
	soft := flag.Duration("soft-shutdown", 10*time.Second, "мягкий дедлайн остановки")
	hard := flag.Duration("hard-shutdown", 10*time.Second, "жёсткий дедлайн остановки")
	flag.Parse()

	if *workers < 1 {
		log.Fatal("workers must be >= 1")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// После первого сигнала возвращаем обработку по умолчанию:
	// второй Ctrl+C завершит процесс сразу.
	go func() {
		<-ctx.Done()
		stop()
	}()

	cfg := config{
		workers:     *workers,
		softTimeout: *soft,
		hardTimeout: *hard,
	}

	if err := run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, cfg config) error {
	ch := make(chan agent.Ticket)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		Run(ch, cfg.workers)
	}()

	srv := server.NewServer(addr, ch)
	srvErr := make(chan error, 1)
	go func() { srvErr <- srv.Run() }()

	select {
	case err := <-srvErr:
		// Сервер не поднялся или остановился сам.
		if err != nil {
			return err
		}
		close(ch)
		wg.Wait()
		return nil
	case <-ctx.Done():
		log.Println("MAIN: получен сигнал, останавливаюсь")
	}

	shutdownStart := time.Now()

	// Мягкая остановка: закрываем listener, ждём активные хендлеры.
	softCtx, cancelSoft := context.WithTimeout(context.Background(), cfg.softTimeout)
	err := srv.Shutdown(softCtx)
	cancelSoft()

	// Жёсткая: хендлеры не успели, ждём ещё.
	if err != nil {
		log.Printf("MAIN: мягкий этап истёк (%v), перехожу к жёсткому", err)
		hardCtx, cancelHard := context.WithTimeout(context.Background(), cfg.hardTimeout)
		err = srv.Shutdown(hardCtx)
		cancelHard()
	}

	if err != nil {
		// close(ch) нельзя: хендлеры могут ещё писать в канал.
		log.Printf("MAIN: жёсткий дедлайн, потеряно в работе ~%d", srv.Accepted()-Processed())
		return err
	}

	close(ch)
	wg.Wait()
	log.Printf("MAIN: остановлен за %s, принято=%d, обработано=%d",
		time.Since(shutdownStart).Round(time.Millisecond), srv.Accepted(), Processed())
	return nil
}
