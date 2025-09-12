package storage_database_postgres

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	storage_core "github.com/gabrielmrtt/taski/internal/storage"
	storage_repositories "github.com/gabrielmrtt/taski/internal/storage/repositories"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	user_database_postgres "github.com/gabrielmrtt/taski/internal/user/database/postgres"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UploadedFileTable struct {
	bun.BaseModel `bun:"table:uploaded_file,alias:uploaded_file"`

	InternalId               string `bun:"internal_id,pk,notnull,type:uuid"`
	PublicId                 string `bun:"public_id,notnull,type:varchar(510)"`
	File                     string `bun:"file,notnull,type:text"`
	FileDirectory            string `bun:"file_directory,notnull,type:text"`
	FileMimeType             string `bun:"file_mime_type,notnull,type:varchar(100)"`
	FileExtension            string `bun:"file_extension,notnull,type:varchar(3)"`
	UserUploadedByInternalId string `bun:"user_uploaded_by_internal_id,notnull,type:uuid"`
	UploadedAt               int64  `bun:"uploaded_at,notnull,type:bigint"`

	UserUploadedBy *user_database_postgres.UserTable `bun:"rel:has-one,join:user_uploaded_by_internal_id=internal_id"`
}

func (u *UploadedFileTable) ToEntity() *storage_core.UploadedFile {
	return &storage_core.UploadedFile{
		Identity:               core.NewIdentityFromInternal(uuid.MustParse(u.InternalId), storage_core.UploadedFileIdentityPrefix),
		File:                   &u.File,
		FileDirectory:          &u.FileDirectory,
		FileMimeType:           &u.FileMimeType,
		FileExtension:          &u.FileExtension,
		UserUploadedByIdentity: core.NewIdentityFromInternal(uuid.MustParse(u.UserUploadedByInternalId), user_core.UserIdentityPrefix),
		UploadedAt:             u.UploadedAt,
	}
}

type UploadedFilePostgresRepository struct {
	db *bun.DB
	tx *core_database_postgres.TransactionPostgres
}

func NewUploadedFilePostgresRepository() *UploadedFilePostgresRepository {
	return &UploadedFilePostgresRepository{db: core_database_postgres.DB, tx: nil}
}

func (r *UploadedFilePostgresRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*core_database_postgres.TransactionPostgres)
	return nil
}

func (r *UploadedFilePostgresRepository) GetUploadedFileByIdentity(params storage_repositories.GetUploadedFileByIdentityParams) (*storage_core.UploadedFile, error) {
	var uploadedFile *UploadedFileTable = new(UploadedFileTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(uploadedFile).Where("internal_id = ?", params.FileIdentity.Internal.String())
	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if uploadedFile.InternalId == "" {
		return nil, nil
	}

	return uploadedFile.ToEntity(), nil
}

func (r *UploadedFilePostgresRepository) StoreUploadedFile(params storage_repositories.StoreUploadedFileParams) (*storage_core.UploadedFile, error) {
	var tx bun.Tx
	var shouldCommit bool = false

	if r.tx != nil && !r.tx.IsClosed() {
		tx = *r.tx.Tx
	} else {
		var err error
		tx, err = r.db.BeginTx(context.Background(), nil)
		shouldCommit = true

		if err != nil {
			return nil, err
		}
	}

	uploadedFileTable := &UploadedFileTable{
		InternalId:               params.UploadedFile.Identity.Internal.String(),
		PublicId:                 params.UploadedFile.Identity.Public,
		File:                     *params.UploadedFile.File,
		FileDirectory:            *params.UploadedFile.FileDirectory,
		FileMimeType:             *params.UploadedFile.FileMimeType,
		FileExtension:            *params.UploadedFile.FileExtension,
		UserUploadedByInternalId: params.UploadedFile.UserUploadedByIdentity.Internal.String(),
		UploadedAt:               params.UploadedFile.UploadedAt,
	}

	_, err := tx.NewInsert().Model(uploadedFileTable).Exec(context.Background())

	if err != nil {
		return nil, err
	}

	if shouldCommit {
		err = tx.Commit()
	}

	return params.UploadedFile, nil
}

func (r *UploadedFilePostgresRepository) DeleteUploadedFile(params storage_repositories.DeleteUploadedFileParams) error {
	var tx bun.Tx
	var shouldCommit bool = false

	if r.tx != nil && !r.tx.IsClosed() {
		tx = *r.tx.Tx
	} else {
		var err error
		tx, err = r.db.BeginTx(context.Background(), nil)
		shouldCommit = true

		if err != nil {
			return err
		}
	}

	_, err := tx.NewDelete().Model(&UploadedFileTable{}).Where("internal_id = ?", params.FileIdentity.Internal.String()).Exec(context.Background())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return err
	}

	if shouldCommit {
		err = tx.Commit()
	}

	return nil
}
