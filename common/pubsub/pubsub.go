package pubsub

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

func PublishProtoMessages(projectID string, topicID string, payload protoreflect.ProtoMessage) error {
	log.Printf("PublishProtoMessages: projectID=%v, topicID=%v", projectID, topicID)

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Printf("pubsub.NewClient: %v", err)
		return err
	}

	t := client.Topic(topicID)
	cfg, err := t.Config(ctx)
	if err != nil {
		log.Printf("topic.Config err: %v", err)
		return err
	}
	encoding := cfg.SchemaSettings.Encoding

	var msg []byte
	switch encoding {
	case pubsub.EncodingBinary:
		msg, err = proto.Marshal(payload)
		if err != nil {
			log.Printf("proto.Marshal err: %v", err)
			return err
		}
	case pubsub.EncodingJSON:
		pjson := protojson.MarshalOptions{
			EmitUnpopulated: true,
		}
		msg, err = pjson.Marshal(payload)
		if err != nil {
			log.Printf("protojson.Marshal err: %v", err)
			return err
		}
	default:
		log.Printf("invalid encoding: %v", encoding)
		return err
	}

	result := t.Publish(ctx, &pubsub.Message{
		Data: msg,
	})
	_, err = result.Get(ctx)
	if err != nil {
		log.Printf("result.Get: %v", err)
		return err
	}

	log.Printf("Published message: %v encoding: %v", encoding, string(msg))
	return nil
}
