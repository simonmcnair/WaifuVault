package thumbnail

import (
	"encoding/base64"
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/rs/zerolog/log"
	"github.com/waifuvault/WaifuVault/thumbnails/pkg/dao"
	"github.com/waifuvault/WaifuVault/thumbnails/pkg/mod"
	"github.com/waifuvault/WaifuVault/thumbnails/pkg/utils"
	"os"
)

type Service interface {
	GenerateThumbnails(files []mod.FileEntry) error
}

type service struct {
	dao dao.Dao
}

func NewService(daoService dao.Dao) Service {
	return &service{
		daoService,
	}
}

func (s *service) GenerateThumbnails(files []mod.FileEntry) error {
	var thumbnailStructs []mod.Thumbnail
	for _, file := range files {
		var thumbnailBytes []byte
		var err error
		if utils.IsImage(file.MediaType) {
			thumbnailBytes, err = generateImageThumbnail(file.FullFileNameOnSystem)
			if err != nil {
				log.Err(err).Msg("Error generating image thumbnail")
				continue
			}
		} else if utils.IsVideo(file.MediaType) {
			thumbnailBytes, err = generateVideoThumbnail(file.FullFileNameOnSystem)
			if err != nil {
				log.Err(err).Msg("Error generating video thumbnail")
				continue
			}
		}

		thumbnailStructs = append(thumbnailStructs, mod.Thumbnail{
			Id:     file.Id,
			Data:   base64.StdEncoding.EncodeToString(thumbnailBytes),
			FileId: 0,
		})
	}
	_, err := s.dao.SaveThumbnails(thumbnailStructs)
	if err != nil {
		return err
	}
	return nil
}

func generateVideoThumbnail(system string) ([]byte, error) {
	return nil, nil
}

func generateImageThumbnail(fileName string) ([]byte, error) {
	vips.Startup(nil)
	defer vips.Shutdown()

	data, err := os.ReadFile(utils.BaseUrl + "/" + fileName)
	if err != nil {
		return nil, err
	}

	intSet := vips.IntParameter{}
	intSet.Set(-1)
	params := vips.NewImportParams()
	params.NumPages = intSet

	image, err := vips.LoadImageFromBuffer(data, params)
	if err != nil {
		return nil, err
	}

	err = applyScale(400.00, image)
	if err != nil {
		return nil, err
	}

	err = image.RemoveMetadata("delay", "dispose", "loop", "loop_count")
	if err != nil {
		return nil, err
	}

	bytes, _, err := image.ExportNative()
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func applyScale(width float64, img *vips.ImageRef) error {
	scale := width / float64(img.Width())
	return img.Resize(scale, vips.KernelAuto)
}

/*package thumbnail

import (
"encoding/base64"
"github.com/davidbyttow/govips/v2/vips"
"github.com/rs/zerolog/log"
"github.com/waifuvault/WaifuVault/thumbnails/pkg/dao"
"github.com/waifuvault/WaifuVault/thumbnails/pkg/mod"
"github.com/waifuvault/WaifuVault/thumbnails/pkg/utils"
"os"
"sync"
)

type Service interface {
	GenerateThumbnails(files []mod.FileEntry) error
}

type service struct {
	dao dao.Dao
}

func NewService(daoService dao.Dao) Service {
	return &service{
		daoService,
	}
}

func (s *service) GenerateThumbnails(files []mod.FileEntry) error {
	numWorkers := 4
	filesChan := make(chan mod.FileEntry)
	resultsChan := make(chan mod.Thumbnail)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range filesChan {
				var thumbnailBytes []byte
				var err error
				if utils.IsImage(file.MediaType) {
					thumbnailBytes, err = generateImageThumbnail(file.FullFileNameOnSystem)
					if err != nil {
						log.Err(err).Msg("Error generating image thumbnail")
						continue
					}
				} else if utils.IsVideo(file.MediaType) {
					thumbnailBytes, err = generateVideoThumbnail(file.FullFileNameOnSystem)
					if err != nil {
						log.Err(err).Msg("Error generating video thumbnail")
						continue
					}
				} else {
					continue
				}
				resultsChan <- mod.Thumbnail{
					Id:     file.Id,
					Data:   base64.StdEncoding.EncodeToString(thumbnailBytes),
					FileId: 0,
				}
			}
		}()
	}

	go func() {
		for _, f := range files {
			filesChan <- f
		}
		close(filesChan)
	}()

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var thumbnailStructs []mod.Thumbnail
	for thumbnail := range resultsChan {
		thumbnailStructs = append(thumbnailStructs, thumbnail)
	}

	_, err := s.dao.SaveThumbnails(thumbnailStructs)
	if err != nil {
		return err
	}

	return nil
}

func generateVideoThumbnail(system string) ([]byte, error) {
	return nil, nil
}

func generateImageThumbnail(fileName string) ([]byte, error) {
	vips.Startup(nil)
	defer vips.Shutdown()

	data, err := os.ReadFile(utils.BaseUrl + "/" + fileName)
	if err != nil {
		return nil, err
	}

	intSet := vips.IntParameter{}
	intSet.Set(-1)
	params := vips.NewImportParams()
	params.NumPages = intSet

	image, err := vips.LoadImageFromBuffer(data, params)
	if err != nil {
		return nil, err
	}

	err = applyScale(400.00, image)
	if err != nil {
		return nil, err
	}

	err = image.RemoveMetadata("delay", "dispose", "loop", "loop_count")
	if err != nil {
		return nil, err
	}

	bytes, _, err := image.ExportNative()
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func applyScale(width float64, img *vips.ImageRef) error {
	scale := width / float64(img.Width())
	return img.Resize(scale, vips.KernelAuto)
}

*/
