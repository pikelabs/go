package rpmutil

import (
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"

	"code.pikelabs.net/go/compress/xz"
	"code.pikelabs.net/go/rpm"
)

const (
	plUncompressed = "uncompressed"
	plGzip         = "gzip"
	plBzip2        = "bzip2"
	plXZ           = "xz"
)

func decompressPkgPayload(p *Package) (io.Reader, error) {
	compressor, err := p.Header.GetString(rpm.TagPayloadCompressor)

	if err != nil {
		return nil, err
	}

	switch compressor {
	case plGzip:
		return gzip.NewReader(p.r)
	case plBzip2:
		return bzip2.NewReader(p.r), nil
	case plXZ:
		return xz.NewReader(p.r), nil
	}
	return nil, fmt.Errorf("unsuppored compression format: %s", compressor)
}
