package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const STARLINK_HOST_PORT = "192.168.100.1:9200"

type StarlinkLocationProvider struct {
	conn             *grpc.ClientConn
	descriptorSource grpcurl.DescriptorSource
}

func NewStarlinkLocationProvider() (LocationProvider, error) {
	conn, err := grpc.Dial(STARLINK_HOST_PORT, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	client := grpcreflect.NewClientAuto(ctx, conn)
	descriptorSource := grpcurl.DescriptorSourceFromServer(ctx, client)

	locationProvider := StarlinkLocationProvider{
		conn,
		descriptorSource,
	}
	return locationProvider, nil
}

func (p StarlinkLocationProvider) GetLocation() (Location, error) {
	ctx := context.Background()
	in := strings.NewReader(`{"get_location": {}}`)
	options := grpcurl.FormatOptions{}
	emptyLocation := Location{}

	rf, _, err := grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, p.descriptorSource, in, options)
	if err != nil {
		return emptyLocation, err
	}

	handler := ResponseHandler{}

	err = grpcurl.InvokeRPC(ctx, p.descriptorSource, p.conn, "SpaceX.API.Device.Device/Handle", []string{}, &handler, rf.Next)
	if err != nil {
		return emptyLocation, err
	}

	if handler.err != nil {
		return emptyLocation, handler.err
	}

	jsonString := handler.response
	starlinkData := struct {
		ApiVersion  string `json:"apiVersion"`
		GetLocation struct {
			Lla struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
				Alt float64 `json:"alt"`
			} `json:"lla"`
			Source string `json:"source"`
		} `json:"getLocation"`
	}{}
	err = json.Unmarshal([]byte(jsonString), &starlinkData)
	if err != nil {
		return emptyLocation, err
	}

	if starlinkData.ApiVersion != "10" {
		return emptyLocation, fmt.Errorf("unexpected api version: %v", starlinkData.ApiVersion)
	}

	return Location{
		Latitude:  starlinkData.GetLocation.Lla.Lat,
		Longitude: starlinkData.GetLocation.Lla.Lon,
	}, nil
}

type ResponseHandler struct {
	response string
	err      error
}

func (h *ResponseHandler) OnResolveMethod(md *desc.MethodDescriptor) {}
func (h *ResponseHandler) OnSendHeaders(md metadata.MD)              {}
func (h *ResponseHandler) OnReceiveHeaders(md metadata.MD)           {}

func (h *ResponseHandler) OnReceiveResponse(resp proto.Message) {
	marshaler := jsonpb.Marshaler{}
	jsonString, err := marshaler.MarshalToString(resp)
	if err != nil {
		h.err = err
		return
	}
	h.response = jsonString
}

func (h *ResponseHandler) OnReceiveTrailers(stat *status.Status, md metadata.MD) {}
