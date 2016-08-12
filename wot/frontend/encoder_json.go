package frontend

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"github.com/conas/tno2/util/str"
)

func init() {
	Encoders.Register(&JsonEncoder{})
}

type JsonEncoder struct{}

func NewJsonEncoder() *JsonEncoder {
	return &JsonEncoder{}
}

func (c *JsonEncoder) Info() string {
	return ENCODING_JSON
}

func (c *JsonEncoder) Marshal(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func (c *JsonEncoder) Unmarshal(r io.Reader, t interface{}) error {
	data, err := ioutil.ReadAll(r)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, t)

	if err != nil {
		return errors.New(str.Concat("Error unmarshaling input using ", c.Info(), " codec."))
	}

	return nil
}
