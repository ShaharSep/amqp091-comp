package converter

import (
	"reflect"
	"testing"

	amqp10 "github.com/Azure/go-amqp"
	"github.com/rabbitmq/amqp091-go"
)

type message struct{}

func Test_10_091(t *testing.T) {
	original := &amqp10.Message{}
	publishing, err := ConvertTo091(original)
	if err != nil {
		t.Fatal(err)
	}
	delivery := &amqp091.Delivery{
		Headers:         publishing.Headers,
		ContentType:     publishing.ContentType,
		ContentEncoding: publishing.ContentEncoding,
		DeliveryMode:    publishing.DeliveryMode,
		Priority:        publishing.Priority,
		CorrelationId:   publishing.CorrelationId,
		ReplyTo:         publishing.ReplyTo,
		Expiration:      publishing.Expiration,
		MessageId:       publishing.MessageId,
		Timestamp:       publishing.Timestamp,
		Type:            publishing.Type,
		UserId:          publishing.UserId,
		AppId:           publishing.AppId,
		Body:            publishing.Body,
	}
	target, err := ConvertTo10(delivery)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(original, target) {
		t.Fatal("original and target are not equal")
	}
}

func Test_091_10() {

}
