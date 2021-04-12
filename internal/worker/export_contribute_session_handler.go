package worker

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"mcm-api/pkg/contribution"
	"mcm-api/pkg/log"
	"mcm-api/pkg/media"
	"mcm-api/pkg/queue"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

func (w worker) exportContributeSessionHandler(ctx context.Context, message *queue.Message) error {
	v, ok := message.Data.(*queue.ExportContributeSessionPayload)
	if !ok {
		return errors.New("unknown message")
	}
	mutex := w.CreateMutex(v.ContributeSessionId)
	// try to lock
	err := mutex.LockContext(ctx)
	if err != nil {
		return err
	}
	// lock oke
	defer func() {
		ok, err = mutex.UnlockContext(ctx)
		if !ok || err != nil {
			log.Logger.Error("failed to unlock", zap.Error(err))
		}
	}()

	// create temp folder to processing
	temp, err := os.MkdirTemp("", "mcm-processing-*")
	if err != nil {
		return err
	}
	defer func() {
		err = os.RemoveAll(temp)
		if err != nil {
			log.Logger.Error("failed to delete folder", zap.String("name", temp))
		}
	}()
	contributions, err := w.contributionService.GetAllAcceptedContributions(ctx, v.ContributeSessionId)
	if err != nil {
		return err
	}
	if len(contributions) == 0 {
		log.Logger.Info("contribution session dont have any contributions", zap.Int("id", v.ContributeSessionId))
		return nil
	} else {
		log.Logger.Info(fmt.Sprintf("found %v contributions", len(contributions)))
	}
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(contributions))
	for _, contrib := range contributions {
		go func(c *contribution.Entity) {
			e := w.downloadContribution(ctx, temp, c)
			if e != nil {
				log.Logger.Error("download contribution failed", zap.Error(e))
			} else {
				log.Logger.Info("finish download contribution", zap.Int("id", c.Id))
			}
			waitGroup.Done()
		}(contrib)
	}
	waitGroup.Wait()
	zipFile, err := recursiveZip(temp)
	if err != nil {
		return err
	}
	_, _ = zipFile.Seek(0, 0)
	defer func() {
		_ = zipFile.Close()
		_ = os.Remove(zipFile.Name())
	}()
	uploadResult, err := w.mediaService.UploadContribution(ctx, &media.ContributionUploadReq{
		File:                zipFile,
		ContributeSessionId: v.ContributeSessionId,
	})
	if err != nil {
		return err
	}
	return w.contributionSessionService.UpdateExportedAsset(
		ctx, v.ContributeSessionId, uploadResult.Key)
}

func recursiveZip(basePath string) (*os.File, error) {
	zipfile, err := os.CreateTemp("", "*.zip")
	if err != nil {
		return nil, err
	}

	w := zip.NewWriter(zipfile)
	defer func() {
		_ = w.Close()
	}()

	walker := func(path string, info os.FileInfo, err error) error {
		log.Logger.Debug("current path", zap.String("path", path))
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			_ = file.Close()
		}()

		// Ensure that `path` is not absolute; it should not start with "/".
		// transforms path into a zip-root relative path.
		f, err := w.Create(strings.TrimPrefix(path, basePath+"/"))
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	err = filepath.Walk(basePath, walker)
	if err != nil {
		return nil, err
	}
	return zipfile, nil
}

func (w worker) downloadContribution(ctx context.Context, basePath string, c *contribution.Entity) error {
	version, err := w.articleService.GetLatestVersionOfArticle(ctx, *c.ArticleId)
	if err != nil {
		return err
	}
	contributionFolder := basePath + "/" + strconv.Itoa(c.Id)
	err = os.Mkdir(contributionFolder, fs.ModePerm)
	if err != nil {
		return err
	}
	articleFile, err := os.Create(contributionFolder + "/article." + strings.Split(version.LinkOriginal, ".")[1])
	if err != nil {
		return err
	}
	articleReader, err := w.mediaService.GetFile(ctx, version.LinkOriginal)
	if err != nil {
		return err
	}
	defer func() {
		_ = articleReader.Close()
	}()
	_, err = io.Copy(articleFile, articleReader)
	if err != nil {
		return err
	}
	if len(c.Images) == 0 {
		return nil
	}
	for _, image := range c.Images {
		imageWriter, er := os.Create(contributionFolder + "/" + image.Key)
		if er != nil {
			return er
		}
		imageReader, er := w.mediaService.GetFile(ctx, image.Key)
		if er != nil {
			return er
		}
		_, er = io.Copy(imageWriter, imageReader)
		if er != nil {
			return er
		}
		_ = imageReader.Close()
	}
	return nil
}

func (w worker) CreateMutex(id int) *redsync.Mutex {
	return w.lock.NewMutex(generateLockKey(id),
		redsync.WithExpiry(time.Hour),
		redsync.WithTries(2),
	)
}

func generateLockKey(id int) string {
	return fmt.Sprintf("session:%v:export-asset-lock", id)
}
