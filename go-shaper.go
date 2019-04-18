package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

func watch(interval time.Duration) (<-chan string, context.CancelFunc) {
	var (
		ch          = make(chan string)
		lastContent string
		ctx, cancel = context.WithCancel(context.Background())
	)

	go func() {
		defer close(ch)

		for {
			time.Sleep(interval)

			select {
			case <-ctx.Done():
				return
			default:
			}

			content, err := clipboard.ReadAll()
			if err != nil || content == "" || content == lastContent {
				continue
			}

			lastContent = content
			ch <- content
		}
	}()

	return ch, cancel
}

func main() {
	ch, _ := watch(200 * time.Millisecond)
	rep := strings.NewReplacer("\n", " ")

	for text := range ch {
		renew := rep.Replace(text)
		if err := clipboard.WriteAll(renew); err != nil {
			log.Fatal(err)
		}
	}
}
