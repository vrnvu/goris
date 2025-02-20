package common

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const PROTOCOL_HEADER int = 4
const MESSAGE_MAX_SIZE int = 4096
const MESSAGE_LENGTH int = PROTOCOL_HEADER + MESSAGE_MAX_SIZE + 1

func ReadRequestFromBuffer(buffer *bytes.Buffer) (data []byte, remaining bool, err error) {

	header := make([]byte, PROTOCOL_HEADER)
	_, err = buffer.Read(header)
	if err != nil {
		return data, remaining, err
	}

	contentLength := binary.LittleEndian.Uint32(header)
	if contentLength <= 0 {
		return data, remaining, errors.New("no requests")
	}
	data = make([]byte, contentLength)

	_, err = buffer.Read(data)
	if err != nil {
		return data, remaining, err
	}

	rv, err := buffer.ReadByte()
	if rv == 0 || err != nil {
		remaining = false
	} else {
		remaining = true
		buffer.UnreadByte()
	}

	// Return the fully formed response; [header][data].
	return append(header, data...), remaining, nil

}

func ReadFromBuffer(buffer *bytes.Buffer) (string, error) {

	data, _, err := ReadRequestFromBuffer(buffer)
	if err != nil {
		return "", err
	}

	// Skip the header; we only want the payload.
	return string(data[PROTOCOL_HEADER:]), nil
}

func AppendToBuffer(msg string, buffer *bytes.Buffer) (*bytes.Buffer, error) {
	if len(msg) > MESSAGE_MAX_SIZE {
		return nil, fmt.Errorf("message length was %d, max size is: %d", len(msg), MESSAGE_MAX_SIZE)
	}

	if buffer == nil {
		buffer = bytes.NewBuffer(make([]byte, 0, PROTOCOL_HEADER+len((msg))))
	}

	header := make([]byte, PROTOCOL_HEADER)
	binary.LittleEndian.PutUint32(header, uint32(len(msg)))

	_, err := buffer.Write(header)
	if err != nil {
		return nil, err
	}

	_, err = buffer.WriteString(msg)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
