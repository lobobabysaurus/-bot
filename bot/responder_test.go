package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TriggerResponder", func() {

	responses := []string{"these", "are", "some", "responses"}

	Context("when trigger is in passed message", func() {

		It("should give a random response", func() {
			responder := NewTriggerResponder("trig", "", responses)
			resp := responder.Respond("trig")
			Expect(resp).ToNot(BeNil())
			Expect(responses).To(ContainElement(*resp))
		})

		It("should prefix the reponse", func() {
			prefix := "prefix:"
			responder := NewTriggerResponder("trig", prefix, responses)
			resp := responder.Respond("trigger")
			Expect(resp).ToNot(BeNil())
			Expect(*resp).To(HavePrefix(prefix))
			prefixlessResp := (*resp)[len(prefix):]
			Expect(responses).To(ContainElement(prefixlessResp))
		})

		It("should be case senstive", func() {
			responder := NewTriggerResponder("trig", "", responses)
			resp := responder.Respond("TrIg")
			Expect(resp).To(BeNil())
		})
	})

	Context("when trigger is not in passed message", func() {
		It("should give no response", func() {
			responder := NewTriggerResponder("trig", "", responses)
			resp := responder.Respond("not a match")
			Expect(resp).To(BeNil())
		})
	})
})

var _ = Describe("RandomResponder", func() {

	responses := []string{"these", "are", "some", "responses"}

	Context("when random value is under the threshold", func() {
		It("should give a random response", func() {
			responder := NewRandomResponder(
				0.5, func() float64 { return 0.49 }, responses)
			resp := responder.Respond("throwaway")
			Expect(resp).ToNot(BeNil())
			Expect(responses).To(ContainElement(*resp))
		})
	})

	Context("when random value is equal or over the threshold", func() {
		It("should give no response", func() {
			responder := NewRandomResponder(
				0.5, func() float64 { return 0.5 }, responses)
			resp := responder.Respond("throwaway")
			Expect(resp).To(BeNil())
		})
	})
})
