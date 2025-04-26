package converter

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	amqp10 "github.com/Azure/go-amqp"
	amqp091 "github.com/rabbitmq/amqp091-go"
)

const (
	propertiesHeaderName    = "x-amqp-1.0-properties"
	appPropertiesHeaderName = "x-amqp-1.0-app-properties"
	appPropertiesPrefix     = "\x00St"
	bodyDelimiter           = "\x00SU"
	codeBin8                = '\xa0'
	codeBin32               = '\xb0'
)

func ConvertTo10(original *amqp091.Delivery) (*amqp10.Message, error) {
	payload := []byte{}
	propertirsHeader := original.Headers[propertiesHeaderName]
	if propertirsHeader == nil {
		return nil, fmt.Errorf("error %s header missing", propertiesHeaderName)
	}
	appPropertirsHeader := original.Headers[appPropertiesHeaderName]
	if appPropertirsHeader == nil {
		return nil, fmt.Errorf("error %s header missing", appPropertiesHeaderName)
	}
	propertiesString := convertHeader(propertirsHeader)
	appPropertiesString := convertHeader(appPropertirsHeader)

	payload = append(payload, propertiesString...)
	payload = append(payload, appPropertiesString...)

	if len(original.Body) > 0 {
		payload = append(payload, encodeBodySize(original.Body)...)
	}

	result := &amqp10.Message{}
	if err := result.UnmarshalBinary(payload); err != nil {
		return nil, fmt.Errorf("error could not UnMarshal message: %s", err.Error())
	}

	return result, nil
}

func ConvertTo091(original *amqp10.Message) (*amqp091.Publishing, error) {
	var buffer []byte
	var err error
	if buffer, err = original.MarshalBinary(); err != nil {
		return nil, err
	}

	body, headerSize := extractBody(buffer)

	bufferStr := string(buffer[:headerSize])
	var encodingIndex int
	if encodingIndex = strings.Index(bufferStr, appPropertiesPrefix); encodingIndex < 0 {
		return nil, fmt.Errorf("coulde not find encoding in headers: %s", bufferStr)
	}

	headers := make(map[string]interface{})
	headers[propertiesHeaderName] = buffer[:encodingIndex]
	headers[appPropertiesHeaderName] = buffer[encodingIndex:headerSize]
	result := &amqp091.Publishing{
		Headers:         headers,
		ContentType:     "applicationJSON",
		ContentEncoding: "",
		DeliveryMode:    0,
		Priority:        0,
		Timestamp:       time.Now().UTC(),
		Body:            body,
	}
	return result, nil
}

func convertHeader(header interface{}) []byte {
	output := ""
	switch reflect.ValueOf(header).Kind() {
	case reflect.Array:
	case reflect.Slice:
		output = string(header.([]uint8))
	case reflect.String:
		output = header.(string)
	}
	return []byte(output)
}

func extractBody(data []byte) ([]byte, int) {
	dataLength := len(data)
	dataStr := string(data)
	headerSize := strings.Index(dataStr, bodyDelimiter)
	if headerSize < 0 {
		return nil, dataLength
	}

	var bodySize int64
	var bodyIndex int
	bodySegment := data[headerSize+len(bodyDelimiter):]
	if len(bodySegment) > 2 && bodySegment[0] == codeBin8 { // type8Bit
		bodyIndex = 1
		bodySize = int64(bodySegment[1])
	} else if len(bodySegment) > 4 && bodySegment[0] == codeBin32 { // type32Bit
		bodyIndex = 4
		bodySize = int64(binary.BigEndian.Uint32(bodySegment[1 : bodyIndex+1]))
	} else {
		return nil, dataLength
	}

	body := bodySegment[bodyIndex+1:]
	if int64(len(body)) != bodySize {
		return nil, dataLength
	}
	return body, bodyIndex
}

func encodeBodySize(body []byte) []byte {
	result := []byte(bodyDelimiter)
	size := len(body)
	if size <= math.MaxUint8 {
		result = append(result, codeBin8, uint8(size))
	} else {
		result = append(result, codeBin32)
		result = binary.BigEndian.AppendUint32(result, uint32(size))
	}

	result = append(result, body...)
	return result
}
