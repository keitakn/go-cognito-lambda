package infrastructure

import (
	"flag"
	"strconv"
	"testing"
)

func IsTestRun() bool {
	testing.Init()
	flag.Parse()
	flg := flag.Lookup("test.v")
	if flg != nil {
		test, err := strconv.ParseBool(flg.Value.String())
		if err != nil {
			panic(err.Error())
		}
		if test {
			return true
		}
	}

	return false
}
