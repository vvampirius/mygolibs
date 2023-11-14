package instagram

import (
	"errors"
	"fmt"
	"github.com/Davincible/goinsta/v3"
	"time"
)

var (
	ErrIsAd              = errors.New(`is AD`)
	ErrIsNotAfter        = errors.New(`is not after`)
	ErrAlreadyDownloaded = errors.New(`already downloaded`)
)

func DownloadItem(item *goinsta.Item, storage Storage, after time.Time) (int, int64, error) {
	postId := item.GetID()
	mediaItemsIds := make([]string, 0)
	isAd := false
	DebugLog.Println(item.AdLink, item.AdLinkType, item.AdLinkHint, item.AdAction, item.AdMetadata, item.AdText, item.AdTitle)
	if item.AdLink != `` {
		isAd = true
	}
	createdAt := time.Unix(item.Caption.CreatedAtUtc, 0)
	if createdAt.Before(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)) {
		createdAt = time.Now()
	}
	fmt.Printf("User: %s, Post ID: %s, type: %d, Created at: %s, Location: %v, AD: %t\n", item.User.Username, postId, item.MediaType, createdAt, item.Location, isAd)
	if createdAt.Before(after) {
		return 0, 0, ErrIsNotAfter
	}
	if isAd {
		return 0, 0, ErrIsAd
	}
	if storage.IsExist(item.User.Username, postId) {
		return 0, 0, ErrAlreadyDownloaded
	}
	var mediaItemsCount int
	var mediaItemsSize int64
	c := DownloadItemMedias(item)
	for {
		x, open := <-c
		if !open {
			break
		}
		if x.Error != nil {
			ErrorLog.Println(x.Error)
			return mediaItemsCount, mediaItemsSize, x.Error
		}
		mediaType := `PHOTO`
		if x.MediaType == 2 {
			mediaType = `VIDEO`
		}
		if err := storage.SaveMedia(item.User.Username, postId, x.Id, mediaType, createdAt, x.Data); err != nil {
			return mediaItemsCount, mediaItemsSize, x.Error
		}
		mediaItemsCount++
		mediaItemsSize = mediaItemsSize + x.Size()
		mediaItemsIds = append(mediaItemsIds, x.Id)
	}
	if err := storage.SavePost(item.User.Username, postId, createdAt, item.Location, mediaItemsIds, item.Caption.Text); err != nil {
		return mediaItemsCount, mediaItemsSize, err
	}
	return mediaItemsCount, mediaItemsSize, nil
}
