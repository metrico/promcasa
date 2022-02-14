package async

import (
	"context"

	"github.com/metrico/promcasa/utils/logger"

	"github.com/metrico/promcasa/model"
)

// Future interface has the method signature for await
type Future interface {
	Await() model.AsyncSqlResult
}

type future struct {
	await func(ctx context.Context) model.AsyncSqlResult
}

func (f future) Await() model.AsyncSqlResult {
	return f.await(context.Background())
}

type fn func(query string, locIndex uint, queryIndex int) model.AsyncSqlResult

// Exec executes the async function
func ExecAsyncSql(f fn, query string, locIndex uint, queryIndex int) Future {
	result := model.AsyncSqlResult{}
	c := make(chan struct{})
	go func() {
		defer close(c)
		result = f(query, locIndex, queryIndex)
	}()
	return future{
		await: func(ctx context.Context) model.AsyncSqlResult {
			select {
			case <-ctx.Done():
				logger.Error("ERROR async: ", ctx.Err())
				result.Err = ctx.Err()
				return result
			case <-c:
				if result.Err != nil {
					logger.Error("ERROR in ASYNC!", result.Err)
				}
				return result
			}
		},
	}
}
