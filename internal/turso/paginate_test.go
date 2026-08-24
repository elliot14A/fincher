package turso_test

import (
	"context"
	"errors"
	"testing"

	"github.com/elliot14A/fincher/internal/turso"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestOrderBy(t *testing.T) {
	pAsc := models.NewPagination(1, 10, "asc", "")
	if got := turso.OrderBy(pAsc, "ASC_ORDER", "DESC_ORDER"); got != "ASC_ORDER" {
		t.Errorf("expected ASC_ORDER, got %s", got)
	}

	pDesc := models.NewPagination(1, 10, "desc", "")
	if got := turso.OrderBy(pDesc, "ASC_ORDER", "DESC_ORDER"); got != "DESC_ORDER" {
		t.Errorf("expected DESC_ORDER, got %s", got)
	}
}

func TestPaginate_Success(t *testing.T) {
	ctx := context.Background()
	p := models.NewPagination(2, 5, "desc", "")

	countFn := func(ctx context.Context) (int, error) {
		return 20, nil
	}

	fetchFn := func(ctx context.Context, limit, offset int) ([]int, error) {
		if limit != 5 || offset != 5 {
			t.Errorf("unexpected limit/offset: %d, %d", limit, offset)
		}
		return []int{6, 7, 8, 9, 10}, nil
	}

	mapFn := func(items []int) []int {
		return items
	}

	res := turso.Paginate(ctx, "test.Paginate", p, countFn, fetchFn, mapFn)
	if res.IsErr() {
		t.Fatalf("expected Ok, got: %v", res.Error())
	}

	paginationRes := res.Unwrap()
	if len(paginationRes.Items) != 5 || paginationRes.TotalItems != 20 {
		t.Errorf("unexpected pagination result: %+v", paginationRes)
	}
	if paginationRes.Page != 2 || paginationRes.TotalPages != 4 {
		t.Errorf("expected page 2 of 4, got page %d of %d", paginationRes.Page, paginationRes.TotalPages)
	}
}

func TestPaginate_CountError(t *testing.T) {
	ctx := context.Background()
	p := models.NewPagination(1, 10, "desc", "")

	countFn := func(ctx context.Context) (int, error) {
		return 0, errors.New("db count failed")
	}
	fetchFn := func(ctx context.Context, limit, offset int) ([]int, error) {
		return nil, nil
	}

	res := turso.Paginate(ctx, "test.Paginate", p, countFn, fetchFn, func(i []int) []int { return i })
	if res.IsOk() {
		t.Fatalf("expected error from count failure")
	}

	domErr, ok := res.Error().(*domainerrors.DomainError)
	if !ok || domErr.Code != domainerrors.CodeInternal {
		t.Errorf("expected CodeInternal for count failure, got: %v", res.Error())
	}
}

func TestPaginate_FetchError(t *testing.T) {
	ctx := context.Background()
	p := models.NewPagination(1, 10, "desc", "")

	countFn := func(ctx context.Context) (int, error) {
		return 10, nil
	}
	fetchFn := func(ctx context.Context, limit, offset int) ([]int, error) {
		return nil, errors.New("db query failed")
	}

	res := turso.Paginate(ctx, "test.Paginate", p, countFn, fetchFn, func(i []int) []int { return i })
	if res.IsOk() {
		t.Fatalf("expected error from fetch failure")
	}

	domErr, ok := res.Error().(*domainerrors.DomainError)
	if !ok || domErr.Code != domainerrors.CodeInternal {
		t.Errorf("expected CodeInternal for fetch failure, got: %v", res.Error())
	}
}
