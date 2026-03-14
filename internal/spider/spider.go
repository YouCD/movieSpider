package spider

import "context"

type Spider interface {
	Run(ctx context.Context)
}
