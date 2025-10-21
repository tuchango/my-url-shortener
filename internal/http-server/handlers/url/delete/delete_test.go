package delete_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"github.com/tuchango/my-url-shortener/internal/http-server/handlers/url/delete"
	"github.com/tuchango/my-url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/tuchango/my-url-shortener/internal/storage"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respCode  int
		mockError error
	}{
		{
			name:     "Success",
			alias:    "test_alias",
			respCode: http.StatusNoContent,
		},
		{
			name:     "Empty alias",
			alias:    "",
			respCode: http.StatusBadRequest,
		},
		{
			name:      "URL not found",
			alias:     "non_existent_alias",
			respCode:  http.StatusNotFound,
			mockError: storage.ErrURLNotFound,
		},
		{
			name:      "Internal error",
			alias:     "test_alias",
			respCode:  http.StatusInternalServerError,
			mockError: errors.New("unexpected database error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlDeleterMock := delete.NewMockURLDeleter(t)

			// Настраиваем mock только если ожидаем вызов
			if tc.alias != "" {
				urlDeleterMock.On("DeleteURL", tc.alias).
					Return(tc.mockError).
					Once()
			}

			// Создаем router и handler
			r := chi.NewRouter()
			handler := delete.New(slogdiscard.NewDiscardLogger(), urlDeleterMock)
			r.Delete("/{alias}", handler)
			r.Delete("/", handler)

			// 3. Создание и выполнение запроса (без создания тестового сервера)
			// Используем любой валидный хост, так как мы вызываем роутер напрямую.
			targetURL := "http://example.com/"
			if tc.alias != "" {
				targetURL += tc.alias
			}

			req, err := http.NewRequest(http.MethodDelete, targetURL, nil)
			require.NoError(t, err)

			// Выполняем запрос
			respRecorder := httptest.NewRecorder()
			r.ServeHTTP(respRecorder, req)

			// Проверяем код ответа
			require.Equal(t, tc.respCode, respRecorder.Code)

			// Для ошибки с пустым alias проверяем содержимое
			if tc.alias == "" {
				require.Contains(t, respRecorder.Body.String(), "invalid request")
			}
		})
	}
}
