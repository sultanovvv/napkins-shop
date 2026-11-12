package category

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
	"napkins-shop/backoffice-executor/internal/customerr"
)

// readJSON сохраняет наличие полей для PATCH и декодирует контракт API.
// Query-параметры и поля формы не участвуют в изменении категории.
func readJSON(c echo.Context, body any, allowedFields ...string) (map[string]json.RawMessage, error) {
	mediaType, _, err := mime.ParseMediaType(c.Request().Header.Get(echo.HeaderContentType))
	if err != nil || mediaType != echo.MIMEApplicationJSON {
		return nil, invalidField("body", "Content-Type must be application/json", err)
	}
	decoder := json.NewDecoder(c.Request().Body)
	var raw json.RawMessage
	if err := decoder.Decode(&raw); err != nil {
		return nil, invalidField("body", "must be a JSON object", err)
	}
	var extra json.RawMessage
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return nil, invalidField("body", "must contain exactly one JSON object", err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil || fields == nil {
		return nil, invalidField("body", "must be a JSON object", err)
	}
	// encoding/json допускает несовпадение регистра имён. Контракт API
	// и определение присутствия PATCH-поля используют точное имя.
	for field := range fields {
		if !slices.Contains(allowedFields, field) {
			return nil, invalidField(field, "unknown field", nil)
		}
	}
	decoder = json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(body); err != nil {
		return nil, invalidField("body", "contains unknown fields or values of an invalid type", err)
	}
	return fields, nil
}

func parseID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 || strconv.FormatInt(id, 10) != value {
		if err == nil {
			err = fmt.Errorf("invalid category id: %q", value)
		}
		return 0, invalidField("id", "must be a positive int64 in decimal notation", err)
	}
	return id, nil
}

func isNull(raw json.RawMessage) bool {
	return bytes.Equal(bytes.TrimSpace(raw), []byte("null"))
}

func invalidField(target, message string, cause error) *customerr.Error {
	err := customerr.NewValidation("Request validation failed", cause)
	err.Details = []customerr.Detail{{Code: "INVALID_VALUE", Message: message, Target: target}}
	return err
}
