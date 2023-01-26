package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	localconfig "desktop-notif-go/util"

	"github.com/go-redis/redis"
	"gopkg.in/toast.v1"
)

func main() {
	dir, _ := os.Getwd() // Get current directory

	// Setup log
	f, err := os.OpenFile("resources/log/main.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Load Config
	config, _ := localconfig.LoadConfig("resources/config.yml")
	if err != nil {
		log.Fatalln(err)
	}

	// Connect to Redis
	var redisClient = redis.NewClient(&redis.Options{
		Addr: config.Redis.Host + ":" + config.Redis.Port,
	})

	// check if redis is connected
	_, err = redisClient.Ping().Result()
	if err != nil {
		fmt.Println("Cannot connect to Redis")
	} else {
		fmt.Println("Redis is connected")
	}

	// Subscribe to channel
	subscriber := redisClient.Subscribe(config.Channel.Name)
	defer subscriber.Close()

	// Listen for messages
	for {
		msg, err := subscriber.ReceiveMessage()
		if err != nil {
			log.Fatalln(err)
		}

		// parse msg json
		var data map[string]interface{}
		json.Unmarshal([]byte(msg.Payload), &data)

		url := fmt.Sprintf("%s:%s%s%s", config.Frontend.Host, config.Frontend.Port, config.Frontend.Path, data["id"])

		notification := toast.Notification{
			AppID:   "Judul App",
			Title:   "Pemberitahuan",
			Message: data["id"].(string),
			Icon:    dir + "/resources/assets/warning.jpg",
			Audio:   toast.LoopingAlarm,
			Actions: []toast.Action{
				{Type: "protocol", Label: "Buka Aplikasi", Arguments: url},
			},
			ActivationArguments: url,
		}

		err = notification.Push()
		if err != nil {
			log.Printf("Error: %v | Payload: %s", err, msg.Payload)
		} else {
			log.Println(msg.Payload)
		}
	}
}
