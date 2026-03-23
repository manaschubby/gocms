package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/manaschubby/gocms/internal/modules/cms/common"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
	"github.com/manaschubby/gocms/internal/modules/cms/repository"
)

func ValidateAccount(r repository.AccountRepository, ctx context.Context, aid uuid.UUID) (int, error) {
	account, err := r.GetAccountByUUID(aid, repository.GetAccountOptions{Context: &ctx})
	if err != nil {
		log.Printf("failed to retrieve account data: %v", err)
		return http.StatusInternalServerError, errors.New("failed to retrieve account data")
	}

	if account == nil {
		return http.StatusBadRequest, errors.New("account not found")
	}
	return 0, nil
}

func ValidateContentType(r repository.ContentTypeRepository, ctx context.Context, ctId uuid.UUID) (*domain.ContentType, int, error) {
	contentType, err := r.GetContentTypeById(ctId, repository.GetContentTypeOptions{Context: &ctx})
	if err != nil {
		log.Printf("failed to fetch contentType from database: %v", err)
		return nil, http.StatusInternalServerError, errors.New("failed to fetch contentType from database")
	}

	if contentType == nil {
		return nil, http.StatusBadRequest, errors.New("contentType does not exist")
	}

	return contentType, 0, nil
}

func ValidateContentData(cd json.RawMessage, sm domain.SchemaDefinitionMap) (json.RawMessage, error) {
	var contentData map[string]any
	err := common.JsonNumberDecode(cd, &contentData)
	if err != nil { // should not happen (callers should check for validation before passing)
		return nil, fmt.Errorf("failed to parse content data: %w", err)
	}

	for key := range contentData {
		if _, exists := sm[key]; !exists {
			return nil, errors.New("extra columns found in contentData")
		}
	}

	parsed := map[string]any{}
	errorColumns := make([]string, 0)
	for k, v := range sm {
		value := contentData[k]
		if v.DefaultValue != nil && value == nil {
			value = v.DefaultValue
		}
		if v.Required && value == nil {
			errorColumns = append(errorColumns, k+": required field and no default value exists")
			continue
		}
		if value == nil {
			continue
		}

		err = v.ValidateAny(value)
		log.Println(value, err)
		if err != nil {
			errorColumns = append(errorColumns, k+": "+err.Error())
		}
		parsed[k] = value
	}

	if len(errorColumns) != 0 {
		return nil, errors.New("schema validation failed for following columns: {" + strings.Join(errorColumns, ", ") + "}")
	}

	cd, err = json.Marshal(parsed)
	if err != nil {
		log.Printf("failed to prepare contentData for db insert: %v", err)
		return nil, errors.New("failed to prepare contentData for db insert")
	}

	return cd, nil
}
