package xid

import (
	"database/sql/driver"
	"fmt"
)

// Value implements the driver.Valuer interface.
func (id ID) Value() (driver.Value, error) {
	if id.IsNil() {
		return nil, nil
	}

	b, err := id.MarshalText()

	return string(b), err
}

// Scan implements the sql.Scanner interface.
func (id *ID) Scan(value interface{}) (err error) {
	switch val := value.(type) {
	case string:
		return id.UnmarshalText([]byte(val))
	case []byte:
		return id.UnmarshalText(val)
	case nil:
		*id = nilID
		return nil
	default:
		return fmt.Errorf("%w: %T", ErrScanUnsupportedType, value)
	}
}
