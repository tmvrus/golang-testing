package service

import "context"

type Authorizator interface {
	Authorized(ctx context.Context, m map[string]string) (bool, error)
}
