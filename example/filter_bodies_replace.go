package main

import (
	"fmt"

	"github.com/envoyproxyx/go-sdk/envoy"
)

// bodiesReplaceHttpFilter implements envoy.HttpFilter.
//
// This is to demonstrate how to use body manipulation APIs.
type bodiesReplaceHttpFilter struct{}

func newbodiesReplaceHttpFilter(string) envoy.HttpFilter { return &bodiesReplaceHttpFilter{} }

// NewHttpFilterInstance implements envoy.HttpFilter.
func (f *bodiesReplaceHttpFilter) NewHttpFilterInstance(envoyFilter envoy.EnvoyFilterInstance) envoy.HttpFilterInstance {
	return &bodiesReplaceHttpFilterInstance{envoyFilter: envoyFilter}
}

// Destroy implements envoy.HttpFilter.
func (f *bodiesReplaceHttpFilter) Destroy() {}

// bodiesReplaceHttpFilterInstance implements envoy.HttpFilterInstance.
type bodiesReplaceHttpFilterInstance struct {
	envoyFilter                                      envoy.EnvoyFilterInstance
	requestAppend, requestPrepend, requestReplace    string
	responseAppend, responsePrepend, responseReplace string
}

// EventHttpRequestHeaders implements envoy.HttpFilterInstance.
func (h *bodiesReplaceHttpFilterInstance) EventHttpRequestHeaders(headers envoy.RequestHeaders, _ bool) envoy.EventHttpRequestHeadersStatus {
	append, ok := headers.Get("append")
	if ok {
		h.requestAppend = append.String()
	}
	prepend, ok := headers.Get("prepend")
	if ok {
		h.requestPrepend = prepend.String()
	}
	replace, ok := headers.Get("replace")
	if ok {
		h.requestReplace = replace.String()
	}
	headers.Remove("content-length") // Remove the content-length header to reset the length.
	return envoy.EventHttpRequestHeadersStatusContinue
}

// EventHttpRequestBody implements envoy.HttpFilterInstance.
func (h *bodiesReplaceHttpFilterInstance) EventHttpRequestBody(body envoy.RequestBodyBuffer, endOfStream bool) envoy.EventHttpRequestBodyStatus {
	if !endOfStream {
		// Wait for the end of the stream to see the full body.
		return envoy.EventHttpRequestBodyStatusStopIterationAndBuffer
	}

	entireBody := h.envoyFilter.GetRequestBodyBuffer()
	if h.requestAppend != "" {
		entireBody.Append([]byte(h.requestAppend))
	}
	if h.requestPrepend != "" {
		entireBody.Prepend([]byte(h.requestPrepend))
	}
	if h.requestReplace != "" {
		entireBody.Replace([]byte(h.requestReplace))
	}
	return envoy.EventHttpRequestBodyStatusContinue
}

// EventHttpResponseHeaders implements envoy.HttpFilterInstance.
func (h *bodiesReplaceHttpFilterInstance) EventHttpResponseHeaders(headers envoy.ResponseHeaders, _ bool) envoy.EventHttpResponseHeadersStatus {
	append, ok := headers.Get("append")
	if ok {
		h.responseAppend = append.String()
	}
	prepend, ok := headers.Get("prepend")
	if ok {
		h.responsePrepend = prepend.String()
	}
	replace, ok := headers.Get("replace")
	if ok {
		h.responseReplace = replace.String()
	}
	headers.Remove("content-length") // Remove the content-length header to reset the length.
	return envoy.EventHttpResponseHeadersStatusContinue
}

// EventHttpResponseBody implements envoy.HttpFilterInstance.
func (h *bodiesReplaceHttpFilterInstance) EventHttpResponseBody(body envoy.ResponseBodyBuffer, endOfStream bool) envoy.EventHttpResponseBodyStatus {
	fmt.Printf("new request body frame: %s\n", string(body.Copy()))
	if !endOfStream {
		// Wait for the end of the stream to see the full body.
		return envoy.EventHttpResponseBodyStatusStopIterationAndBuffer
	}

	entireBody := h.envoyFilter.GetResponseBodyBuffer()
	if h.responseAppend != "" {
		entireBody.Append([]byte(h.responseAppend))
	}
	if h.responsePrepend != "" {
		entireBody.Prepend([]byte(h.responsePrepend))
	}
	if h.responseReplace != "" {
		entireBody.Replace([]byte(h.responseReplace))
	}
	return envoy.EventHttpResponseBodyStatusContinue
}

// EventHttpDestroy implements envoy.HttpFilterInstance.
func (h *bodiesReplaceHttpFilterInstance) EventHttpDestroy() {}
