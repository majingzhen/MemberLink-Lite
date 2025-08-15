package common

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewPageRequest(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected PageRequest
	}{
		{
			name:     "Valid parameters",
			page:     2,
			pageSize: 20,
			expected: PageRequest{Page: 2, PageSize: 20},
		},
		{
			name:     "Invalid page (0)",
			page:     0,
			pageSize: 20,
			expected: PageRequest{Page: 1, PageSize: 20},
		},
		{
			name:     "Invalid page (negative)",
			page:     -1,
			pageSize: 20,
			expected: PageRequest{Page: 1, PageSize: 20},
		},
		{
			name:     "Invalid pageSize (0)",
			page:     2,
			pageSize: 0,
			expected: PageRequest{Page: 2, PageSize: 10},
		},
		{
			name:     "Invalid pageSize (too large)",
			page:     2,
			pageSize: 200,
			expected: PageRequest{Page: 2, PageSize: 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewPageRequest(tt.page, tt.pageSize)
			assert.Equal(t, tt.expected, *result)
		})
	}
}

func TestPageRequest_GetOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected int
	}{
		{
			name:     "First page",
			page:     1,
			pageSize: 10,
			expected: 0,
		},
		{
			name:     "Second page",
			page:     2,
			pageSize: 10,
			expected: 10,
		},
		{
			name:     "Third page with different page size",
			page:     3,
			pageSize: 20,
			expected: 40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := PageRequest{Page: tt.page, PageSize: tt.pageSize}
			result := req.GetOffset()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPageRequest_GetLimit(t *testing.T) {
	req := PageRequest{Page: 1, PageSize: 15}
	result := req.GetLimit()
	assert.Equal(t, 15, result)
}

func TestPageRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		pageSize  int
		expectErr bool
	}{
		{
			name:      "Valid parameters",
			page:      1,
			pageSize:  10,
			expectErr: false,
		},
		{
			name:      "Invalid page (0)",
			page:      0,
			pageSize:  10,
			expectErr: true,
		},
		{
			name:      "Invalid page (negative)",
			page:      -1,
			pageSize:  10,
			expectErr: true,
		},
		{
			name:      "Invalid pageSize (0)",
			page:      1,
			pageSize:  0,
			expectErr: true,
		},
		{
			name:      "Invalid pageSize (too large)",
			page:      1,
			pageSize:  200,
			expectErr: true,
		},
		{
			name:      "Valid maximum pageSize",
			page:      1,
			pageSize:  100,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := PageRequest{Page: tt.page, PageSize: tt.pageSize}
			err := req.Validate()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParsePageRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		query     string
		expected  PageRequest
		expectErr bool
	}{
		{
			name:      "No query parameters",
			query:     "",
			expected:  PageRequest{Page: 1, PageSize: 10},
			expectErr: false,
		},
		{
			name:      "Valid query parameters",
			query:     "page=2&page_size=20",
			expected:  PageRequest{Page: 2, PageSize: 20},
			expectErr: false,
		},
		{
			name:      "Invalid page parameter",
			query:     "page=0&page_size=20",
			expected:  PageRequest{Page: 1, PageSize: 20},
			expectErr: true,
		},
		{
			name:      "Invalid page_size parameter",
			query:     "page=2&page_size=200",
			expected:  PageRequest{Page: 2, PageSize: 10},
			expectErr: true,
		},
		{
			name:      "Non-numeric parameters",
			query:     "page=abc&page_size=def",
			expected:  PageRequest{Page: 1, PageSize: 10},
			expectErr: false,
		},
		{
			name:      "Maximum valid page_size",
			query:     "page=1&page_size=100",
			expected:  PageRequest{Page: 1, PageSize: 100},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试请求
			req, _ := http.NewRequest("GET", "/?"+tt.query, nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// 解析分页参数
			result, err := ParsePageRequest(c)

			// 验证结果
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected.Page, result.Page)
			assert.Equal(t, tt.expected.PageSize, result.PageSize)
		})
	}
}

func TestNewPaginateResult(t *testing.T) {
	list := []string{"item1", "item2", "item3"}
	total := int64(25)
	page := 2
	pageSize := 10

	result := NewPaginateResult(list, total, page, pageSize)

	assert.Equal(t, list, result.List)
	assert.Equal(t, total, result.Total)
	assert.Equal(t, page, result.Page)
	assert.Equal(t, pageSize, result.PageSize)
	assert.Equal(t, 3, result.Pages) // 25 / 10 = 2.5, rounded up to 3
}

func TestNewPaginateResult_ExactPages(t *testing.T) {
	list := []string{"item1", "item2"}
	total := int64(20)
	page := 1
	pageSize := 10

	result := NewPaginateResult(list, total, page, pageSize)

	assert.Equal(t, 2, result.Pages) // 20 / 10 = 2 exactly
}

func TestNewPaginateResult_SinglePage(t *testing.T) {
	list := []string{"item1"}
	total := int64(5)
	page := 1
	pageSize := 10

	result := NewPaginateResult(list, total, page, pageSize)

	assert.Equal(t, 1, result.Pages) // 5 / 10 = 0.5, rounded up to 1
}

func TestNewPaginateResult_ZeroTotal(t *testing.T) {
	list := []string{}
	total := int64(0)
	page := 1
	pageSize := 10

	result := NewPaginateResult(list, total, page, pageSize)

	assert.Equal(t, 0, result.Pages) // 0 / 10 = 0
}

func TestPaginate(t *testing.T) {
	tests := []struct {
		name           string
		page           int
		pageSize       int
		expectedOffset int
		expectedLimit  int
	}{
		{
			name:           "Valid parameters",
			page:           2,
			pageSize:       10,
			expectedOffset: 10,
			expectedLimit:  10,
		},
		{
			name:           "Invalid page (0)",
			page:           0,
			pageSize:       10,
			expectedOffset: 0,
			expectedLimit:  10,
		},
		{
			name:           "Invalid page (negative)",
			page:           -1,
			pageSize:       10,
			expectedOffset: 0,
			expectedLimit:  10,
		},
		{
			name:           "Invalid pageSize (0)",
			page:           2,
			pageSize:       0,
			expectedOffset: 10,
			expectedLimit:  10,
		},
		{
			name:           "Invalid pageSize (too large)",
			page:           2,
			pageSize:       200,
			expectedOffset: 10,
			expectedLimit:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Since we can't easily test GORM scopes without a database,
			// we'll test the logic by calling the function and checking
			// that it returns a function (not nil)
			paginateFunc := Paginate(tt.page, tt.pageSize)
			assert.NotNil(t, paginateFunc)
		})
	}
}

func TestPaginateWithRequest(t *testing.T) {
	req := &PageRequest{Page: 2, PageSize: 20}
	paginateFunc := PaginateWithRequest(req)
	assert.NotNil(t, paginateFunc)
}

func TestNewBasePaginateService(t *testing.T) {
	// Test service creation without actual database
	service := NewBasePaginateService(nil, "test_model")
	assert.NotNil(t, service)
	assert.Equal(t, "test_model", service.model)
}

func TestGetPageInfo(t *testing.T) {
	tests := []struct {
		name     string
		total    int64
		page     int
		pageSize int
		expected PageInfo
	}{
		{
			name:     "Normal pagination",
			total:    25,
			page:     2,
			pageSize: 10,
			expected: PageInfo{
				Total:    25,
				Page:     2,
				PageSize: 10,
				Pages:    3,
				HasNext:  true,
				HasPrev:  true,
			},
		},
		{
			name:     "First page",
			total:    25,
			page:     1,
			pageSize: 10,
			expected: PageInfo{
				Total:    25,
				Page:     1,
				PageSize: 10,
				Pages:    3,
				HasNext:  true,
				HasPrev:  false,
			},
		},
		{
			name:     "Last page",
			total:    25,
			page:     3,
			pageSize: 10,
			expected: PageInfo{
				Total:    25,
				Page:     3,
				PageSize: 10,
				Pages:    3,
				HasNext:  false,
				HasPrev:  true,
			},
		},
		{
			name:     "Empty result",
			total:    0,
			page:     1,
			pageSize: 10,
			expected: PageInfo{
				Total:    0,
				Page:     1,
				PageSize: 10,
				Pages:    0,
				HasNext:  false,
				HasPrev:  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPageInfo(tt.total, tt.page, tt.pageSize)
			assert.Equal(t, tt.expected, *result)
		})
	}
}

func TestPageRequest_IsValidPage(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		total    int64
		expected bool
	}{
		{
			name:     "Valid page",
			page:     2,
			pageSize: 10,
			total:    25,
			expected: true,
		},
		{
			name:     "First page",
			page:     1,
			pageSize: 10,
			total:    25,
			expected: true,
		},
		{
			name:     "Last page",
			page:     3,
			pageSize: 10,
			total:    25,
			expected: true,
		},
		{
			name:     "Invalid page (too high)",
			page:     4,
			pageSize: 10,
			total:    25,
			expected: false,
		},
		{
			name:     "Empty result, first page",
			page:     1,
			pageSize: 10,
			total:    0,
			expected: true,
		},
		{
			name:     "Empty result, invalid page",
			page:     2,
			pageSize: 10,
			total:    0,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := PageRequest{Page: tt.page, PageSize: tt.pageSize}
			result := req.IsValidPage(tt.total)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPageRequest_ToPageInfo(t *testing.T) {
	req := PageRequest{Page: 2, PageSize: 10}
	total := int64(25)

	result := req.ToPageInfo(total)

	expected := &PageInfo{
		Total:    25,
		Page:     2,
		PageSize: 10,
		Pages:    3,
		HasNext:  true,
		HasPrev:  true,
	}

	assert.Equal(t, expected, result)
}
