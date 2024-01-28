package mongo

import (
	"fmt"
	"log"
	"time"
)

const (
	metricType = "mongo"
)

type Interceptor func(oldProcessFn processFn) (newProcessFn processFn)

func InterceptorChain(interceptors ...Interceptor) func(oldProcess processFn) processFn {
	build := func(interceptor Interceptor, oldProcess processFn) processFn {
		return interceptor(oldProcess)
	}

	return func(oldProcess processFn) processFn {
		chain := oldProcess
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = build(interceptors[i], chain)
		}
		return chain
	}
}

func debugInterceptor(compName string, c *config) func(processFn) processFn {
	return func(oldProcess processFn) processFn {
		return func(cmd *cmd) error {
			if !eapp.IsDevelopmentMode() {
				return oldProcess(cmd)
			}

			beg := time.Now()
			err := oldProcess(cmd)
			cost := time.Since(beg)
			if err != nil {
				log.Println("emongo.response", xdebug.MakeReqAndResError(fileWithLineNum(), compName,
					fmt.Sprintf("%v", c.keyName), cost, fmt.Sprintf("%s %v", cmd.name, mustJsonMarshal(cmd.req)), err.Error()),
				)
			} else {
				log.Println("emongo.response", xdebug.MakeReqAndResInfo(fileWithLineNum(), compName,
					fmt.Sprintf("%v", c.keyName), cost, fmt.Sprintf("%s %v", cmd.name, mustJsonMarshal(cmd.req)), fmt.Sprintf("%v", cmd.res)),
				)
			}
			return err
		}
	}
}
