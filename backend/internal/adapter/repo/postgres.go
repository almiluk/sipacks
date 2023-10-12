package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/almiluk/sipacks/internal/entity"
	"github.com/almiluk/sipacks/pkg/postgres"
	"github.com/jackc/pgx/v4"
)

type PostgresRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *PostgresRepo {
	return &PostgresRepo{pg}
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
		Columns("name", "author_id", "creation_date", "file_size", "downloads_num").
		Values(pack.Name, pack.Author.Id, pack.CreationDate, pack.FileSize, 0).
		Suffix("RETURNING \"id\"").
		ToSql()
	if err != nil {
		return fmt.Errorf("PostgresRepo - AddPack - Insert: %w", err)
	}

	err = tx.QueryRow(ctx, sql, args...).Scan(&pack.Id)
	if err != nil {
		return fmt.Errorf("PostgresRepo - AddPack - tx.QueryRow: %w", err)
	}

	// Link tags to pack
	packTags := [][]interface{}{}
	for _, tag := range pack.Tags {
		packTags = append(packTags, []interface{}{pack.Id, tag.Id})
	}

	insertedTagsCnt, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"pack_tag"},
		[]string{"pack_id", "tag_id"},
		pgx.CopyFromRows(packTags),
	)
	if err != nil {
		return fmt.Errorf("PostgresRepo - AddPack - tx.CopyFrom: %w", err)
	}
	if insertedTagsCnt != int64(len(pack.Tags)) {
		return fmt.Errorf("PostgresRepo - AddPack - tx.CopyFrom: insertedTagsCnt != len(pack.Tags)")
	}

	tx.Commit(ctx)

	return nil
}

func (pg *PostgresRepo) GetPacks(ctx context.Context, filter entity.PackFilter) ([]entity.Pack, error) {
	sql, args, err := pg.Builder.
		Select("pack.id", "pack.name", "author.nickname", "pack.creation_date", "pack.file_size", "pack.downloads_num").
		From("pack").
		Join("author ON pack.author_id = author.id").
		Where("pack.name LIKE ?", "%"+filter.Name+"%").
		OrderBy("pack.id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PostgresRepo - GetPacks - Select: %w", err)
	}

	rows, err := pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PostgresRepo - GetPacks - pg.Pool.Query: %w", err)
	}
	defer rows.Close()

	packs := []entity.Pack{}
	for rows.Next() {
		var pack entity.Pack
		err = rows.Scan(&pack.Id, &pack.Name, &pack.Author.Nickname, &pack.CreationDate, &pack.FileSize, &pack.DownloadsNum)
		if err != nil {
			return nil, fmt.Errorf("PostgresRepo - GetPacks - rows.Scan: %w", err)
		}
		packs = append(packs, pack)
	}

	return packs, nil
}

func (pg *PostgresRepo) getAuthor(ctx context.Context, tx pgx.Tx, nickname string) (entity.Author, error) {
	sql, args, err := pg.Builder.Select("*").From("author").Where("nickname = ?", nickname).ToSql()
	if err != nil {
		return entity.Author{}, fmt.Errorf("PostgresRepo - getAuthor - Select: %w", err)
	}

	var author entity.Author
	row := tx.QueryRow(ctx, sql, args...)
	err = row.Scan(&author.Id, &author.Nickname)
	return author, err
}

func (pg *PostgresRepo) addAuthor(ctx context.Context, tx pgx.Tx, author *entity.Author) error {
	sql, args, err := pg.Builder.Insert("author").Columns("nickname").Values(author.Nickname).Suffix("RETURNING \"id\"").ToSql()
	if err != nil {
		return fmt.Errorf("PostgresRepo - addAuthor - Insert: %w", err)
	}

	err = tx.QueryRow(ctx, sql, args...).Scan(&author.Id)
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

	var tag entity.Tag
	row := tx.QueryRow(ctx, sql, args...)
	err = row.Scan(&tag.Id, &tag.Name)
	return tag, err
}

func (pg *PostgresRepo) addTag(ctx context.Context, tx pgx.Tx, tag *entity.Tag) error {
	sql, args, err := pg.Builder.Insert("tag").Columns("name").Values(tag.Name).Suffix("RETURNING \"id\"").ToSql()
	if err != nil {
		return fmt.Errorf("PostgresRepo - addTag - Insert: %w", err)
	}

	err = tx.QueryRow(ctx, sql, args...).Scan(&tag.Id)
	if err != nil {
		return fmt.Errorf("PostgresRepo - addTag - tx.QueryRow: %w", err)
	}
	return err
}
