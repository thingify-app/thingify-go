module github.com/thingify-app/thingify-go

go 1.16

require github.com/thingify-app/thing-rtc-go v0.0.0

require (
	github.com/thingify-app/thingify-schema/golang v0.0.0
	google.golang.org/protobuf v1.28.0
)

replace github.com/thingify-app/thing-rtc-go => ../thing-rtc-go

replace github.com/thingify-app/thingify-schema/golang => ../thingify-schema/golang