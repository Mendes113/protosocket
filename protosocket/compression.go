package protosocket

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"

	"google.golang.org/protobuf/proto"
)

type CompressionType int

const (
	NoCompression CompressionType = iota
	GzipCompression
)

func compressData(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	if _, err := writer.Write(data); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decompressData(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

// CompressionMiddleware aplica compressão nas mensagens
func CompressionMiddleware() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, msg proto.Message) error {
			data, err := proto.Marshal(msg)
			if err != nil {
				return err
			}

			// Comprime os dados
			compressed, err := compressData(data)
			if err != nil {
				return err
			}

			// Descomprime antes de passar adiante
			decompressed, err := decompressData(compressed)
			if err != nil {
				return err
			}

			// Reconstrói a mensagem
			if err := proto.Unmarshal(decompressed, msg); err != nil {
				return err
			}

			return next(ctx, msg)
		}
	}
}
