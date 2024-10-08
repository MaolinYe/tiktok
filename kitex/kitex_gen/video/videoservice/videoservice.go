// Code generated by Kitex v0.11.3. DO NOT EDIT.

package videoservice

import (
	"context"
	"errors"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	video "tiktok/kitex/kitex_gen/video"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"Feed": kitex.NewMethodInfo(
		feedHandler,
		newVideoServiceFeedArgs,
		newVideoServiceFeedResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"PublishAction": kitex.NewMethodInfo(
		publishActionHandler,
		newVideoServicePublishActionArgs,
		newVideoServicePublishActionResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"PublishList": kitex.NewMethodInfo(
		publishListHandler,
		newVideoServicePublishListArgs,
		newVideoServicePublishListResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
}

var (
	videoServiceServiceInfo                = NewServiceInfo()
	videoServiceServiceInfoForClient       = NewServiceInfoForClient()
	videoServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return videoServiceServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return videoServiceServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return videoServiceServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(false, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "VideoService"
	handlerType := (*video.VideoService)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "video",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.11.3",
		Extra:           extra,
	}
	return svcInfo
}

func feedHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*video.VideoServiceFeedArgs)
	realResult := result.(*video.VideoServiceFeedResult)
	success, err := handler.(video.VideoService).Feed(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newVideoServiceFeedArgs() interface{} {
	return video.NewVideoServiceFeedArgs()
}

func newVideoServiceFeedResult() interface{} {
	return video.NewVideoServiceFeedResult()
}

func publishActionHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*video.VideoServicePublishActionArgs)
	realResult := result.(*video.VideoServicePublishActionResult)
	success, err := handler.(video.VideoService).PublishAction(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newVideoServicePublishActionArgs() interface{} {
	return video.NewVideoServicePublishActionArgs()
}

func newVideoServicePublishActionResult() interface{} {
	return video.NewVideoServicePublishActionResult()
}

func publishListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*video.VideoServicePublishListArgs)
	realResult := result.(*video.VideoServicePublishListResult)
	success, err := handler.(video.VideoService).PublishList(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newVideoServicePublishListArgs() interface{} {
	return video.NewVideoServicePublishListArgs()
}

func newVideoServicePublishListResult() interface{} {
	return video.NewVideoServicePublishListResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) Feed(ctx context.Context, req *video.FeedRequest) (r *video.FeedResponse, err error) {
	var _args video.VideoServiceFeedArgs
	_args.Req = req
	var _result video.VideoServiceFeedResult
	if err = p.c.Call(ctx, "Feed", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) PublishAction(ctx context.Context, req *video.PublishActionRequest) (r *video.PublishActionResponse, err error) {
	var _args video.VideoServicePublishActionArgs
	_args.Req = req
	var _result video.VideoServicePublishActionResult
	if err = p.c.Call(ctx, "PublishAction", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) PublishList(ctx context.Context, req *video.PublishListRequest) (r *video.PublishListResponse, err error) {
	var _args video.VideoServicePublishListArgs
	_args.Req = req
	var _result video.VideoServicePublishListResult
	if err = p.c.Call(ctx, "PublishList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
