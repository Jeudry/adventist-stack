// Package paginationtest holds the behaviour every paginated repository owes its callers.
//
// Three services implement the same pair of queries against three different tables, so the rules
// live here once instead of being re-proven per service: the count agrees with the rows, the pages
// cover everything exactly once, a soft-deleted row leaves both, an offset past the end is empty
// rather than an error, and an empty search does not filter.
//
// The contract is the List/Count pair; everything else in Subject is the fixture, which is the only
// part a service has to supply.
package paginationtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Jeudry/adventist-stack/pkg/pagination"
)

type Subject[T any] struct {
	List  func(ctx context.Context, q pagination.Query) ([]T, error)
	Count func(ctx context.Context, q pagination.Query) (int, error)

	// Create must store a row the given name finds through Search, and register its own cleanup
	// on the testing.T it receives, which is the subtest's and not the caller's.
	Create func(t *testing.T, name string) T
	Delete func(t *testing.T, item T)
	ID     func(T) uuid.UUID
}

const (
	rows     = 3
	pageSize = 2
)

// Run checks the whole contract. Each case seeds its own rows under a search term nobody else
// shares, because the tables outlive the tests and carry dev data.
func Run[T any](t *testing.T, s Subject[T]) {
	t.Helper()

	t.Run("counts every matching row", func(t *testing.T) {
		search, _ := seed(t, s, rows)

		count, err := s.Count(t.Context(), pagination.Query{Search: search})
		require.NoError(t, err, "count")
		assert.Equal(t, rows, count, "count")
	})

	t.Run("pages cover every row exactly once", func(t *testing.T) {
		search, created := seed(t, s, rows)

		page := func(offset int) []T {
			t.Helper()
			got, err := s.List(t.Context(), pagination.Query{Search: search, Limit: pageSize, Offset: offset})
			require.NoError(t, err, "page at offset %d", offset)
			return got
		}

		first, second := page(0), page(pageSize)
		assert.Len(t, first, pageSize, "a full page")
		assert.Len(t, second, rows-pageSize, "the last page carries the remainder")

		// Coverage without overlap is what pagination promises. Asserting an exact order would go
		// red on its own whenever two rows land on the same created_at.
		seen := map[uuid.UUID]bool{}
		for _, got := range [][]T{first, second} {
			for _, item := range got {
				id := s.ID(item)
				assert.True(t, created[id], "a page returned an unrelated row: %v", id)
				assert.False(t, seen[id], "a row appeared on both pages: %v", id)
				seen[id] = true
			}
		}
		assert.Len(t, seen, rows, "rows covered by both pages")
	})

	t.Run("a deleted row leaves the list and the count", func(t *testing.T) {
		search, _ := seed(t, s, 2)

		before, err := s.List(t.Context(), pagination.Query{Search: search, Limit: 10})
		require.NoError(t, err, "list before deleting")
		require.Len(t, before, 2, "list before deleting")

		s.Delete(t, before[0])

		after, err := s.List(t.Context(), pagination.Query{Search: search, Limit: 10})
		require.NoError(t, err, "list after deleting")
		require.Len(t, after, 1, "list after deleting")
		assert.Equal(t, s.ID(before[1]), s.ID(after[0]), "the surviving row")

		// Listing and counting are separate queries: fixing the filter in one and forgetting the
		// other shows up as a list of 1 next to a total of 2.
		count, err := s.Count(t.Context(), pagination.Query{Search: search})
		require.NoError(t, err, "count after deleting")
		assert.Equal(t, 1, count, "count after deleting")
	})

	t.Run("an offset past the end is empty and not an error", func(t *testing.T) {
		search, _ := seed(t, s, 1)

		got, err := s.List(t.Context(), pagination.Query{Search: search, Limit: pageSize, Offset: 500})
		require.NoError(t, err, "reading past the last page must not fail")
		assert.Empty(t, got, "a page past the end")
	})

	t.Run("an empty search does not filter", func(t *testing.T) {
		search, _ := seed(t, s, 1)

		scoped, err := s.Count(t.Context(), pagination.Query{Search: search})
		require.NoError(t, err, "scoped count")

		// The table is shared, so the exact total is unknowable. What this pins down is the branch
		// that matters: an empty search must widen the result, never collapse it to zero.
		all, err := s.Count(t.Context(), pagination.Query{})
		require.NoError(t, err, "unfiltered count")
		assert.GreaterOrEqual(t, all, scoped, "an empty search must match at least what a term matches")
		assert.NotZero(t, all, "an empty search must not filter everything out")
	})
}

func seed[T any](t *testing.T, s Subject[T], n int) (search string, created map[uuid.UUID]bool) {
	t.Helper()

	search = "pg" + uuid.NewString()
	created = make(map[uuid.UUID]bool, n)
	for i := range n {
		created[s.ID(s.Create(t, fmt.Sprintf("%s #%d", search, i)))] = true
	}
	return search, created
}
