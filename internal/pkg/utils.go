package pkg

import "github.com/oklog/ulid/v2"


// ULIDToString returns the canonical string representation of a ULID.
func ULIDToString(id ulid.ULID) string {
	return id.String()
}

// StringToULID parses a ULID string and returns the ULID object.
// Returns an error if the string is invalid.
func StringToULID(s string) (ulid.ULID, error) {
	return ulid.Parse(s)
}