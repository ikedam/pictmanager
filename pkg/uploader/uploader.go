package uploader

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/disintegration/imaging"
	"github.com/ikedam/pictmanager/pkg/config"
	"github.com/ikedam/pictmanager/pkg/log"
	"github.com/ikedam/pictmanager/pkg/model"
	"github.com/ikedam/pictmanager/pkg/simplestore"
	"github.com/ikedam/pictmanager/pkg/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Uploader struct {
	config        *config.Config
	gcsBucketName string
	gcsBasePath   string
}

func New(config *config.Config) (*Uploader, error) {
	gcsURL, err := url.Parse(config.GCS)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse: %v", config.GCS)
	}
	if gcsURL.Scheme != "gs" || gcsURL.Host == "" {
		return nil, errors.Wrapf(err, "unexpected URL: %v", config.GCS)
	}
	bucketName := gcsURL.Host
	basePath := strings.TrimPrefix(gcsURL.Path, "/")
	return &Uploader{
		config:        config,
		gcsBucketName: bucketName,
		gcsBasePath:   basePath,
	}, nil
}

func (u *Uploader) Scan(ctx context.Context, dir string) error {
	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}
	bucket := gcsClient.Bucket(u.gcsBucketName)

	fsClient, err := simplestore.New(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create firestore client")
	}

	log.Debugf(ctx, "scanning %v...", dir)
	stat, err := os.Stat(dir)
	if err != nil {
		return errors.Wrapf(err, "failed to open %v", dir)
	}
	if !stat.IsDir() {
		return errors.Wrapf(err, "not a directory: %v", dir)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return errors.Wrapf(err, "failed to read %v", dir)
	}
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			err := u.Scan(ctx, path)
			if err != nil {
				return err
			}
			continue
		}
		err := u.upload(ctx, fsClient, bucket, path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Uploader) upload(ctx context.Context, fsClient *simplestore.Client, bucket *storage.BucketHandle, path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return errors.Wrapf(err, "failed to stat %v", path)
	}
	filename := filepath.Base(path)
	objectPath, err := url.JoinPath(u.gcsBasePath, filename)
	if err != nil {
		return errors.Wrapf(err, "failed to join path: %v + %v", u.gcsBasePath, filename)
	}
	thumbnailPath, err := url.JoinPath(u.gcsBasePath, "thumbnail", filename)
	if err != nil {
		return errors.Wrapf(err, "failed to join path: %v + %v", u.gcsBasePath, filename)
	}

	object := bucket.Object(objectPath)
	thumbnail := bucket.Object(thumbnailPath)

	imageInfo := &model.Image{
		ID: filename,
	}
	exists, err := fsClient.Exists(ctx, imageInfo)
	if err != nil {
		return errors.Wrapf(err, "failed to check existence of Image %v", imageInfo.ID)
	}
	if exists {
		log.Info(ctx, "skip as already exists", zap.String("localfile", objectPath))
		return nil
	}

	var imageData image.Image
	// "png", "jpeg"
	var imageFormat string
	var thumbnailBuffer bytes.Buffer
	err = func() error {
		f, err := os.Open(path)
		if err != nil {
			return errors.Wrapf(err, "failed to open %v", path)
		}
		defer f.Close()
		imageData, imageFormat, err = image.Decode(f)
		if err != nil {
			return errors.Wrapf(err, "failed to read %v", path)
		}
		// 200 x 200 に収まるようにリサイズ
		// 長い方の辺を指定する
		width := 200
		height := 200
		if imageData.Bounds().Size().X > imageData.Bounds().Size().Y {
			height = 0
		} else {
			width = 0
		}
		thumbnailData := imaging.Resize(imageData, width, height, imaging.Lanczos)
		if imageFormat == "png" {
			err := png.Encode(&thumbnailBuffer, thumbnailData)
			if err != nil {
				return errors.Wrapf(err, "failed to write thumbnail for %v", path)
			}
		} else if imageFormat == "jpeg" {
			err := jpeg.Encode(&thumbnailBuffer, thumbnailData, &jpeg.Options{
				Quality: 60, // 1-100
			})
			if err != nil {
				return errors.Wrapf(err, "failed to write thumbnail for %v", path)
			}
		} else {
			return errors.Errorf("unkown error format: %v", imageFormat)
		}
		return nil
	}()
	if err != nil {
		return nil
	}
	contentType := "image/" + imageFormat
	log.Info(
		ctx,
		"uploading",
		zap.String("localfile", path),
		zap.String("gcs_path", objectPath),
		zap.String("content_type", contentType),
	)

	err = func() error {
		w := thumbnail.NewWriter(ctx)
		defer w.Close()
		_, err := io.Copy(w, &thumbnailBuffer)
		if err != nil {
			return errors.Wrapf(err, "failed to upload thumbnail of %v to %v", path, thumbnailPath)
		}
		err = w.Close()
		if err != nil {
			return errors.Wrapf(err, "failed to write thumbnail to %v", thumbnailPath)
		}
		return nil
	}()
	if err != nil {
		return err
	}
	_, err = thumbnail.Update(ctx, storage.ObjectAttrsToUpdate{
		ContentType: contentType,
	})
	if err != nil {
		return errors.Wrapf(
			err,
			"failed to upload thumbnail of %v to %v (setting content-type)",
			path,
			objectPath,
		)
	}

	err = func() error {
		r, err := os.Open(path)
		if err != nil {
			return errors.Wrapf(err, "failed to open %v", path)
		}
		defer r.Close()
		w := object.NewWriter(ctx)
		defer w.Close()
		_, err = io.Copy(w, r)
		if err != nil {
			return errors.Wrapf(err, "failed to upload %v to %v", path, objectPath)
		}
		err = w.Close()
		if err != nil {
			return errors.Wrapf(err, "failed to write thumbnail to %v", objectPath)
		}
		return nil
	}()
	if err != nil {
		return err
	}
	_, err = object.Update(ctx, storage.ObjectAttrsToUpdate{
		ContentType: contentType,
	})
	if err != nil {
		return errors.Wrapf(
			err,
			"failed to upload %v to %v (setting content-type)",
			path,
			objectPath,
		)
	}

	now := util.GetCurrentTime()
	imageInfo = &model.Image{
		ID:          filename,
		PublishTime: stat.ModTime(),
		CreateTime:  now,
		UpdateTime:  now,
	}
	err = fsClient.Put(ctx, imageInfo)
	if err != nil {
		return errors.Wrapf(err, "failed to put Image %v", filename)
	}
	return nil
}
