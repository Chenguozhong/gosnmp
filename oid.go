package gosnmp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// A bunch of commonly used MIB-2 oids.
var (
	SYS_DESCRoid     = ObjectIdentifier{1, 3, 6, 1, 2, 1, 1, 1, 0}
	SYS_OBJECT_IDoid = ObjectIdentifier{1, 3, 6, 1, 2, 1, 1, 2, 0}
	SYS_UPTIMEoid    = ObjectIdentifier{1, 3, 6, 1, 2, 1, 1, 3, 0}
	SYS_CONTACToid   = ObjectIdentifier{1, 3, 6, 1, 2, 1, 1, 4, 0}
	SYS_NAMEoid      = ObjectIdentifier{1, 3, 6, 1, 2, 1, 1, 5, 0}
	SYS_LOCATIONoid  = ObjectIdentifier{1, 3, 6, 1, 2, 1, 1, 6, 0}
)

func parseOid(oidString string) (oid []int, err error) {
	ids := strings.Split(oidString, ".")
	if len(ids) < 2 {
		return nil, errors.New(fmt.Sprintf("The object identifer \"%s\" doesn't contain at least 2 sub identifiers", oidString))
	}
	oid = make([]int, len(ids))
	for i := 0; i < len(ids); i++ {
		if oid[i], err = strconv.Atoi(ids[i]); err != nil {
			return nil, errors.New(fmt.Sprintf("Sub identifier %d in \"%s\" couldn't be parsed", i+1, oidString))
		}
	}
	return oid, nil
}
