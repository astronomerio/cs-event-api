package types

import(
	"strconv"
)

type ConvertibleString string

func (str *ConvertibleString) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		*str = ConvertibleString(data)
	} else {
		*str = ConvertibleString(s)
	}

	return nil
}
