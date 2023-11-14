package instagram

import (
	"errors"
	"fmt"
	"github.com/Davincible/goinsta/v3"
	"time"
)

func DownloadProfile(name string, client *goinsta.Instagram, storage Storage, after time.Time, delay time.Duration) (int, int, int64, error) {
	if name == `` || client == nil || storage == nil {
		return 0, 0, 0, errors.New(`Bad input`)
	}
	profile, err := client.Profiles.ByName(name)
	if err != nil {
		ErrorLog.Println(name, err.Error())
		return 0, 0, 0, err
	}
	postsCount := 0
	mediaItemsCount := 0
	var mediaItemsSize int64
	feed := profile.Feed()
	pagination := true
	for pagination {
		if !feed.Next() {
			break
		}
		if err := feed.Error(); err != nil && err.Error() == `no more posts availible, page end has been reached` {
			pagination = false
		}
		alreadyDownloadedPosts := 0
		DebugLog.Printf("status: %s, more available: %t, items: %d, error: %v", feed.Status, feed.MoreAvailable, len(feed.Items), feed.Error())
		for _, post := range feed.Items {
			createdAt := time.Unix(post.Caption.CreatedAtUtc, 0)
			if createdAt.Before(after) {
				pagination = false
				break
			}
			postId := post.GetID()
			mediaItemsIds := make([]string, 0)
			DebugLog.Printf("Post ID: %s, type: %d, Created at: %s, Location: %v, Text: %s", postId, post.MediaType, createdAt, post.Location, post.Caption.Text)
			if storage.IsExist(name, postId) {
				DebugLog.Println(`Already downloaded`)
				alreadyDownloadedPosts++
				continue
			}
			time.Sleep(delay)
			switch post.MediaType {
			case 8:
				for _, cm := range post.CarouselMedia {
					mediaId := fmt.Sprintf("%v", cm.ID)
					mediaType := `PHOTO`
					if cm.MediaType == 2 {
						mediaType = `VIDEO`
					}
					if cm.MediaType != 1 && cm.MediaType != 2 {
						ErrorLog.Printf(`Unknown media type: %d`, cm.MediaType)
						mediaType = `UNKNOWN`
					}
					DebugLog.Printf("ID: %s, Type: %s, URL: %s", mediaId, mediaType, cm.Images.GetBest())
					if mediaType == `UNKNOWN` {
						return postsCount, mediaItemsCount, mediaItemsSize, err
					}
					data, err := cm.Download()
					if err != nil {
						ErrorLog.Println(err.Error())
						return postsCount, mediaItemsCount, mediaItemsSize, err
					}
					if storage.SaveMedia(name, postId, mediaId, mediaType, createdAt, data) != nil {
						return postsCount, mediaItemsCount, mediaItemsSize, err
					}
					mediaItemsCount++
					mediaItemsSize = mediaItemsSize + int64(len(data))
					mediaItemsIds = append(mediaItemsIds, mediaId)
				}
			case 1:
				DebugLog.Printf("ID: %s, Type: PHOTO, URL: %s", postId, post.Images.GetBest())
				data, err := post.Download()
				if err != nil {
					ErrorLog.Println(err.Error())
					return postsCount, mediaItemsCount, mediaItemsSize, err
				}
				if storage.SaveMedia(name, postId, postId, `PHOTO`, createdAt, data) != nil {
					return postsCount, mediaItemsCount, mediaItemsSize, err
				}
				mediaItemsCount++
				mediaItemsSize = mediaItemsSize + int64(len(data))
				mediaItemsIds = append(mediaItemsIds, postId)
			case 2:
				DebugLog.Printf("ID: %s, Type: VIDEO, URL: %s", postId, post.Images.GetBest())
				data, err := post.Download()
				if err != nil {
					ErrorLog.Println(err.Error())
					return postsCount, mediaItemsCount, mediaItemsSize, err
				}
				if storage.SaveMedia(name, postId, postId, `VIDEO`, createdAt, data) != nil {
					return postsCount, mediaItemsCount, mediaItemsSize, err
				}
				mediaItemsCount++
				mediaItemsSize = mediaItemsSize + int64(len(data))
				mediaItemsIds = append(mediaItemsIds, postId)
			default:
				err = errors.New(`Unknow type`)
				ErrorLog.Println(err.Error())
				return postsCount, mediaItemsCount, mediaItemsSize, err
			}
			if err := storage.SavePost(name, postId, createdAt, post.Location, mediaItemsIds, post.Caption.Text); err != nil {
				return postsCount, mediaItemsCount, mediaItemsSize, err
			}
		}
		if pagination && alreadyDownloadedPosts == feed.NumResults {
			pagination = false
		}
	}
	time.Sleep(delay)
	storyMedia, err := profile.Stories()
	if err != nil {
		ErrorLog.Println(err.Error())
		return 0, mediaItemsCount, mediaItemsSize, nil
	}
	DebugLog.Println(storyMedia.Status, storyMedia.Broadcasts, len(storyMedia.Reel.MediaIDs))
	for _, reelMediaId := range storyMedia.Reel.MediaIDs {
		time.Sleep(delay)
		feedMedia, err := client.GetMedia(reelMediaId)
		if err != nil {
			ErrorLog.Println(err.Error())
			continue
		}
		DebugLog.Printf("Reel Media ID: %d\tItems:%d", reelMediaId, len(feedMedia.Items))
		for _, item := range feedMedia.Items {
			itemId := item.GetID()
			DebugLog.Printf("%s\tItem ID:%s\tType:%d\tSeen:%t", item.User.Username, itemId, item.MediaType, item.IsSeen)
			if storage.IsExist(item.User.Username, item.GetID()) {
				DebugLog.Println(`Already downloaded`)
				continue
			}
			n, s, err := DownloadItem(item, storage, time.Time{})
			if err != nil {
				ErrorLog.Println(err.Error())
				continue
			}
			DebugLog.Printf("Got %d items (%d bytes)", n, s)
		}
	}

	return 0, mediaItemsCount, mediaItemsSize, nil
}
