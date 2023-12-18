package encoding

import (
	"encoding/json"
)

// Encoding is an implementation of a format that structures can be marshaled
// to and unmarshaled from.
type Encoding interface {
	// ContentType returns the MIME-compatible name of the encoding.
	ContentType() string

	// Encode marshals src into the encoding as bytes.
	Encode(src any) ([]byte, error)

	// Decode unmarshals src bytes into a structure using the encoding.
	Decode(src []byte, dst any) error
}

// JSON implements JavaScript Object Notation encoding.
type JSON struct{}

func (j *JSON) ContentType() string {
	return "application/json"
}

func (j *JSON) Encode(src any) ([]byte, error) {
	return json.Marshal(src)
}

func (j *JSON) Decode(src []byte, dst any) error {
	return json.Unmarshal(src, dst)
}
