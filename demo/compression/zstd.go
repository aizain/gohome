package compression

import "github.com/klauspost/compress/zstd"

func zstdEncode(str []byte) []byte {
	var encoder, _ = zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.EncoderLevelFromZstd(3)))
	return encoder.EncodeAll(str, make([]byte, 0, len(str)))
}

func zstdDecode(str []byte) ([]byte, error) {
	var decoder, _ = zstd.NewReader(nil)
	return decoder.DecodeAll(str, nil)
}
