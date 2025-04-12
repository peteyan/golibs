package math

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"github.com/pierrec/lz4/v4"
)

func CompressGzip(input string) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write([]byte(input)); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecompressGzip(compressed []byte) (string, error) {
	var buf bytes.Buffer
	buf.Write(compressed)
	gz, err := gzip.NewReader(&buf)
	if err != nil {
		return "", err
	}
	defer gz.Close()
	var result bytes.Buffer
	if _, err := result.ReadFrom(gz); err != nil {
		return "", err
	}
	return result.String(), nil
}

func CompressZlib(input string) ([]byte, error) {
	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	if _, err := zw.Write([]byte(input)); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecompressZlib(compressed []byte) (string, error) {
	var buf bytes.Buffer
	buf.Write(compressed)
	zr, err := zlib.NewReader(&buf)
	if err != nil {
		return "", err
	}
	defer zr.Close()
	var result bytes.Buffer
	if _, err := result.ReadFrom(zr); err != nil {
		return "", err
	}
	return result.String(), nil
}

func CompressLz4(input string) ([]byte, error) {
	compressed := make([]byte, lz4.CompressBlockBound(len(input)))
	n, err := lz4.CompressBlock([]byte(input), compressed, nil)
	if err != nil {
		return nil, err
	}
	return compressed[:n], nil
}

func DecompressLz4(compressed []byte) (string, error) {
	decompressed := make([]byte, len(compressed)*3) // 估计解压后大小
	n, err := lz4.UncompressBlock(compressed, decompressed)
	if err != nil {
		return "", err
	}
	return string(decompressed[:n]), nil
}
