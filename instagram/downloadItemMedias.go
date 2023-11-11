package instagram

import (
	"errors"
	"github.com/Davincible/goinsta/v3"
)

type DownloadItemMediasResult struct {
	Id        string
	MediaType int
	Url       string
	Data      []byte
	Error     error
}

func (dimr *DownloadItemMediasResult) Size() int64 {
	if dimr.Data == nil {
		return 0
	}
	return int64(len(dimr.Data))
}

func DownloadItemMedias(item *goinsta.Item) chan DownloadItemMediasResult {
	c := make(chan DownloadItemMediasResult, 0)
	go func() {
		defer close(c)
		result := DownloadItemMediasResult{
			Id:        item.GetID(),
			MediaType: item.MediaType,
		}
		switch item.MediaType {
		case 1, 2:
			result.Url = item.Images.GetBest()
			result.Data, result.Error = item.Download()
			c <- result
		case 8:
			for _, cm := range item.CarouselMedia {
				result := DownloadItemMediasResult{
					Id:        cm.GetID(),
					MediaType: cm.MediaType,
					Url:       cm.Images.GetBest(),
				}
				if cm.MediaType != 1 && cm.MediaType != 2 {
					result.Error = errors.New(`Unknown media type`)
					ErrorLog.Println(cm.MediaType, result.Error.Error())
					c <- result
					return
				}
				result.Data, result.Error = cm.Download()
				c <- result
			}
		default:
			result.Error = errors.New(`Unknow media type`)
			ErrorLog.Println(item.MediaType, result.Error.Error())
			c <- result
		}
	}()
	return c
}
