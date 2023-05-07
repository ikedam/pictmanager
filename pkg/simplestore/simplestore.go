package simplestore

import (
	"context"
	stderrors "errors"
	"os"
	"reflect"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrProgramming = stderrors.New("programing error")
)

type Client struct {
	FirestoreClient *firestore.Client
}

type SimpleStoreLoader interface {
	OnLoadFromFirestore() error
}

type SimpleStoreSaver interface {
	OnSaveToFirestore() error
}

var projectEnvNameList = []string{
	"CLOUDSDK_CORE_PROJECT",
	"GOOGLE_CLOUD_PROJECT",
	"GCP_PROJECT",
	"GCLOUD_PROJECT",
	"PROJECT_ID",
}

func getProjectID() string {
	for _, name := range projectEnvNameList {
		projectID := os.Getenv(name)
		if projectID != "" {
			return projectID
		}
	}
	return ""
}

func New(ctx context.Context) (*Client, error) {
	client, err := firestore.NewClient(ctx, getProjectID())
	if err != nil {
		return nil, err
	}
	return &Client{
		FirestoreClient: client,
	}, nil
}

func (c *Client) getDoc(ctx context.Context, target any) (*firestore.DocumentRef, bool, error) {
	v := reflect.ValueOf(target)
	if v.Type().Kind() != reflect.Pointer {
		return nil, false, errors.WithMessagef(ErrProgramming, "target must be a pointer to a struct, but %T", target)
	}
	if v.IsNil() {
		return nil, false, errors.WithMessage(ErrProgramming, "target must not be nil")
	}
	v = v.Elem()
	if v.Type().Kind() != reflect.Struct {
		return nil, false, errors.WithMessagef(ErrProgramming, "target must be a pointer to a struct, but %T", target)
	}
	typeName := v.Type().Name()
	idField, ok := v.Type().FieldByName("ID")
	if !ok {
		return nil, false, errors.WithMessagef(ErrProgramming, "failed to stat ID field for %v", typeName)
	}
	if idField.Type.Kind() != reflect.String {
		return nil, false, errors.WithMessagef(ErrProgramming, "ID field of %v is not a string", typeName)
	}
	id := v.FieldByName(idField.Name).String()
	if id == "" {
		return c.FirestoreClient.Collection(typeName).NewDoc(), true, nil
	}
	return c.FirestoreClient.Collection(typeName).Doc(id), false, nil
}

func (c *Client) getSnapshot(ctx context.Context, target any) (*firestore.DocumentSnapshot, error) {
	doc, _, err := c.getDoc(ctx, target)
	if err != nil {
		return nil, err
	}
	return doc.Get(ctx)
}

func (c *Client) Get(ctx context.Context, target any) error {
	snapshot, err := c.getSnapshot(ctx, target)
	if err != nil {
		return err
	}
	err = snapshot.DataTo(target)
	if err != nil {
		return errors.WithMessagef(ErrProgramming, "failed to serialize data to %T from %+v: %+v", target, snapshot.Data(), err)
	}
	loader, ok := target.(SimpleStoreLoader)
	if ok {
		err := loader.OnLoadFromFirestore()
		if err != nil {
			return errors.Wrap(err, "failed in OnLoadFromFirestore")
		}
	}
	return nil
}

func (c *Client) Exists(ctx context.Context, target any) (bool, error) {
	snapshot, err := c.getSnapshot(ctx, target)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false, err
	}
	return snapshot.Exists(), nil
}

func (c *Client) Put(ctx context.Context, target any) error {
	doc, isNew, err := c.getDoc(ctx, target)
	if err != nil {
		return err
	}
	if isNew {
		v := reflect.ValueOf(target).Elem()
		v.FieldByName("ID").Set(reflect.ValueOf(doc.ID))
	}
	saver, ok := target.(SimpleStoreSaver)
	if ok {
		err := saver.OnSaveToFirestore()
		if err != nil {
			return errors.Wrap(err, "failed in OnSaveToFirestore")
		}
	}
	_, err = doc.Set(ctx, target)
	if err != nil && isNew {
		v := reflect.ValueOf(target).Elem()
		v.FieldByName("ID").Set(reflect.ValueOf(""))
	}
	return err
}
