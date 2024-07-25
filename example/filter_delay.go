package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/envoyproxyx/go-sdk/envoy"
)

// delayHttpFilter implements envoy.HttpFilter.
//
// This is to demonstrate how to delay the request and response by using Gorooutines.
type delayHttpFilter struct{ requestCounts atomic.Int32 }

func newDelayHttpFilter(string) envoy.HttpFilter { return &delayHttpFilter{} }

// NewInstance implements envoy.HttpFilter.
func (m *delayHttpFilter) NewInstance(e envoy.EnvoyFilterInstance) envoy.HttpFilterInstance {
	// NewInstance is called for each new Http request, so we can use a counter to track the number of requests.
	// On the other hand, that means this function must be thread-safe.
	id := m.requestCounts.Add(1)
	return &delayHttpFilterInstance{id: id, envoyFilter: e}
}

// Destroy implements envoy.HttpFilter.
func (m *delayHttpFilter) Destroy() {
	fmt.Println("Destroy called")
}

// delayHttpFilterInstance implements envoy.HttpFilterInstance.
type delayHttpFilterInstance struct {
	id          int32
	envoyFilter envoy.EnvoyFilterInstance
}

// RequestHeaders implements envoy.HttpFilterInstance.
func (h *delayHttpFilterInstance) RequestHeaders(_ envoy.RequestHeaders, _ bool) envoy.RequestHeadersStatus {
	if h.id == 1 {
		go func() {
			fmt.Println("blocking for 1 second at EventHttpRequestHeaders with id", h.id)
			time.Sleep(1 * time.Second)
			fmt.Println("calling ContinueRequest with id", h.id)
			if h.envoyFilter != nil { // By calling ContinueResponse, the filter might have been destroyed.
				h.envoyFilter.ContinueRequest()
			}
		}()
		fmt.Println("EventHttpRequestHeaders returning StopAllIterationAndBuffer with id", h.id)
		return envoy.RequestHeadersStatusStopAllIterationAndBuffer
	}
	fmt.Println("EventHttpRequestHeaders called with id", h.id)
	return envoy.HeadersStatusContinue
}

// RequestBody implements envoy.HttpFilterInstance.
func (h *delayHttpFilterInstance) RequestBody(_ envoy.RequestBodyBuffer, _ bool) envoy.RequestBodyStatus {
	if h.id == 2 {
		go func() {
			fmt.Println("blocking for 1 second at EventHttpRequestBody with id", h.id)
			time.Sleep(1 * time.Second)
			fmt.Println("calling ContinueRequest with id", h.id)
			if h.envoyFilter != nil { // By calling ContinueResponse, the filter might have been destroyed.
				h.envoyFilter.ContinueRequest()
			}
		}()
		fmt.Println("EventHttpRequestBody returning StopIterationAndBuffer with id", h.id)
		return envoy.RequestBodyStatusStopIterationAndBuffer
	}
	fmt.Println("EventHttpRequestBody called with id", h.id)
	return envoy.RequestBodyStatusContinue
}

// ResponseHeaders implements envoy.HttpFilterInstance.
func (h *delayHttpFilterInstance) ResponseHeaders(_ envoy.ResponseHeaders, _ bool) envoy.ResponseHeadersStatus {
	if h.id == 3 {
		go func() {
			fmt.Println("blocking for 1 second at EventHttpResponseHeaders with id", h.id)
			time.Sleep(1 * time.Second)
			fmt.Println("calling ContinueResponse with id", h.id)
			if h.envoyFilter != nil { // By calling ContinueResponse, the filter might have been destroyed.
				h.envoyFilter.ContinueResponse()
			}
		}()
		fmt.Println("EventHttpResponseHeaders returning StopAllIterationAndBuffer with id", h.id)
		return envoy.ResponseHeadersStatusStopAllIterationAndBuffer
	}
	fmt.Println("EventHttpResponseHeaders called with id", h.id)
	return envoy.ResponseHeadersStatusContinue
}

// ResponseBody implements envoy.HttpFilterInstance.
func (h *delayHttpFilterInstance) ResponseBody(_ envoy.ResponseBodyBuffer, _ bool) envoy.ResponseBodyStatus {
	if h.id == 4 {
		go func() {
			fmt.Println("blocking for 1 second at EventHttpResponseBody with id", h.id)
			time.Sleep(1 * time.Second)
			fmt.Println("calling ContinueResponse with id", h.id)
			if h.envoyFilter != nil { // By calling ContinueResponse, the filter might have been destroyed.
				h.envoyFilter.ContinueResponse()
			}
		}()
		fmt.Println("EventHttpResponseBody returning StopIterationAndBuffer with id", h.id)
		return envoy.ResponseBodyStatusStopIterationAndBuffer
	}
	fmt.Println("EventHttpResponseBody called with id", h.id)
	return envoy.ResponseBodyStatusContinue
}

// Destroy implements envoy.HttpFilterInstance.
func (h *delayHttpFilterInstance) Destroy() {
	// After the request is done, we can clean up the filter instance.
	h.envoyFilter = nil
}
