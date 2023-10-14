package repo_test

// generate tests for AddPack

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/almiluk/sipacks/internal/adapter/repo"
	"github.com/almiluk/sipacks/internal/entity"
	"github.com/almiluk/sipacks/pkg/postgres"
	"github.com/jackc/pgx/v4"
)

const (
	pgURL   = "postgres://user:pass@localhost:5432/postgres"
	poolMax = 2
)

var pgRepo *repo.PostgresRepo

func TestMain(m *testing.M) {
	// setup
	pg, err := postgres.New(pgURL, postgres.MaxPoolSize(poolMax))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	pgRepo = repo.New(pg)
	// run tests
	os.Exit(m.Run())
	// teardown
}

func TestPostgresRepo_AddPack(t *testing.T) {
	ctx := context.Background()
	batch := pgx.Batch{}
	// Delete pack
	sql, args, err := pgRepo.Builder.Delete("pack").Where("name = ?", "name").ToSql()
	if err != nil {
		t.Fatal(err)
	}
	batch.Queue(sql, args...)

	// Delete author
	sql, args, err = pgRepo.Builder.Delete("author").Where("nickname = ?", "author").ToSql()
	if err != nil {
		t.Fatal(err)
	}
	batch.Queue(sql, args...)

	// Delete tags
	sql, args, err = pgRepo.Builder.Delete("tag").Where("name = ? OR name = ?", "tag1", "tag2").ToSql()
	if err != nil {
		t.Fatal(err)
	}
	batch.Queue(sql, args...)

	result := pgRepo.Pool.SendBatch(ctx, &batch)
	defer func() {
		err := result.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	tags := []entity.Tag{{Name: "tag1"}, {Name: "tag2"}}

	pack := entity.Pack{
		Name:         "name",
		Author:       entity.Author{Nickname: "author"},
		CreationDate: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		FileSize:     0,
		DownloadsNum: 0,
		Tags:         tags,
	}

	err = pgRepo.AddPack(ctx, &pack)
	if err != nil {
		t.Fatal(err)
	}

	// Check insertion

	log.Printf("pack: %+v\n", pack)
	// Check author
	sql, args, err = pgRepo.Builder.Select("id", "nickname").From("author").Where("nickname = ?", "author").ToSql()
	if err != nil {
		t.Fatal(err)
	}

	var author entity.Author
	err = pgRepo.Pool.QueryRow(ctx, sql, args...).Scan(&author.Id, &author.Nickname)
	if err != nil {
		t.Fatal(err)
	}

	if author != pack.Author {
		t.Errorf("author != pack.Author: %+v != %+v", author, pack.Author)
	}

	// Check tags
	sql, args, err = pgRepo.Builder.Select("id", "name").From("tag").Where("name = ? OR name = ?", "tag1", "tag2").OrderBy("id").ToSql()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := pgRepo.Pool.Query(ctx, sql, args...)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	readTags := []entity.Tag{}
	for rows.Next() {
		var tag entity.Tag
		err = rows.Scan(&tag.Id, &tag.Name)
		if err != nil {
			t.Fatal(err)
		}
		readTags = append(readTags, tag)
	}

	for i, tag := range readTags {
		if tag != pack.Tags[i] {
			t.Errorf("tag != pack.Tags[i]: %+v != %+v", tag, pack.Tags[i])
		}
	}

	// Check pack_tag
	sql, args, err = pgRepo.Builder.Select("pack_id", "tag_id").From("pack_tag").Where("pack_id = ?", pack.Id).OrderBy("tag_id").ToSql()
	if err != nil {
		t.Fatal(err)
	}

	rows, err = pgRepo.Pool.Query(ctx, sql, args...)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	readPackTag := [][]int{}
	for rows.Next() {
		var packId, tagId int
		err = rows.Scan(&packId, &tagId)
		if err != nil {
			t.Fatal(err)
		}
		readPackTag = append(readPackTag, []int{packId, tagId})
	}

	for i, packTag := range readPackTag {
		if packTag[0] != int(pack.Id) || packTag[1] != int(pack.Tags[i].Id) {
			t.Errorf("incorrect pack_tag: pack.Id==%d, pack.Tags[i].Id==%d, packTag[0]==%d, packTag[1]==%d", pack.Id, pack.Tags[i].Id, packTag[0], packTag[1])
		}
	}

	pack2 := entity.Pack{
		Name:         "name2",
		Author:       entity.Author{Nickname: "author"},
		CreationDate: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		FileSize:     0,
		DownloadsNum: 0,
		Tags:         tags,
	}

	err = pgRepo.AddPack(ctx, &pack2)
	if err != nil {
		t.Fatal(err)
	}

	if pack2.Author != pack.Author {
		t.Errorf("pack2.Author.Id != pack.Author.Id: %d != %d", pack2.Author.Id, pack.Author.Id)
	}

	for i, tag := range pack2.Tags {
		if tag.Id != pack.Tags[i].Id {
			t.Errorf("tag.Id != pack.Tags[i].Id: %d != %d", tag.Id, pack.Tags[i].Id)
		}
	}
}

func TestPostgresRepo_GetPacks(t *testing.T) {
	packs := []entity.Pack{
		{
			Name:         "name1",
			Author:       entity.Author{Nickname: "author"},
			CreationDate: time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
			FileSize:     0,
			DownloadsNum: 0,
			Tags:         []entity.Tag{{Name: "tag1"}, {Name: "tag2"}},
		},
		{
			Name:         "name2",
			Author:       entity.Author{Nickname: "author"},
			CreationDate: time.Date(1975, 1, 1, 0, 0, 0, 0, time.UTC),
			FileSize:     0,
			DownloadsNum: 0,
			Tags:         []entity.Tag{{Name: "tag1"}},
		},
		{
			Name:         "name3",
			Author:       entity.Author{Nickname: "author2"},
			CreationDate: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			FileSize:     0,
			DownloadsNum: 0,
			Tags:         []entity.Tag{{Name: "tag2"}},
		},
	}

	type testCase struct {
		name   string
		filter entity.PackFilter
		want   []entity.Pack
	}

	strPtr := func(s string) *string {
		return &s
	}

	// Clear pack db table
	ctx := context.Background()
	sql, args, err := pgRepo.Builder.Delete("pack").ToSql()
	if err != nil {
		t.Fatal(err)
	}
	_, err = pgRepo.Pool.Exec(ctx, sql, args...)
	if err != nil {
		t.Fatal(err)
	}

	// Add packs

	for i := range packs {
		err = pgRepo.AddPack(ctx, &packs[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	time1975 := time.Date(1975, 1, 1, 0, 0, 0, 0, time.UTC)
	time1979 := time.Date(1979, 1, 1, 0, 0, 0, 0, time.UTC)

	testCases := []testCase{
		{
			name:   "no filter",
			filter: entity.PackFilter{},
			want:   packs,
		},
		{
			name:   "filter by name",
			filter: entity.PackFilter{Name: strPtr("name")},
			want:   packs,
		},
		{
			name:   "filter by author",
			filter: entity.PackFilter{Author: strPtr("author")},
			want:   packs[:2],
		},
		{
			name:   "filter by tags (tag1)",
			filter: entity.PackFilter{Tags: []string{"tag1"}},
			want:   packs[:2],
		},
		{
			name:   "filter by tags (tag2)",
			filter: entity.PackFilter{Tags: []string{"tag2"}},
			want:   []entity.Pack{packs[0], packs[2]},
		},
		{
			name:   "filter by tags (tag1, tag2)",
			filter: entity.PackFilter{Tags: []string{"tag1", "tag2"}},
			want:   packs[0:1],
		},
		{
			name:   "filter by min_creation_date",
			filter: entity.PackFilter{MinCreationDate: &time1975},
			want:   packs[:2],
		},
		{
			name:   "filter by max_creation_date",
			filter: entity.PackFilter{MaxCreationDate: &time1975},
			want:   packs[1:3],
		},
		{
			name:   "filter by min_creation_date and max_creation_date",
			filter: entity.PackFilter{MinCreationDate: &time1975, MaxCreationDate: &time1979},
			want:   packs[1:2],
		},
	}

	for _, tc := range testCases {
		result, err := pgRepo.GetPacks(ctx, tc.filter)
		if err != nil {
			t.Fatal(fmt.Errorf("%s: %w", tc.name, err))
		}

		if len(result) != len(tc.want) {
			t.Errorf("%s: len(result) != len(tc.want): %d != %d", tc.name, len(result), len(tc.want))
			continue
		}

		for i, pack := range result {
			if !ComparePacks(pack, tc.want[i]) {
				t.Errorf("%s: pack != tc.want: %+v != %+v", tc.name, pack, tc.want[i])
			}
		}
	}
}

func ComparePacks(p1, p2 entity.Pack) bool {
	return p1.Id == p2.Id &&
		p1.Name == p2.Name &&
		p1.Author.Id == p2.Author.Id &&
		p1.Author.Nickname == p2.Author.Nickname &&
		p1.CreationDate == p2.CreationDate &&
		p1.FileSize == p2.FileSize &&
		p1.DownloadsNum == p2.DownloadsNum
}
