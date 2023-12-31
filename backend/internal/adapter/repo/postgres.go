package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/almiluk/sipacks/internal/entity"
	"github.com/almiluk/sipacks/pkg/postgres"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type PostgresRepo struct {
	*postgres.Postgres
}

func NewPGRepo(url string, opts ...postgres.Option) (*PostgresRepo, error) {
	pg, err := postgres.New(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("PostgresRepo - New - postgres.New: %w", err)
	}

	return &PostgresRepo{
		Postgres: pg,
	}, nil
}

func (pg *PostgresRepo) AddPack(ctx context.Context, pack *entity.Pack) error {
	// Begin transaction
	tx, err := pg.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("PostgresRepo - AddPack - pg.Pool.BeginTx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check if author exists and add if not
	author, err := pg.getAuthor(ctx, tx, pack.Author.Nickname)
	if errors.Is(err, pgx.ErrNoRows) {
		err = pg.addAuthor(ctx, tx, &pack.Author)
		if err != nil {
			return fmt.Errorf("PostgresRepo - AddPack - pg.addAuthor: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("PostgresRepo - AddPack - pg.getAuthor: %w", err)
	} else {
		pack.Author = author
	}

	// Check if all tags exist and add if not
	for i, tag := range pack.Tags {
		tag, err = pg.getTag(ctx, tx, tag.Name)
		if errors.Is(err, pgx.ErrNoRows) {
			err = pg.addTag(ctx, tx, &pack.Tags[i])
			if err != nil {
				return fmt.Errorf("PostgresRepo - AddPack - pg.addTag: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("PostgresRepo - AddPack - pg.getTag: %w", err)
		} else {
			pack.Tags[i] = tag
		}
	}

	// Add pack
	sql, args, err := pg.Builder.
		Insert("pack").
		Columns("name", "author_id", "creation_date", "file_size", "downloads_num", "guid").
		Values(pack.Name, pack.Author.ID, pack.CreationDate, pack.FileSize, 0, pack.GUID).
		Suffix("RETURNING \"id\"").
		ToSql()
	if err != nil {
		return fmt.Errorf("PostgresRepo - AddPack - Insert: %w", err)
	}

	err = tx.QueryRow(ctx, sql, args...).Scan(&pack.ID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return fmt.Errorf("PostgresRepo - AddPack - tx.QueryRow: %w", entity.ErrPackAlreadyExists)
	} else if err != nil {
		return fmt.Errorf("PostgresRepo - AddPack - tx.QueryRow: %w", err)
	}

	if err != nil {
		return fmt.Errorf("PostgresRepo - AddPack - tx.QueryRow: %w", err)
	}

	// Link tags to pack
	packTags := [][]interface{}{}
	for _, tag := range pack.Tags {
		packTags = append(packTags, []interface{}{pack.ID, tag.ID})
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"pack_tag"},
		[]string{"pack_id", "tag_id"},
		pgx.CopyFromRows(packTags),
	)
	if err != nil {
		return fmt.Errorf("PostgresRepo - AddPack - tx.CopyFrom: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("PostgresRepo - AddPack - tx.Commit: %w", err)
	}

	return nil
}

func (pg *PostgresRepo) GetPacks(ctx context.Context, filter entity.PackFilter) ([]entity.Pack, error) {
	tx, err := pg.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("PostgresRepo - GetPacks - pg.Pool.BeginTx: %w", err)
	}
	defer tx.Rollback(ctx)

	builder := pg.Builder.
		Select("pack.id", "pack.name", "author.id", "author.nickname", "pack.creation_date", "pack.file_size", "pack.downloads_num", "pack.guid").
		From("pack").
		Join("author ON pack.author_id = author.id")

	if filter.Name != nil {
		builder = builder.Where("LOWER(pack.name) LIKE '%' || ?|| '%'", filter.Name)
	}

	if filter.Author != nil {
		builder = builder.Where("author.nickname = ?", filter.Author)
	}

	if filter.MinCreationDate != nil {
		builder = builder.Where("pack.creation_date >= ?", filter.MinCreationDate)
	}

	if filter.MaxCreationDate != nil {
		builder = builder.Where("pack.creation_date <= ?", filter.MaxCreationDate)
	}

	if filter.Tags != nil && len(filter.Tags) > 0 {
		builder = builder.Join("pack_tag ON pack.id = pack_tag.pack_id").
			Join("tag ON pack_tag.tag_id = tag.id").
			Where("tag.name =ANY (?)", filter.Tags).
			GroupBy("pack.id", "author.id").
			Having("COUNT(pack_tag.tag_id) = ?", len(filter.Tags))
	}

	switch {
	case filter.SortBy == nil || *filter.SortBy == "":
		builder = builder.OrderBy("pack.id")
	case *filter.SortBy == "creation_date":
		builder = builder.OrderBy("pack.creation_date DESC")
	case *filter.SortBy == "downloads_num":
		builder = builder.OrderBy("pack.downloads_num DESC")
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("PostgresRepo - GetPacks - builder.ToSql: %w", err)
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PostgresRepo - GetPacks - pg.Pool.Query: %w", err)
	}
	defer rows.Close()

	pack := entity.Pack{}
	packs := []entity.Pack{}
	for rows.Next() {
		err = rows.Scan(
			&pack.ID, &pack.Name, &pack.Author.ID, &pack.Author.Nickname,
			&pack.CreationDate, &pack.FileSize, &pack.DownloadsNum, &pack.GUID,
		)

		if err != nil {
			return nil, fmt.Errorf("PostgresRepo - GetPacks - rows.Scan: %w", err)
		}

		packs = append(packs, pack)
	}

	for i := range packs {
		packs[i].Tags, err = pg.getPackTags(ctx, tx, packs[i].ID)
		if err != nil {
			return nil, fmt.Errorf("PostgresRepo - GetPacks - pg.getPackTags: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("PostgresRepo - GetPacks - tx.Commit: %w", err)
	}

	return packs, nil
}

func (pg *PostgresRepo) getAuthor(ctx context.Context, tx pgx.Tx, nickname string) (entity.Author, error) {
	sql, args, err := pg.Builder.Select("*").From("author").Where("nickname = ?", nickname).ToSql()
	if err != nil {
		return entity.Author{}, fmt.Errorf("PostgresRepo - getAuthor - Select: %w", err)
	}

	row := tx.QueryRow(ctx, sql, args...)

	var author entity.Author
	err = row.Scan(&author.ID, &author.Nickname)

	return author, err
}

func (pg *PostgresRepo) addAuthor(ctx context.Context, tx pgx.Tx, author *entity.Author) error {
	sql, args, err := pg.Builder.Insert("author").Columns("nickname").Values(author.Nickname).Suffix("RETURNING \"id\"").ToSql()
	if err != nil {
		return fmt.Errorf("PostgresRepo - addAuthor - Insert: %w", err)
	}

	err = tx.QueryRow(ctx, sql, args...).Scan(&author.ID)
	if err != nil {
		return fmt.Errorf("PostgresRepo - addAuthor - tx.QueryRow: %w", err)
	}

	return err
}

func (pg *PostgresRepo) getTag(ctx context.Context, tx pgx.Tx, name string) (entity.Tag, error) {
	sql, args, err := pg.Builder.Select("*").From("tag").Where("name = ?", name).ToSql()
	if err != nil {
		return entity.Tag{}, fmt.Errorf("PostgresRepo - getTag - Select: %w", err)
	}

	row := tx.QueryRow(ctx, sql, args...)

	var tag entity.Tag
	err = row.Scan(&tag.ID, &tag.Name)
	return tag, err
}

func (pg *PostgresRepo) addTag(ctx context.Context, tx pgx.Tx, tag *entity.Tag) error {
	sql, args, err := pg.Builder.Insert("tag").Columns("name").Values(tag.Name).Suffix("RETURNING \"id\"").ToSql()
	if err != nil {
		return fmt.Errorf("PostgresRepo - addTag - Insert: %w", err)
	}

	err = tx.QueryRow(ctx, sql, args...).Scan(&tag.ID)
	if err != nil {
		return fmt.Errorf("PostgresRepo - addTag - tx.QueryRow: %w", err)
	}

	return err
}

func (pg *PostgresRepo) getPackTags(ctx context.Context, tx pgx.Tx, id uint32) ([]entity.Tag, error) {
	sql, args, err := pg.Builder.Select("tag.id", "tag.name").
		From("pack_tag").
		Join("tag ON pack_tag.tag_id = tag.id").
		Where("pack_tag.pack_id = ?", id).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PostgresRepo - getPackTags - Select: %w", err)
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PostgresRepo - getPackTags - tx.Query: %w", err)
	}
	defer rows.Close()

	tag := entity.Tag{}
	tags := []entity.Tag{}
	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, fmt.Errorf("PostgresRepo - getPackTags - rows.Scan: %w", err)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (pg *PostgresRepo) IncreaseDownloadsCounter(ctx context.Context, guid string) error {
	sql, args, err := pg.Builder.Update("pack").
		Set("downloads_num", squirrel.Expr("downloads_num + 1")).
		Where("guid = ?", guid).
		ToSql()
	if err != nil {
		return fmt.Errorf("PostgresRepo - IncreaseDownloadsCounter - ToSql: %w", err)
	}

	_, err = pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PostgresRepo - IncreaseDownloadsCounter - pg.Pool.Exec: %w", err)
	}

	return nil
}
