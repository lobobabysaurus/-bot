package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suiite")
}

var _ = Describe("Handler", func() {

	var callbackRequest events.APIGatewayProxyRequest

	Context("When a user sends a message to GroupMe", func() {

		BeforeEach(func() {
			callbackRequest = wrapCallback(&GroupMeCallback{
				Name:       "testing",
				SenderType: "user",
			})
		})

		Context("And a random value under the Ilya Rate happens", func() {

			BeforeEach(func() {
				RandFunc = func() float64 {
					return 0.019
				}
			})

			It("should send a bot message to GroupMe", func() {
				BotId = "test id"
				called := false
				h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					reqMessage, err := ioutil.ReadAll(r.Body)
					Expect(err).To(BeNil())
					r.Body.Close()

					var gmMsg *GroupMeMessage
					err = json.Unmarshal(reqMessage, &gmMsg)
					Expect(err).To(BeNil())
					Expect(gmMsg.BotId).To(Equal("test id"))
					Expect(gmMsg.Text).To(Equal(
						"When I talk English, I’m still computing."))

					called = true
					w.WriteHeader(200)
				})
				s := rigServer(h)
				defer s.Close()

				resp, err := samHandler(callbackRequest)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
				Expect(called).To(BeTrue())
			})

			It("should proxy the response from GroupMe", func() {
				h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(403)
				})
				s := rigServer(h)
				defer s.Close()

				resp, err := samHandler(callbackRequest)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(403))
			})
		})

		Context("And a random value over the Ilya Rate happens", func() {

			BeforeEach(func() {
				RandFunc = func() float64 {
					return 0.021
				}
			})

			It("should not send any messages", func() {
				h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					Fail("This server should not be hit")
				})
				s := rigServer(h)
				defer s.Close()

				resp, err := samHandler(callbackRequest)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
			})
		})
	})

	Context("When a bot sends a message to GroupMe", func() {

		BeforeEach(func() {
			callbackRequest = wrapCallback(&GroupMeCallback{
				Name:       "testing",
				SenderType: "bot",
			})
		})

		It("should not send any messages", func() {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Fail("This server should not be hit")
			})
			s := rigServer(h)
			defer s.Close()

			resp, err := samHandler(callbackRequest)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
		})
	})
})

func rigServer(h http.Handler) *httptest.Server {
	s := httptest.NewServer(h)
	ChatURL = s.URL
	return s
}

func wrapCallback(body *GroupMeCallback) events.APIGatewayProxyRequest {
	bodyBytes, err := json.Marshal(body)
	Expect(err).To(BeNil())
	return events.APIGatewayProxyRequest{
		Body: string(bodyBytes),
	}
}
