package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"music_downloader_daemon/errors"
)

var (
	Connection         amqp.Connection
	Channel            amqp.Channel
	MusicRequestQueue  amqp.Queue
	MusicDownloadQueue amqp.Queue
)

const (
	hostname string = "localhost"
	port     int32  = 6572
	username string = ""
	password string = "password"
)

// Connect : Connects to RabbitMQ
func Connect() {
	connection, err := amqp.Dial("amqp://admin:jukebox_123!@212.224.88.157:5672/")
	errors.CaptureErr(err, "Failed to connect to RabbitMQ")
	Connection = *connection
	channel, err := Connection.Channel()
	errors.CaptureErr(err, "Failed to open a channel")
	Channel = *channel

	MusicRequestQueue, err = Channel.QueueDeclare(
		"music_request", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	errors.CaptureErr(err, "Failed to declare queue")
	MusicDownloadQueue, err = Channel.QueueDeclare(
		"music_download", // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	errors.CaptureErr(err, "Failed to declare queue")
	fmt.Println("Connected to RabbitMQ")
}

func Disconnect() {
	Connection.Close()
	Channel.Close()
}
