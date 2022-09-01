package main

import (
	"fmt"

	thingrtc "github.com/thingify-app/thing-rtc-go"
	schema "github.com/thingify-app/thingify-schema/golang"
	"google.golang.org/protobuf/proto"

	"github.com/thingify-app/thing-rtc-go/codec/mmal"
	_ "github.com/thingify-app/thing-rtc-go/driver/camera"
)

func main() {
	spi, err := InitSpi()
	if err != nil {
		panic(err)
	}
	defer spi.Close()

	tokenGenerator := thingrtc.BasicTokenGenerator{
		Role:        "responder",
		ResponderId: "123",
	}

	videoSource := thingrtc.CreateVideoMediaSource(640, 480)
	codec, err := mmal.NewCodec(1_000_000)
	if err != nil {
		panic(err)
	}

	peer := thingrtc.NewPeer("wss://thingify-test.herokuapp.com", codec, videoSource)

	peer.OnConnectionStateChange(func(connectionState int) {
		fmt.Printf("Connection state changed: %v", connectionState)
	})

	peer.OnBinaryMessage(func(message []byte) {
		cmd, err := parseCommand(message)
		if err != nil {
			fmt.Printf("Error parsing command: %v\n", err)
			return
		}
		fmt.Printf("Command received: %v, %v\n", cmd.ValueL, cmd.ValueR)

		err = spi.WritePwm(byte(cmd.ValueL), byte(cmd.ValueR))
		if err != nil {
			fmt.Printf("Error writing PWM: %v\n", err)
		}
	})

	err = peer.Connect(tokenGenerator)
	if err != nil {
		panic(err)
	}
	defer peer.Disconnect()

	select {}
}

func parseCommand(bytes []byte) (schema.Command, error) {
	command := schema.Command{}
	err := proto.Unmarshal(bytes, &command)
	return command, err
}
