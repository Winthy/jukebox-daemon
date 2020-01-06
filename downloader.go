package main

import (
	"log"
	"music_downloader_daemon/errors"
	"music_downloader_daemon/rabbitmq"
	"music_downloader_daemon/upload"
	"os"
	"os/exec"

	"github.com/streadway/amqp"
)

func main() {
	rabbitmq.Connect()

	messages, err := rabbitmq.Channel.Consume(
		rabbitmq.MusicRequestQueue.Name, // queue
		"",                              // consumer
		true,                            // auto-ack
		false,                           // exclusive
		false,                           // no-local
		false,                           // no-wait
		nil,                             // args
	)
	errors.CaptureErr(err, "Failed to register a consumer")

	go upload.FileServer()

	for msg := range messages {
		id := string(msg.Body)
		url := "https://www.youtube.com/watch?v=" + id
		filePath := "cache/" + id + ".%(ext)s"
		log.Printf("Received a request with id: %s", id)
		if !fileExists("cache/" + id + ".mp3") {
			cmd := exec.Command("youtube-dl", "--extract-audio", "--output", filePath, "--audio-format", "mp3", url)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			errors.CaptureErr(err, "Failed to download audio")
		}
		err = rabbitmq.Channel.Publish(
			"",                               // exchange
			rabbitmq.MusicDownloadQueue.Name, // routing key
			false,                            // mandatory
			false,                            // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(id),
			})
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
