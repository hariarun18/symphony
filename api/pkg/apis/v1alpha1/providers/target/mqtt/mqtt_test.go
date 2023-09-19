/*

	MIT License

	Copyright (c) Microsoft Corporation.

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE

*/

package mqtt

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/azure/symphony/api/pkg/apis/v1alpha1/model"
	"github.com/azure/symphony/api/pkg/apis/v1alpha1/providers/target/conformance"
	"github.com/azure/symphony/coa/pkg/apis/v1alpha2"
	gmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
)

func TestDoubleIni(t *testing.T) {
	testMQTT := os.Getenv("TEST_MQTT")
	if testMQTT == "" {
		t.Skip("Skipping because TES_MQTT enviornment variable is not set")
	}
	config := MQTTTargetProviderConfig{
		Name:          "me",
		BrokerAddress: "tcp://20.118.146.198:1883",
		ClientID:      "coa-test2",
		RequestTopic:  "coa-request",
		ResponseTopic: "coa-response",
	}
	provider := MQTTTargetProvider{}
	err := provider.Init(config)
	assert.Nil(t, err)
	err = provider.Init(config)
	assert.Nil(t, err)
}

func TestGet(t *testing.T) {
	testMQTT := os.Getenv("TEST_MQTT")
	if testMQTT == "" {
		t.Skip("Skipping because TES_MQTT enviornment variable is not set")
	}
	config := MQTTTargetProviderConfig{
		Name:          "me",
		BrokerAddress: "tcp://20.118.146.198:1883",
		ClientID:      "coa-test2",
		RequestTopic:  "coa-request",
		ResponseTopic: "coa-response",
	}
	provider := MQTTTargetProvider{}
	err := provider.Init(config)
	assert.Nil(t, err)

	opts := gmqtt.NewClientOptions().AddBroker(config.BrokerAddress).SetClientID("test-sender")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := gmqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := c.Subscribe(config.RequestTopic, 0, func(client gmqtt.Client, msg gmqtt.Message) {
		var response v1alpha2.COAResponse
		ret := make([]model.ComponentSpec, 0)
		data, _ := json.Marshal(ret)
		response.State = v1alpha2.OK
		response.Metadata = make(map[string]string)
		response.Metadata["call-context"] = "TargetProvider-Get"
		response.Body = data
		data, _ = json.Marshal(response)
		token := c.Publish(config.ResponseTopic, 0, false, data) //sending COARequest directly doesn't seem to work
		token.Wait()

	}); token.Wait() && token.Error() != nil {
		if token.Error().Error() != "subscription exists" {
			panic(token.Error())
		}
	}

	arr, err := provider.Get(context.Background(), model.DeploymentSpec{}, nil)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(arr))
}
func TestGetBad(t *testing.T) {
	testMQTT := os.Getenv("TEST_MQTT")
	if testMQTT == "" {
		t.Skip("Skipping because TES_MQTT enviornment variable is not set")
	}
	config := MQTTTargetProviderConfig{
		Name:          "me",
		BrokerAddress: "tcp://20.118.146.198:1883",
		ClientID:      "coa-test2",
		RequestTopic:  "coa-request",
		ResponseTopic: "coa-response",
	}
	provider := MQTTTargetProvider{}
	err := provider.Init(config)
	assert.Nil(t, err)

	opts := gmqtt.NewClientOptions().AddBroker(config.BrokerAddress).SetClientID("test-sender")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := gmqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := c.Subscribe(config.RequestTopic, 0, func(client gmqtt.Client, msg gmqtt.Message) {
		var response v1alpha2.COAResponse
		response.State = v1alpha2.InternalError
		response.Metadata = make(map[string]string)
		response.Metadata["call-context"] = "TargetProvider-Get"
		response.Body = []byte("BAD!!")
		data, _ := json.Marshal(response)
		token := c.Publish(config.ResponseTopic, 0, false, data) //sending COARequest directly doesn't seem to work
		token.Wait()

	}); token.Wait() && token.Error() != nil {
		if token.Error().Error() != "subscription exists" {
			panic(token.Error())
		}
	}

	_, err = provider.Get(context.Background(), model.DeploymentSpec{}, nil)

	assert.NotNil(t, err)
	assert.Equal(t, "BAD!!", err.Error())
}
func TestApply(t *testing.T) {
	testMQTT := os.Getenv("TEST_MQTT")
	if testMQTT == "" {
		t.Skip("Skipping because TES_MQTT enviornment variable is not set")
	}
	config := MQTTTargetProviderConfig{
		Name:          "me",
		BrokerAddress: "tcp://20.118.146.198:1883",
		ClientID:      "coa-test2",
		RequestTopic:  "coa-request",
		ResponseTopic: "coa-response",
	}
	provider := MQTTTargetProvider{}
	err := provider.Init(config)
	assert.Nil(t, err)

	opts := gmqtt.NewClientOptions().AddBroker(config.BrokerAddress).SetClientID("test-sender")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := gmqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := c.Subscribe(config.RequestTopic, 0, func(client gmqtt.Client, msg gmqtt.Message) {
		var response v1alpha2.COAResponse
		response.State = v1alpha2.OK
		response.Metadata = make(map[string]string)
		response.Metadata["call-context"] = "TargetProvider-Apply"
		data, _ := json.Marshal(response)
		token := c.Publish(config.ResponseTopic, 0, false, data) //sending COARequest directly doesn't seem to work
		token.Wait()

	}); token.Wait() && token.Error() != nil {
		if token.Error().Error() != "subscription exists" {
			panic(token.Error())
		}
	}

	_, err = provider.Apply(context.Background(), model.DeploymentSpec{}, model.DeploymentStep{}, false) //TODO: this is probably broken: the step should contain at least a component

	assert.Nil(t, err)
}
func TestApplyBad(t *testing.T) {
	testMQTT := os.Getenv("TEST_MQTT")
	if testMQTT == "" {
		t.Skip("Skipping because TES_MQTT enviornment variable is not set")
	}
	config := MQTTTargetProviderConfig{
		Name:          "me",
		BrokerAddress: "tcp://20.118.146.198:1883",
		ClientID:      "coa-test2",
		RequestTopic:  "coa-request",
		ResponseTopic: "coa-response",
	}
	provider := MQTTTargetProvider{}
	err := provider.Init(config)
	assert.Nil(t, err)

	opts := gmqtt.NewClientOptions().AddBroker(config.BrokerAddress).SetClientID("test-sender")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := gmqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := c.Subscribe(config.RequestTopic, 0, func(client gmqtt.Client, msg gmqtt.Message) {
		var response v1alpha2.COAResponse
		response.State = v1alpha2.InternalError
		response.Metadata = make(map[string]string)
		response.Metadata["call-context"] = "TargetProvider-Apply"
		data, _ := json.Marshal(response)
		token := c.Publish(config.ResponseTopic, 0, false, data) //sending COARequest directly doesn't seem to work
		token.Wait()

	}); token.Wait() && token.Error() != nil {
		if token.Error().Error() != "subscription exists" {
			panic(token.Error())
		}
	}

	_, err = provider.Apply(context.Background(), model.DeploymentSpec{}, model.DeploymentStep{}, false) //TODO: this is probably broken - the step should contain at least one component

	assert.NotNil(t, err)
}

func TestARemove(t *testing.T) {
	testMQTT := os.Getenv("TEST_MQTT")
	if testMQTT == "" {
		t.Skip("Skipping because TES_MQTT enviornment variable is not set")
	}
	config := MQTTTargetProviderConfig{
		Name:          "me",
		BrokerAddress: "tcp://20.118.146.198:1883",
		ClientID:      "coa-test2",
		RequestTopic:  "coa-request",
		ResponseTopic: "coa-response",
	}
	provider := MQTTTargetProvider{}
	err := provider.Init(config)
	assert.Nil(t, err)

	opts := gmqtt.NewClientOptions().AddBroker(config.BrokerAddress).SetClientID("test-sender")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := gmqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := c.Subscribe(config.RequestTopic, 0, func(client gmqtt.Client, msg gmqtt.Message) {
		var response v1alpha2.COAResponse
		response.State = v1alpha2.OK
		response.Metadata = make(map[string]string)
		response.Metadata["call-context"] = "TargetProvider-Remove"
		data, _ := json.Marshal(response)
		token := c.Publish(config.ResponseTopic, 0, false, data) //sending COARequest directly doesn't seem to work
		token.Wait()

	}); token.Wait() && token.Error() != nil {
		if token.Error().Error() != "subscription exists" {
			panic(token.Error())
		}
	}

	_, err = provider.Apply(context.Background(), model.DeploymentSpec{}, model.DeploymentStep{}, false)
	assert.Nil(t, err)
}
func TestARemoveBad(t *testing.T) {
	testMQTT := os.Getenv("TEST_MQTT")
	if testMQTT == "" {
		t.Skip("Skipping because TES_MQTT enviornment variable is not set")
	}
	config := MQTTTargetProviderConfig{
		Name:          "me",
		BrokerAddress: "tcp://20.118.146.198:1883",
		ClientID:      "coa-test2",
		RequestTopic:  "coa-request",
		ResponseTopic: "coa-response",
	}
	provider := MQTTTargetProvider{}
	err := provider.Init(config)
	assert.Nil(t, err)

	opts := gmqtt.NewClientOptions().AddBroker(config.BrokerAddress).SetClientID("test-sender")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := gmqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := c.Subscribe(config.RequestTopic, 0, func(client gmqtt.Client, msg gmqtt.Message) {
		var response v1alpha2.COAResponse
		response.State = v1alpha2.InternalError
		response.Metadata = make(map[string]string)
		response.Metadata["call-context"] = "TargetProvider-Remove"
		data, _ := json.Marshal(response)
		token := c.Publish(config.ResponseTopic, 0, false, data) //sending COARequest directly doesn't seem to work
		token.Wait()

	}); token.Wait() && token.Error() != nil {
		if token.Error().Error() != "subscription exists" {
			panic(token.Error())
		}
	}

	_, err = provider.Apply(context.Background(), model.DeploymentSpec{}, model.DeploymentStep{}, false) //TODO: this is probably broken, a step should have at least one component

	assert.NotNil(t, err)
}
func TestGetApply(t *testing.T) {
	testMQTT := os.Getenv("TEST_MQTT")
	if testMQTT == "" {
		t.Skip("Skipping because TES_MQTT enviornment variable is not set")
	}
	config := MQTTTargetProviderConfig{
		Name:          "me",
		BrokerAddress: "tcp://20.118.146.198:1883",
		ClientID:      "coa-test2",
		RequestTopic:  "coa-request",
		ResponseTopic: "coa-response",
	}
	provider := MQTTTargetProvider{}
	err := provider.Init(config)
	assert.Nil(t, err)

	opts := gmqtt.NewClientOptions().AddBroker(config.BrokerAddress).SetClientID("test-sender")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := gmqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := c.Subscribe(config.RequestTopic, 0, func(client gmqtt.Client, msg gmqtt.Message) {
		var request v1alpha2.COARequest
		json.Unmarshal(msg.Payload(), &request)
		var response v1alpha2.COAResponse
		response.Metadata = make(map[string]string)
		if request.Method == "GET" {
			response.Metadata["call-context"] = "TargetProvider-Get"
			ret := make([]model.ComponentSpec, 0)
			data, _ := json.Marshal(ret)
			response.State = v1alpha2.OK
			response.Body = data
		} else {
			response.Metadata["call-context"] = "TargetProvider-Apply"
			response.State = v1alpha2.OK
		}

		data, _ := json.Marshal(response)
		token := c.Publish(config.ResponseTopic, 0, false, data) //sending COARequest directly doesn't seem to work
		token.Wait()

	}); token.Wait() && token.Error() != nil {
		if token.Error().Error() != "subscription exists" {
			panic(token.Error())
		}
	}

	arr, err := provider.Get(context.Background(), model.DeploymentSpec{}, nil)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(arr))

	err = provider.Init(config)
	assert.Nil(t, err)

	_, err = provider.Apply(context.Background(), model.DeploymentSpec{}, model.DeploymentStep{}, false) //TODO: this is probably broken - a step should have at least one component
	assert.Nil(t, err)
}

// Conformance: you should call the conformance suite to ensure provider conformance
func TestConformanceSuite(t *testing.T) {
	provider := &MQTTTargetProvider{}
	_ = provider.Init(MQTTTargetProviderConfig{})
	// assert.Nil(t, err) okay if provider is not fully initialized
	conformance.ConformanceSuite(t, provider)
}
