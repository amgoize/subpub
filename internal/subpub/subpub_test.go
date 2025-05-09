package subpub

import (
	"context"
	"testing"
	"time"
)

func TestSubPub(t *testing.T) {
	sp := NewSubPub()

	var msg interface{}
	handler := func(m interface{}) {
		msg = m
	}

	sub, err := sp.Subscripe("sub1", handler)
	if err != nil {
		t.Fatalf("Ошибка при подписке: %v", err)
	}

	err = sp.Publish("sub1", "message")
	if err != nil {
		t.Fatalf("Ошибка при публикации: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	if msg != "message" {
		t.Errorf("Ожидалось 'message', а получено: %v", msg)
	}

	sub.Unsubscribe()

	err = sp.Publish("sub1", "new message")
	if err != nil {
		t.Fatalf("Ошибка при публикации: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	if msg == "new message" {
		t.Errorf("Подписчик получил сообщение после отписки: %v", msg)
	}
}

func TestClose(t *testing.T) {
	sp := NewSubPub()

	sub, err := sp.Subscripe("sub1", func(m interface{}) {})
	if err != nil {
		t.Fatalf("Ошибка при подписке: %v", err)
	}

	err = sp.Close(context.Background())
	if err != nil {
		t.Fatalf("Ошибка при закрытии: %v", err)
	}

	_, err = sp.Subscripe("sub1", func(m interface{}) {})
	if err == nil {
		t.Error("Ожидалась ошибка при подписке после закрытия, но ошибка не произошла")
	}

	sub.Unsubscribe()
}

func TestMultipleSubs(t *testing.T) {
	sp := NewSubPub()

	var msg1, msg2 interface{}

	sub1, err := sp.Subscripe("sub1", func(m interface{}) {
		msg1 = m
	})
	if err != nil {
		t.Fatalf("Ошибка при подписке: %v", err)
	}

	sub2, err := sp.Subscripe("sub1", func(m interface{}) {
		msg2 = m
	})
	if err != nil {
		t.Fatalf("Ошибка при подписке: %v", err)
	}

	err = sp.Publish("sub1", "message")
	if err != nil {
		t.Fatalf("Ошибка при публикации: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	if msg1 != "message" {
		t.Errorf("Ожидалось 'message', а получено: %v для подписчика 1", msg1)
	}
	if msg2 != "message" {
		t.Errorf("Ожидалось 'message', а получено: %v для подписчика 2", msg2)
	}

	sub1.Unsubscribe()
	sub2.Unsubscribe()
}
