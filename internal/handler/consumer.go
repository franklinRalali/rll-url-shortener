// Package handler
package handler

import (
	"context"

	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/internal/consts"
	uContract "github.com/ralali/rll-url-shortener/internal/ucase/contract"
	"github.com/ralali/rll-url-shortener/pkg/awssqs"
)

// SQSConsumerHandler sqs consumer message processor handler
func SQSConsumerHandler(msgHandler uContract.MessageProcessor) awssqs.MessageProcessorFunc {
	return func(decoder *awssqs.MessageDecoder) error {
		return msgHandler.Serve(context.Background(), &appctx.ConsumerData{
			Body:        []byte(*decoder.Body),
			Key:         []byte(*decoder.MessageId),
			ServiceType: consts.ServiceTypeConsumer,
		})
	}
}
