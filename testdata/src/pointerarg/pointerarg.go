package pointerarg

import (
	"database/sql"
)

func f1(a *string, tx *sql.Tx) {
}

func f2(a string, tx sql.Tx) { // want "no \\*string type arg, no \\*database/sql.Tx type arg found for func f2"
}
