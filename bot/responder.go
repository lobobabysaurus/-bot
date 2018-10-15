package main

import (
	"math/rand"
	"strings"
)

type Responder interface {
	Respond(string) *string
}

type TriggerResponder struct {
	trigger   string
	prefix    string
	responses []string
}

func NewTriggerResponder(
	trigger string, prefix string, responses []string) Responder {

	return &TriggerResponder{
		trigger:   trigger,
		prefix:    prefix,
		responses: responses,
	}
}

func (r *TriggerResponder) Respond(in string) *string {
	if strings.Contains(in, r.trigger) {
		resp := r.prefix + randomElement(r.responses)
		return &resp
	} else {
		return nil
	}
}

type randomFunc func() float64

type RandomResponder struct {
	threshold float64
	responses []string
	random    randomFunc
}

func NewRandomResponder(
	threshold float64, random randomFunc, responses []string) Responder {

	return &RandomResponder{
		threshold: threshold,
		random:    random,
		responses: responses,
	}
}

func (r *RandomResponder) Respond(in string) *string {
	if r.random() < r.threshold {
		resp := randomElement(r.responses)
		return &resp
	} else {
		return nil
	}
}

func randomElement(s []string) string {
	i := rand.Int() % len(s)
	return s[i]
}
