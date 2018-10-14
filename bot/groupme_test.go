package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GroupMeParse", func() {

	It("should parse a GroupMe callback from text", func() {
		cb, err := GroupMeParse(
			`{"name":"test","sender_type":"bot","text":"dumb joke"}`)
		Expect(err).To(BeNil())
		Expect(cb.Name).To(Equal("test"))
		Expect(cb.SenderType).To(Equal("bot"))
		Expect(cb.Text).To(Equal("dumb joke"))
	})

	It("should bubble up any errors", func() {
		_, err := GroupMeParse("blej{}[]")
		Expect(err).ToNot(BeNil())
	})
})

var _ = Describe("GroupMeMessenger", func() {

	It("should send a bot message to GroupMe", func() {
		called := false
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqMsg, err := ioutil.ReadAll(r.Body)
			Expect(err).To(BeNil())
			r.Body.Close()

			var gmMsg *GroupMeMessage
			err = json.Unmarshal(reqMsg, &gmMsg)
			Expect(err).To(BeNil())
			Expect(gmMsg.BotId).To(Equal("test id"))
			Expect(gmMsg.Text).To(Equal("Posted"))

			called = true
			w.WriteHeader(200)
		})
		s := httptest.NewServer(h)
		u, _ := url.Parse(s.URL)

		gm := NewGroupMeMessenger("test id", u)
		code, err := gm.Post("Posted")
		Expect(err).To(BeNil())
		Expect(code).To(Equal(200))
		Expect(called).To(BeTrue())
	})

	It("should proxy the status code from GroupMe", func() {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(403)
		})
		s := httptest.NewServer(h)
		u, _ := url.Parse(s.URL)

		gm := NewGroupMeMessenger("test id", u)
		code, err := gm.Post("Posted")
		Expect(err).To(BeNil())
		Expect(code).To(Equal(403))
	})
})
