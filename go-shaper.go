package main

import (
	"context"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

func main() {
	log.Println("running...")

	ch, _ := watch(200 * time.Millisecond)

	replacer := strings.NewReplacer("-\n", "", "- ", "", "\n", " ")
	if runtime.GOOS == "windows" {
		replacer = strings.NewReplacer("-\r\n", "", "- ", "", "\r\n", " ")
	}

	for text := range ch {
		renew := replacer.Replace(text)
		if err := clipboard.WriteAll(renew); err != nil {
			log.Fatal(err)
		}
	}
}

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
