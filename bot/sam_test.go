package main

import (
	"errors"

	"github.com/aws/aws-lambda-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	WhoopsErr = errors.New("whoops")
)

type MockSuccesfullGroupMe struct {
	calls   int
	lastMsg string
}

func (g *MockSuccesfullGroupMe) Post(text string) (int, error) {
	g.calls++
	g.lastMsg = text
	return 123, nil
}

type MockFailedGroupMe struct {
	calls int
}

func (g *MockFailedGroupMe) Post(text string) (int, error) {
	g.calls++
	return 0, WhoopsErr
}

var _ = Describe("SAMHandler", func() {

	Context("on callback parsing errors", func() {
		It("Should bubble up the error", func() {
			h := NewSAMHandler(&MockSuccesfullGroupMe{})
			req := events.APIGatewayProxyRequest{Body: "{not valid"}
			resp, err := h.Handle(req)
			Expect(resp).To(Equal(EmptyResponse))
			Expect(err).ToNot(BeNil())
		})
	})

	Context("when sender is not a user", func() {
		It("Should exit early", func() {
			r := NewTriggerResponder("trigger", "", []string{"response"})
			g := &MockSuccesfullGroupMe{}
			h := NewSAMHandler(g, r)
			req := events.APIGatewayProxyRequest{
				Body: `{"name":"test","sender_type":"bot","text":"trigger"}`,
			}
			resp, err := h.Handle(req)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			Expect(g.calls).To(Equal(0))
		})
	})

	Context("when responders give a response", func() {

		It("should post the first successful one", func() {
			r1 := NewTriggerResponder("trogger", "", []string{"response1"})
			r2 := NewTriggerResponder("trigger", "", []string{"response2"})
			r3 := NewTriggerResponder("trigger", "", []string{"response3"})
			g := &MockSuccesfullGroupMe{}
			h := NewSAMHandler(g, r1, r2, r3)
			req := events.APIGatewayProxyRequest{
				Body: `{"name":"test","sender_type":"user","text":"trigger"}`,
			}
			_, err := h.Handle(req)
			Expect(err).To(BeNil())
			Expect(g.calls).To(Equal(1))
			Expect(g.lastMsg).To(Equal("response2"))
		})

		It("Should bubble the status code from GroupMe", func() {
			r := NewTriggerResponder("trigger", "", []string{"response"})
			h := NewSAMHandler(&MockSuccesfullGroupMe{}, r)
			req := events.APIGatewayProxyRequest{
				Body: `{"name":"test","sender_type":"user","text":"trigger"}`,
			}
			resp, err := h.Handle(req)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(123))
		})

		It("Should bubble up any errors", func() {
			r := NewTriggerResponder("trigger", "", []string{"response"})
			h := NewSAMHandler(&MockFailedGroupMe{}, r)
			req := events.APIGatewayProxyRequest{
				Body: `{"name":"test","sender_type":"user","text":"trigger"}`,
			}
			_, err := h.Handle(req)
			Expect(err).To(Equal(WhoopsErr))
		})
	})

	It("Should exit cleanly if no responders respond", func() {
		h := NewSAMHandler(&MockFailedGroupMe{})
		req := events.APIGatewayProxyRequest{
			Body: `{"name":"test","sender_type":"user","text":"trigger"}`,
		}
		resp, err := h.Handle(req)
		Expect(err).To(BeNil())
		Expect(resp.StatusCode).To(Equal(200))
	})
})
