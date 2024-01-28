package mongo

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/apus-run/sea-kit/zlog"
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

type Builder struct {
	log zlog.Logger
}

func NewBuilder(log zlog.Logger) *Builder {
	return &Builder{
		log,
	}
}

func (b *Builder) BuildDebugInterceptor() func(processFn) processFn {
	return func(oldProcess processFn) processFn {
		return func(cmd *cmd) error {
			beg := time.Now()
			err := oldProcess(cmd)
			cost := time.Since(beg)
			var fields = make([]zlog.Field, 0)
			if err != nil {
				fields = append(
					fields,
					zlog.String("type", "mongo.response"),
					zlog.Duration("cost", cost),
					zlog.String("data", fmt.Sprintf("%s %v", cmd.name, mustJsonMarshal(cmd.req))),
					zlog.String("error", fmt.Sprintf("%v", err.Error())),
				)
			} else {
				fields = append(
					fields,
					zlog.String("type", "mongo.response"),
					zlog.Duration("cost", cost),
					zlog.String("data", fmt.Sprintf("%s %v", cmd.name, mustJsonMarshal(cmd.req))),
					zlog.String("res", fmt.Sprintf("%v", cmd.res)),
					zlog.String("dbName", cmd.dbName),
					zlog.String("collName", cmd.collName),
					zlog.String("cmdName", cmd.name),
				)
			}

			b.log.Info("mongodb", fields...)
			return err
		}
	}
}

func mustJsonMarshal(val interface{}) string {
	res, _ := json.Marshal(val)
	return string(res)
}
