package types

import (
	"fmt"
	"io"
	"strconv"
)

type ID uint64

func (id ID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

func (id *ID) UnmarshalGQL(v interface{}) error {
	idstr, ok := v.(string)
	if !ok {
		return fmt.Errorf("id must be a string")
	}

	u, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to ParseUint: %w", err)
	}

	*id = ID(u)

	return nil
}

func (id ID) MarshalGQL(w io.Writer) {
	_, _ = fmt.Fprint(w, strconv.Quote(id.String()))
}

func (id ID) Raw() interface{} {
	return uint64(id)
}
