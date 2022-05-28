package main

import (
	"fmt"

	thingrtc "github.com/thingify-app/thing-rtc-go"
)

func main() {
	tokenGenerator := thingrtc.BasicTokenGenerator{
		Role:        "responder",
		ResponderId: "123",
	}

	peer := thingrtc.NewPeer("wss://thingify-test.herokuapp.com")

	peer.OnConnectionStateChange(func(connectionState int) {
		fmt.Printf("Connection state changed: %v", connectionState)
	})

	peer.OnStringMessage(func(message string) {
		fmt.Printf("String message received: %v\n", message)
	})
	peer.OnBinaryMessage(func(message []byte) {
		fmt.Printf("Binary message received: %v\n", message)
	})

	err := peer.Connect(tokenGenerator)
	if err != nil {
		panic(err)
	}
	defer peer.Disconnect()

	select {}
}
