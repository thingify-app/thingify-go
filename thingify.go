package main

import (
	"fmt"
	"image"
	"time"

	thingrtc "github.com/thingify-app/thing-rtc-go"
	thingrtc_pairing "github.com/thingify-app/thing-rtc-go/pairing"
	schema "github.com/thingify-app/thingify-schema/golang"
	"google.golang.org/protobuf/proto"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/io/video"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/thingify-app/thing-rtc-go/codec/mmal"
	_ "github.com/thingify-app/thing-rtc-go/driver/camera"
)

const PAIRING_SERVER_URL = "https://thingify-xpo4wgiz5a-ts.a.run.app/pairing"
const SIGNALLING_SERVER_URL = "wss://thingify-xpo4wgiz5a-ts.a.run.app/signalling"

func main() {
	spi, err := InitSpi()
	if err != nil {
		panic(err)
	}
	defer spi.Close()

	pairing := thingrtc_pairing.NewPairing(PAIRING_SERVER_URL, "pairing.json")
	pairingIds := pairing.GetAllPairingIds()

	if len(pairingIds) == 0 {
		// Pairing does not exist, start pairing flow.
		fmt.Println("No pairings found, starting pairing flow.")
		doPairing(&pairing)

		// We should have succeeded at this point - repopulate pairingIds.
		pairingIds = pairing.GetAllPairingIds()
	}

	// Pairing should now exist (either existing or just setup above), so try connecting.
	fmt.Println("Attempting to connect to peer...")
	tokenGenerator, err := pairing.GetTokenGenerator(pairingIds[0])
	if err != nil {
		panic(err)
	}
	connect(spi, tokenGenerator)
}

func connect(spi *Spi, tokenGenerator thingrtc.TokenGenerator) {
	videoSource := thingrtc.CreateVideoMediaSource(640, 480)
	codec, err := mmal.NewCodec(1_000_000)
	if err != nil {
		panic(err)
	}

	peer := thingrtc.NewPeer(SIGNALLING_SERVER_URL, codec, videoSource)

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

	peer.Connect(tokenGenerator)
	defer peer.Disconnect()

	select {}
}

func parseCommand(bytes []byte) (schema.Command, error) {
	command := schema.Command{}
	err := proto.Unmarshal(bytes, &command)
	return command, err
}

func doPairing(pairing *thingrtc_pairing.Pairing) {
	stream, _ := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Video: func(constraint *mediadevices.MediaTrackConstraints) {
			constraint.Width = prop.Int(800)
			constraint.Height = prop.Int(600)
		},
	})

	track := stream.GetVideoTracks()[0].(*mediadevices.VideoTrack)
	defer track.Close()

	videoReader := track.NewReader(false)

	for {
		shortcode := findFrame(videoReader)
		fmt.Printf("Shortcode found: '%v', trying to respond to pairing...\n", shortcode)
		_, err := pairing.RespondToPairing(shortcode)
		if err != nil {
			fmt.Printf("Error responding to pairing: %v\n", err)
			continue
		} else {
			fmt.Printf("Pairing succeeded!\n")
			break
		}
	}
}

func findFrame(videoReader video.Reader) string {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	for range ticker.C {
		frame, release, err := videoReader.Read()
		if err != nil {
			fmt.Printf("Error reading frame: %v\n", err)
			continue
		}

		qrText, err := parseFrame(frame)
		release()
		if err != nil {
			fmt.Printf("Error parsing frame: %v\n", err)
			continue
		}

		return qrText
	}

	// Shouldn't get here because loop should not break.
	return ""
}

func parseFrame(frame image.Image) (string, error) {
	bmp, err := gozxing.NewBinaryBitmapFromImage(frame)
	if err != nil {
		fmt.Printf("Error converting frame: %v\n", err)
		return "", err
	}

	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		fmt.Printf("Error decoding frame: %v\n", err)
		return "", err
	}

	return result.GetText(), nil
}
