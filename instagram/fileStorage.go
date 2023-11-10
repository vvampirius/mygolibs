package instagram

import (
	"github.com/Davincible/goinsta/v3"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"time"
)

type Storage interface {
	IsExist(profileName string, postId string) bool
	SavePost(profileName string, postId string, createdAt time.Time, location goinsta.Location, mediaIds []string, text string) error
	SaveMedia(profileName string, postId string, mediaId string, mediaType string, createdAt time.Time, data []byte) error
}

type FileStorageProfilePost struct {
	ID        string           `yaml:"id"`
	CreatedAt time.Time        `yaml:"created_at"`
	Location  goinsta.Location `yaml:"location"`
	Text      string           `yaml:"text"`
	MediaIds  []string         `yaml:"media_ids"`
}

type FileStorageProfile struct {
	path        string
	ProfileName string                   `yaml:"profile_name"`
	Posts       []FileStorageProfilePost `yaml:"posts"`
}

func (fsp *FileStorageProfile) GetPost(postId string) *FileStorageProfilePost {
	if fsp.Posts == nil {
		return nil
	}
	for n, post := range fsp.Posts {
		if post.ID == postId {
			return &fsp.Posts[n]
		}
	}
	return nil
}

func (fsp *FileStorageProfile) Load() error {
	f, err := os.Open(fsp.path)
	if err != nil {
		return err
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(fsp); err != nil {
		ErrorLog.Println(fsp.path, err.Error())
		return err
	}
	return nil
}

func (fsp *FileStorageProfile) Save() error {
	if err := os.MkdirAll(path.Dir(fsp.path), 0755); err != nil {
		ErrorLog.Println(fsp.path, err.Error())
		return err
	}
	f, err := os.OpenFile(fsp.path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		ErrorLog.Println(fsp.path, err.Error())
		return err
	}
	defer f.Close()
	encoder := yaml.NewEncoder(f)
	if err := encoder.Encode(*fsp); err != nil {
		ErrorLog.Println(fsp.path, err.Error())
		return err
	}
	return nil
}

type FileStorage struct {
	StoragePath string
}

func (fs *FileStorage) GetProfile(profileName string) (*FileStorageProfile, error) {
	fsp := FileStorageProfile{
		path:        path.Join(fs.StoragePath, profileName, `.profile.yaml`),
		ProfileName: profileName,
	}
	return &fsp, fsp.Load()
}

func (fs *FileStorage) IsExist(profileName string, postId string) bool {
	profile, err := fs.GetProfile(profileName)
	if err != nil {
		return false
	}
	if post := profile.GetPost(postId); post != nil {
		return true
	}
	return false
}

func (fs *FileStorage) SaveMedia(profileName string, _ string, mediaId string, mediaType string, createdAt time.Time, data []byte) error {
	fileName := mediaId
	if mediaType == `PHOTO` {
		fileName = fileName + `.jpg`
	}
	if mediaType == `VIDEO` {
		fileName = fileName + `.mp4`
	}
	profilePath := path.Join(fs.StoragePath, profileName)
	if err := os.MkdirAll(profilePath, 0755); err != nil {
		ErrorLog.Println(err.Error())
		return err
	}
	filePath := path.Join(profilePath, fileName)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		ErrorLog.Println(filePath, err.Error())
		return err
	}
	defer os.Chtimes(filePath, time.Now(), createdAt)
	defer f.Close()
	_, err = f.Write(data)
	return err
}

func (fs *FileStorage) SavePost(profileName string, postId string, createdAt time.Time, location goinsta.Location, mediaIds []string, text string) error {
	profile, err := fs.GetProfile(profileName)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if post := profile.GetPost(postId); post != nil {
		post.CreatedAt = createdAt
		post.Location = location
		post.MediaIds = mediaIds
		post.Text = text
	} else {
		post := FileStorageProfilePost{
			ID:        postId,
			CreatedAt: createdAt,
			Location:  location,
			MediaIds:  mediaIds,
			Text:      text,
		}
		profile.Posts = append(profile.Posts, post)
	}
	return profile.Save()
}

func NewFileStorage(storagePath string) (*FileStorage, error) {
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		ErrorLog.Println(err.Error())
		return nil, err
	}
	fs := FileStorage{
		StoragePath: storagePath,
	}
	return &fs, nil
}
