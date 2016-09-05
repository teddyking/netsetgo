// This file was generated by counterfeiter
package netsetgofakes

import (
	"sync"

	"github.com/teddyking/netsetgo"
)

type FakeHostConfigurer struct {
	ApplyStub        func(netConfig netsetgo.NetworkConfig, pid int) error
	applyMutex       sync.RWMutex
	applyArgsForCall []struct {
		netConfig netsetgo.NetworkConfig
		pid       int
	}
	applyReturns struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeHostConfigurer) Apply(netConfig netsetgo.NetworkConfig, pid int) error {
	fake.applyMutex.Lock()
	fake.applyArgsForCall = append(fake.applyArgsForCall, struct {
		netConfig netsetgo.NetworkConfig
		pid       int
	}{netConfig, pid})
	fake.recordInvocation("Apply", []interface{}{netConfig, pid})
	fake.applyMutex.Unlock()
	if fake.ApplyStub != nil {
		return fake.ApplyStub(netConfig, pid)
	} else {
		return fake.applyReturns.result1
	}
}

func (fake *FakeHostConfigurer) ApplyCallCount() int {
	fake.applyMutex.RLock()
	defer fake.applyMutex.RUnlock()
	return len(fake.applyArgsForCall)
}

func (fake *FakeHostConfigurer) ApplyArgsForCall(i int) (netsetgo.NetworkConfig, int) {
	fake.applyMutex.RLock()
	defer fake.applyMutex.RUnlock()
	return fake.applyArgsForCall[i].netConfig, fake.applyArgsForCall[i].pid
}

func (fake *FakeHostConfigurer) ApplyReturns(result1 error) {
	fake.ApplyStub = nil
	fake.applyReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeHostConfigurer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.applyMutex.RLock()
	defer fake.applyMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeHostConfigurer) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ netsetgo.HostConfigurer = new(FakeHostConfigurer)
