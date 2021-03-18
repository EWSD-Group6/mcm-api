package worker

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"mcm-api/pkg/contribution"
	"mcm-api/pkg/log"
	"mcm-api/pkg/queue"
	"os"
	"strconv"
	"strings"
	"sync"
)

func (w worker) exportContributeSessionHandler(ctx context.Context, message *queue.Message) error {
	temp, err := os.MkdirTemp("", "mcm-processing-*")
	if err != nil {
		return err
	}
	defer func() {
		err = os.Remove(temp)
		if err != nil {
			log.Logger.Error("failed to delete folder", zap.String("name", temp))
		}
	}()
	if v, ok := message.Data.(*queue.ExportContributeSessionPayload); ok {
		contributions, err := w.contributionService.GetAllAcceptedContributions(ctx, v.ContributeSessionId)
		if err != nil {
			return err
		}
		waitGroup := sync.WaitGroup{}
		waitGroup.Add(len(contributions))
		for _, contrib := range contributions {
			go func(c *contribution.Entity) {
				e := w.downloadContribution(ctx, temp, c)
				if e != nil {
					log.Logger.Error("download contribution failed", zap.Error(e))
				}
				waitGroup.Done()
			}(contrib)
		}
		waitGroup.Wait()

	} else {
		return errors.New("unknown message")
	}
}

func (w worker) downloadContribution(ctx context.Context, basePath string, c *contribution.Entity) error {
	version, err := w.articleService.GetLatestVersionOfArticle(ctx, *c.ArticleId)
	if err != nil {
		return err
	}
	contributionFolder := basePath + "/" + strconv.Itoa(c.Id)
	err = os.Mkdir(contributionFolder, fs.ModeDir)
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
		_ := articleReader.Close()
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

		}
		_ = imageReader.Close()
	}
}
