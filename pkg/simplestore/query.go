package simplestore

import (
	"context"
	"reflect"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

type Query struct {
	firestore.Query
	c          *Client
	targetList reflect.Value
	targetType reflect.Type
}

func (c *Client) Query(target any) (*Query, error) {
	v := reflect.ValueOf(target)
	if v.Type().Kind() != reflect.Pointer {
		return nil, errors.WithMessagef(ErrProgramming, "target must be a pointer to slice of pointer to struct(*[]*Type), but %T", target)
	}
	if v.IsNil() {
		return nil, errors.WithMessage(ErrProgramming, "target must not be nil")
	}
	v = v.Elem()
	if v.Type().Kind() != reflect.Slice {
		return nil, errors.WithMessagef(ErrProgramming, "target must be a pointer to slice of pointer to struct(*[]*Type), but %T", target)
	}
	t := v.Type().Elem()
	if t.Kind() != reflect.Pointer {
		return nil, errors.WithMessagef(ErrProgramming, "target must be a pointer to slice of pointer to struct(*[]*Type), but %T", target)
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return nil, errors.WithMessagef(ErrProgramming, "target must be a pointer to slice of pointer to struct(*[]*Type), but %T", target)
	}
	typeName := t.Name()
	collection := c.FirestoreClient.Collection(typeName)
	return &Query{
		Query:      collection.Query,
		c:          c,
		targetList: v,
		targetType: t,
	}, nil
}

func (c *Client) MustQuery(target any) *Query {
	q, err := c.Query(target)
	if err != nil {
		panic(err)
	}
	return q
}

func (q *Query) OrderBy(path string, dir firestore.Direction) *Query {
	newQ := *q
	newQ.Query = q.Query.OrderBy(path, dir)
	return &newQ
}

func (q *Query) Limit(n int) *Query {
	newQ := *q
	newQ.Query = q.Query.Limit(n)
	return &newQ
}

func (q *Query) Do(ctx context.Context) error {
	iter := q.Query.Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		v := reflect.New(q.targetType)
		err = doc.DataTo(v.Interface())
		if err != nil {
			return errors.WithMessagef(ErrProgramming, "failed to serialize data to %v:  %+v", q.targetType.Name, doc.Data(), err)
		}
		q.targetList.Set(reflect.Append(q.targetList, v))
	}
	return nil
}

func (q *Query) DoAfterByID(ctx context.Context, id string) error {
	doc := q.c.FirestoreClient.Collection(q.targetType.Name()).Doc(id)
	snapshot, err := doc.Get(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to stat %v/%v", q.targetType.Name(), id)
	}
	newQ := *q
	newQ.Query = q.Query.StartAfter(snapshot)
	return newQ.Do(ctx)
}
