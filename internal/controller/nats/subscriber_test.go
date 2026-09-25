package nats

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	natsgo "github.com/nats-io/nats.go"
)

func TestSubscribeAndPublishIntent(t *testing.T) {
	natsURL := natsgo.DefaultURL

	// 1. Establish SUBSCRIBER
	sub, err := NewSubscriber(natsURL)
	if err != nil {
		t.Fatalf("Cannot establish Subscriber: %v", err)
	}
	defer sub.Close()

	subjectTest := "orders.intent.test"
	queueGroupTest := "workers.group.test"
	done := make(chan bool, 1)

	// 2. Subscribe and waiting
	err = sub.SubscribeQueue(subjectTest, queueGroupTest, func(data []byte) {
		var payload IntentEnvelope
		if err := json.Unmarshal(data, &payload); err != nil {
			log.Printf("[TEST LOGGER] invalid payload: %v", err)
			done <- false
			return
		}

		log.Printf("[TEST LOGGER] received intent successfully: %+v\n", payload)
		done <- true
	})
	if err != nil {
		t.Fatalf("SubscribeQueue failed: %v", err)
	}

	// simulate delay
	time.Sleep(3 * time.Second)

	// 3. Simulating publish data to nats
	publisherConn, err := natsgo.Connect(natsURL)
	if err != nil {
		t.Fatalf("Cannot connect to publisher: %v", err)
	}
	defer publisherConn.Close()

	// data paylod will be published to nats
	mockPayload := []byte(`
		{
			"device_name": "test_device",
			"type":"bgp",
			"data": {
				"name": "eth1/1",
				"ip": "192.168.1.1",
				"description": "uplink",
				"duplex": "full",
				"sub-index": 0,
				"mtu": 1500,
				"speed": 1000,
				"mask": 24,
				"enable": true
			}
		}
	`)

	err = publisherConn.Publish(subjectTest, mockPayload)
	if err != nil {
		t.Fatalf("Failed to publish to NATS: %v", err)
	}
	publisherConn.Flush()

	// 4. Wait and confirm the result
	select {
	case isSuccess := <-done:
		if !isSuccess {
			t.Error("Failed to Unmarshal")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout")
	}
}
