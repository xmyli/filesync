package delta

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"io"
)

func Encode(in interface{}) ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(in)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func compress(in []byte) ([]byte, error) {
	buffer := bytes.Buffer{}
	writer := gzip.NewWriter(&buffer)
	_, err := writer.Write(in)
	if err != nil {
		return nil, err
	}
	writer.Close()
	return buffer.Bytes(), nil
}

func Decode[T interface{}](in []byte, to *T) error {
	decoder := gob.NewDecoder(bytes.NewReader(in))
	err := decoder.Decode(to)
	if err != nil {
		return err
	}
	return nil
}

func decompress(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}
	out, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	reader.Close()
	return out, nil
}
