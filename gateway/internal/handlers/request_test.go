package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseListQuery(t *testing.T) {
	cases := []struct {
		name     string
		url      string
		expected listQuery
	}{
		{
			name:     "reads every parameter",
			url:      "/members?page=2&page_size=50&search=ana",
			expected: listQuery{Page: 2, PageSize: 50, Search: "ana"},
		},
		{
			name:     "missing parameters stay zero so the service applies its defaults",
			url:      "/members",
			expected: listQuery{Page: 0, PageSize: 0, Search: ""},
		},
		{
			name:     "non numeric values stay zero instead of failing the request",
			url:      "/members?page=abc&page_size=",
			expected: listQuery{Page: 0, PageSize: 0, Search: ""},
		},
		{
			name:     "negative values are forwarded untouched",
			url:      "/members?page=-3&page_size=-1",
			expected: listQuery{Page: -3, PageSize: -1, Search: ""},
		},
		{
			name:     "search keeps its spacing for the service to trim",
			url:      "/members?search=%20ana%20",
			expected: listQuery{Page: 0, PageSize: 0, Search: " ana "},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, testCase.url, nil)

			actual := parseListQuery(request)

			if actual != testCase.expected {
				t.Errorf("expected %+v, got %+v", testCase.expected, actual)
			}
		})
	}
}
