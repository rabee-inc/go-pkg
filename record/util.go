package record

import (
	"fmt"

	"github.com/rabee-inc/go-pkg/timeutil"
)

// Start ... 速度計測開始
func Start() int64 {
	return timeutil.NowUnix()
}

// End ... 速度計測終了
func End(start int64, msg string) int64 {
	end := timeutil.NowUnix()
	df := end - start
	fmt.Printf("%s %dms\n", msg, df)
	return df
}
