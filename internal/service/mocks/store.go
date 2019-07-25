// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"
	"time"
)

type FakeStore struct {
	DelStub        func(string) error
	delMutex       sync.RWMutex
	delArgsForCall []struct {
		arg1 string
	}
	delReturns struct {
		result1 error
	}
	delReturnsOnCall map[int]struct {
		result1 error
	}
	GetStub        func(string) (string, error)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		arg1 string
	}
	getReturns struct {
		result1 string
		result2 error
	}
	getReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	RandomKeyStub        func() (string, error)
	randomKeyMutex       sync.RWMutex
	randomKeyArgsForCall []struct {
	}
	randomKeyReturns struct {
		result1 string
		result2 error
	}
	randomKeyReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	SetStub        func(string, string, time.Duration) error
	setMutex       sync.RWMutex
	setArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 time.Duration
	}
	setReturns struct {
		result1 error
	}
	setReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeStore) Del(arg1 string) error {
	fake.delMutex.Lock()
	ret, specificReturn := fake.delReturnsOnCall[len(fake.delArgsForCall)]
	fake.delArgsForCall = append(fake.delArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Del", []interface{}{arg1})
	fake.delMutex.Unlock()
	if fake.DelStub != nil {
		return fake.DelStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.delReturns
	return fakeReturns.result1
}

func (fake *FakeStore) DelCallCount() int {
	fake.delMutex.RLock()
	defer fake.delMutex.RUnlock()
	return len(fake.delArgsForCall)
}

func (fake *FakeStore) DelCalls(stub func(string) error) {
	fake.delMutex.Lock()
	defer fake.delMutex.Unlock()
	fake.DelStub = stub
}

func (fake *FakeStore) DelArgsForCall(i int) string {
	fake.delMutex.RLock()
	defer fake.delMutex.RUnlock()
	argsForCall := fake.delArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeStore) DelReturns(result1 error) {
	fake.delMutex.Lock()
	defer fake.delMutex.Unlock()
	fake.DelStub = nil
	fake.delReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeStore) DelReturnsOnCall(i int, result1 error) {
	fake.delMutex.Lock()
	defer fake.delMutex.Unlock()
	fake.DelStub = nil
	if fake.delReturnsOnCall == nil {
		fake.delReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.delReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeStore) Get(arg1 string) (string, error) {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Get", []interface{}{arg1})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeStore) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *FakeStore) GetCalls(stub func(string) (string, error)) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = stub
}

func (fake *FakeStore) GetArgsForCall(i int) string {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	argsForCall := fake.getArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeStore) GetReturns(result1 string, result2 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeStore) GetReturnsOnCall(i int, result1 string, result2 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeStore) RandomKey() (string, error) {
	fake.randomKeyMutex.Lock()
	ret, specificReturn := fake.randomKeyReturnsOnCall[len(fake.randomKeyArgsForCall)]
	fake.randomKeyArgsForCall = append(fake.randomKeyArgsForCall, struct {
	}{})
	fake.recordInvocation("RandomKey", []interface{}{})
	fake.randomKeyMutex.Unlock()
	if fake.RandomKeyStub != nil {
		return fake.RandomKeyStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.randomKeyReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeStore) RandomKeyCallCount() int {
	fake.randomKeyMutex.RLock()
	defer fake.randomKeyMutex.RUnlock()
	return len(fake.randomKeyArgsForCall)
}

func (fake *FakeStore) RandomKeyCalls(stub func() (string, error)) {
	fake.randomKeyMutex.Lock()
	defer fake.randomKeyMutex.Unlock()
	fake.RandomKeyStub = stub
}

func (fake *FakeStore) RandomKeyReturns(result1 string, result2 error) {
	fake.randomKeyMutex.Lock()
	defer fake.randomKeyMutex.Unlock()
	fake.RandomKeyStub = nil
	fake.randomKeyReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeStore) RandomKeyReturnsOnCall(i int, result1 string, result2 error) {
	fake.randomKeyMutex.Lock()
	defer fake.randomKeyMutex.Unlock()
	fake.RandomKeyStub = nil
	if fake.randomKeyReturnsOnCall == nil {
		fake.randomKeyReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.randomKeyReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeStore) Set(arg1 string, arg2 string, arg3 time.Duration) error {
	fake.setMutex.Lock()
	ret, specificReturn := fake.setReturnsOnCall[len(fake.setArgsForCall)]
	fake.setArgsForCall = append(fake.setArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 time.Duration
	}{arg1, arg2, arg3})
	fake.recordInvocation("Set", []interface{}{arg1, arg2, arg3})
	fake.setMutex.Unlock()
	if fake.SetStub != nil {
		return fake.SetStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.setReturns
	return fakeReturns.result1
}

func (fake *FakeStore) SetCallCount() int {
	fake.setMutex.RLock()
	defer fake.setMutex.RUnlock()
	return len(fake.setArgsForCall)
}

func (fake *FakeStore) SetCalls(stub func(string, string, time.Duration) error) {
	fake.setMutex.Lock()
	defer fake.setMutex.Unlock()
	fake.SetStub = stub
}

func (fake *FakeStore) SetArgsForCall(i int) (string, string, time.Duration) {
	fake.setMutex.RLock()
	defer fake.setMutex.RUnlock()
	argsForCall := fake.setArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeStore) SetReturns(result1 error) {
	fake.setMutex.Lock()
	defer fake.setMutex.Unlock()
	fake.SetStub = nil
	fake.setReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeStore) SetReturnsOnCall(i int, result1 error) {
	fake.setMutex.Lock()
	defer fake.setMutex.Unlock()
	fake.SetStub = nil
	if fake.setReturnsOnCall == nil {
		fake.setReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeStore) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.delMutex.RLock()
	defer fake.delMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	fake.randomKeyMutex.RLock()
	defer fake.randomKeyMutex.RUnlock()
	fake.setMutex.RLock()
	defer fake.setMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeStore) recordInvocation(key string, args []interface{}) {
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
