package mongo

import (
	"context"
	"github.com/soukengo/gopkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	mgoopts "go.mongodb.org/mongo-driver/mongo/options"
	"rockimserver/apis/rockim/shared"
)

var (
	DocumentNotFound = errors.NotFound("DATA_NOT_FOUND", "data not found")
)

func (c *Client) FindOne(ctx context.Context, collection string, filter any, dist any, opts ...*mgoopts.FindOneOptions) (err error) {
	err = c.Collection(collection).FindOne(ctx, filter, opts...).Decode(dist)
	if IsErrNoDocuments(err) {
		err = DocumentNotFound
	}
	return
}

func (c *Client) FindList(ctx context.Context, collection string, filter any, result any, opts ...*mgoopts.FindOptions) (err error) {
	cursor, err := c.Collection(collection).Find(ctx, filter, opts...)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)
	return cursor.All(ctx, result)
}

func (c *Client) Paginate(ctx context.Context, collection string, query interface{}, paginate *shared.Paginating, opts ...*mgoopts.FindOptions) (cursor *mongo.Cursor, p *shared.Paginated, err error) {
	offset := paginate.Offset()
	limit := paginate.Limit()
	opts = append(opts, &mgoopts.FindOptions{
		Limit: &limit,
		Skip:  &offset,
	})
	co := c.Collection(collection)
	total, err := co.CountDocuments(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	p = &shared.Paginated{Total: total}
	cursor, err = co.Find(ctx, query, opts...)
	return
}

func ScanCursor[T any](ctx context.Context, cursor *mongo.Cursor) (results []*T, err error) {
	for cursor.Next(ctx) {
		var record = new(T)
		if err = cursor.Decode(record); err != nil {
			return
		}
		results = append(results, record)
	}
	return
}
