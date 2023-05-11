package schema

import (
	"fmt"
	"time"
)

var ErrInvalidChangeset = fmt.Errorf("InvalidChangeset")

const defaultTimeout = time.Second * 30

var defaultSchema = "public"
