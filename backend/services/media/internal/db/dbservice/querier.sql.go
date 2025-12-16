package dbservice

import (
	"context"
	ct "social-network/shared/go/customtypes"
)

type Querier interface {
	CreateFile(ctx context.Context, fm File) (fileId ct.Id, err error)

	GetFileById(ctx context.Context, fileId ct.Id) (fm File, err error)

	GetFiles(
		ctx context.Context, ids ct.Ids,
	) ([]File, error)

	UpdateFileStatus(
		ctx context.Context,
		fileId ct.Id,
		status ct.UploadStatus,
	) error

	CreateVariant(ctx context.Context, fm File) (fileId ct.Id, err error)

	GetVariant(ctx context.Context, fileId ct.Id,
		variant ct.FileVariant) (fm File, err error)

	GetVariants(
		ctx context.Context,
		fileIds ct.Ids,
		variant ct.FileVariant,
	) (fms []File, notComplete []ct.Id, err error)

	UpdateVariantStatus(ctx context.Context,
		fileId ct.Id,
		variant ct.FileVariant,
		status ct.UploadStatus,
	) error
}

var _ Querier = (*Queries)(nil)
