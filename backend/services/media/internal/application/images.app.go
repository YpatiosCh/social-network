package application

import (
	"context"
	"fmt"
	"net/url"
	"social-network/services/media/internal/db/dbservice"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
	"time"

	"github.com/google/uuid"
)

// Provides a fileId and an upload url targeted on bucket Originals defined on configs.
// Creates all variant entries provided in []variants for workers to later
// create asynchronously the compressed files.
func (m *MediaService) UploadImage(ctx context.Context,
	fm md.FileMeta, exp time.Duration, variants []ct.FileVariant) (fileId ct.Id, upUrl string, err error) {
	// if err := ct.ValidateStruct(fm); err != nil {
	// 	return fileId, "", err
	// }

	var url *url.URL
	errTx := m.txRunner.RunTx(ctx,
		func(q dbservice.Querier) error {
			fm.ObjectKey = uuid.NewString()
			fileId, err = m.Queries.CreateFile(ctx, dbservice.File{
				Filename:   fm.Filename,
				MimeType:   fm.MimeType,
				SizeBytes:  fm.SizeBytes,
				Visibility: fm.Visibility,
				Bucket:     m.Cfgs.FileService.Buckets.Originals,
				ObjectKey:  fm.ObjectKey,
				Status:     ct.Complete,
				Variant:    ct.Original,
			})

			if err != nil {
				return err
			}

			for _, v := range variants {
				_, err := m.Queries.CreateVariant(ctx, dbservice.File{
					Filename:   fm.Filename,
					MimeType:   fm.MimeType,
					SizeBytes:  fm.SizeBytes,
					Bucket:     m.Cfgs.FileService.Buckets.Variants,
					ObjectKey:  fm.ObjectKey + "/" + v.String(),
					Visibility: fm.Visibility,
					Status:     ct.Pending,
					Variant:    v,
				})
				if err != nil {
					return fmt.Errorf(
						"internal database error: %v failed to create variant %v for file with id: %v",
						err, v, fileId)
				}
			}

			url, err = m.Clients.GenerateUploadURL(ctx, fm.Bucket, fm.ObjectKey, exp)
			if err != nil {
				return err
			}
			return nil
		},
	)

	if errTx != nil {
		return 0, "", errTx
	}
	return fileId, url.String(), nil
}

// Returns an image download URL for the requested imageId and Variant.
// If the variant is not available it falls back to the original file.
func (m *MediaService) GetImage(ctx context.Context,
	imgId ct.Id, variant ct.FileVariant,
) (downUrl string, err error) {
	if !imgId.IsValid() || !variant.IsValid() {
		return "", ct.ErrValidation
	}

	var fm dbservice.File
	var url *url.URL

	errTx := m.txRunner.RunTx(ctx,
		func(q dbservice.Querier) error {
			switch variant {
			case ct.Original:
				fm, err = m.Queries.GetFileById(ctx, imgId)
			default:
				fm, err = m.Queries.GetVariant(ctx, imgId, variant)
				if fm.Status != ct.Complete {
					fm, err = m.Queries.GetFileById(ctx, imgId)
					if fm.Status != ct.Complete {
						return fmt.Errorf("file validation status is %v", fm.Status)
					}
				}
			}
			if err != nil {
				return err
			}

			url, err = m.Clients.GenerateDownloadURL(ctx, fm.Bucket, fm.ObjectKey, fm.Visibility.SetExp())
			if err != nil {
				return err
			}
			return nil
		},
	)

	if errTx != nil {
		return "", err
	}
	return url.String(), err
}

// Not allowing originals in batch request
func (m *MediaService) GetImages(ctx context.Context,
	imgIds ct.Ids, variant ct.FileVariant,
) (downUrls map[ct.Id]string, err error) {
	if !imgIds.IsValid() || !variant.IsValid() || variant == ct.Original {
		return nil, ct.ErrValidation
	}
	var na ct.Ids
	var fms []dbservice.File

	errTx := m.txRunner.RunTx(ctx,
		func(q dbservice.Querier) error {

			fms, na, err = m.Queries.GetVariants(ctx, uniqueIds(imgIds), variant)
			if err != nil {
				return err
			}
			if len(na) != 0 {
				originals, err := m.Queries.GetFiles(ctx, na)
				if err != nil {
					return err
				}
				fms = append(fms, originals...)
			}
			return nil
		},
	)

	if errTx != nil {
		return nil, err
	}

	downUrls = make(map[ct.Id]string, len(fms))
	for _, fm := range fms {
		if fm.Status != ct.Complete {
			fmt.Printf("requested file %v validation status is %v", fm.Id, fm.Status)
			continue
		}
		url, err := m.Clients.GenerateDownloadURL(ctx, fm.Bucket, fm.ObjectKey, fm.Visibility.SetExp())
		if err != nil {
			return nil, err
		}
		downUrls[fm.Id] = url.String()
	}

	return downUrls, nil
}

func (m *MediaService) ValidateUpload(ctx context.Context,
	upload md.FileMeta) error {
	if err := ct.ValidateStruct(upload); err != nil {
		return err
	}

	if err := m.Clients.ValidateUpload(ctx, upload); err != nil {
		m.Queries.UpdateFileStatus(ctx, upload.Id, ct.Failed)
		return err
	}

	m.Queries.UpdateFileStatus(ctx, upload.Id, ct.Complete)

	return nil
}

func uniqueIds(ids ct.Ids) ct.Ids {
	uniq := make(map[ct.Id]struct{}, len(ids))
	cleaned := make([]ct.Id, 0, len(ids))
	for _, id := range ids {
		if _, ok := uniq[id]; !ok {
			uniq[id] = struct{}{}
			cleaned = append(cleaned, id)
		}
	}
	return ct.Ids(cleaned)
}
