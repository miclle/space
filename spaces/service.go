package spaces

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fox-gonic/fox/database"
	"gorm.io/gorm"

	"github.com/miclle/space/models"
	"github.com/miclle/space/pkg/markdown"
	"github.com/miclle/space/spaces/params"
)

// Service for spaces interface
type Service interface {
	CreateSpace(context.Context, *params.CreateSpace) (*models.Space, error)
	DescribeSpaces(context.Context, *params.DescribeSpaces) (*database.Pagination[*models.Space], error)
	DescribeSpace(context.Context, *params.DescribeSpace) (*models.Space, error)
	UpdateSpace(context.Context, *params.UpdateSpace) (*models.Space, error)

	CreatePage(context.Context, *params.CreatePage) (*models.Page, error)
	DescribePages(context.Context, *params.DescribePages) ([]*models.Page, error)
	DescribePage(context.Context, *params.DescribePage) (*models.Page, error)
	UpdatePage(context.Context, *params.UpdatePage) (*models.Page, error)

	Serach(context.Context, *params.Search) (*database.Pagination[*models.Page], error)
}

// NewService return default implement spaces service
func NewService(database *database.Database) (Service, error) {

	service := &service{
		Database: database,
	}

	return service, nil
}

var _ Service = &service{}

type service struct {
	Database *database.Database
}

func (s *service) CreateSpace(ctx context.Context, params *params.CreateSpace) (*models.Space, error) {

	var (
		database = s.Database.WithContext(ctx)
		space    *models.Space
	)

	if err := params.Status.IsValid(); err != nil {
		return nil, err
	}

	err := database.Transaction(func(tx *gorm.DB) error {

		space = &models.Space{
			Name:         params.Name,
			Key:          params.Key,
			Lang:         params.Lang,
			FallbackLang: params.FallbackLang,
			Description:  params.Description,
			Avatar:       params.Avatar,
			Status:       params.Status,
			CreatorID:    params.CreatorID,
		}

		err := tx.Create(space).Error
		if err != nil {
			return err
		}

		page := &models.Page{
			SpaceID: space.ID,
			Lang:    params.Lang,
			// Version:   space.Version, // TODO(m) space default version
			Status:     models.PageStatusPublished,
			Title:      space.Name,
			ShortTitle: space.Name,
			Body:       space.Description,
			HTML:       space.Description,
			CreatorID:  params.CreatorID,
		}

		err = tx.Create(page).Error
		if err != nil {
			return err
		}

		err = tx.Model(space).Update("homepage_id", page.PageID).Error
		if err != nil {
			return err
		}

		space.Homepage = page

		return nil
	})

	return space, err
}

func (s *service) DescribeSpaces(ctx context.Context, params *params.DescribeSpaces) (*database.Pagination[*models.Space], error) {

	var (
		database   = s.Database.WithContext(ctx)
		pagination = &params.Pagination
	)

	if q := strings.TrimSpace(params.Q); q != "" {
		like := fmt.Sprintf("%%%s%%", q)
		database = database.Where("`spaces`.`name` LIKE ? OR `spaces`.`key` LIKE ?", like, like)
	}

	if err := database.Model(&pagination.Items).Count(&pagination.Total).Error; err != nil {
		return nil, err
	}

	// Pagination
	database = database.Scopes(pagination.Paginate())

	// TODO(m) preload homepage space lang or fallback lang
	if err := database.Preload("Homepage", func(db *gorm.DB) *gorm.DB {
		if params.Lang != "" {
			db = db.Where("`lang` = ?", params.Lang)
		}
		return db
	}).Find(&pagination.Items).Error; err != nil {
		return nil, err
	}

	return pagination, nil
}

func (s *service) DescribeSpace(ctx context.Context, params *params.DescribeSpace) (*models.Space, error) {

	var (
		database = s.Database.WithContext(ctx)
		space    *models.Space
	)

	// find space
	err := database.Where("`key` = ?", params.Key).First(&space).Error
	if err != nil {
		return nil, err
	}

	lang := params.Lang
	if lang == "" {
		lang = space.Lang
	}

	if len(params.Version) > 0 {
		database = database.Where("`version` = ?", params.Version)
	}

	// find space homepage
	err = database.Where("`space_id` = ? AND `page_id` = ? AND `lang` = ?", space.ID, space.HomepageID, lang).First(&space.Homepage).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) && lang != space.FallbackLang {
			err = database.Where("`space_id` = ? AND `page_id` = ? AND `lang` = ?", space.ID, space.HomepageID, space.FallbackLang).First(&space.Homepage).Error
		}
		return nil, err
	}

	return space, err
}

func (s *service) UpdateSpace(ctx context.Context, params *params.UpdateSpace) (*models.Space, error) {

	var (
		database = s.Database.WithContext(ctx)
		space    *models.Space
	)

	err := database.Where("`key` = ?", params.Key).First(&space).Error
	if err != nil {
		return nil, err
	}

	if params.Name != nil {
		space.Name = *params.Name
	}
	if params.Lang != nil {
		space.Lang = *params.Lang
	}
	if params.FallbackLang != nil {
		space.FallbackLang = *params.FallbackLang
	}
	if params.HomepageID != nil {
		space.HomepageID = *params.HomepageID
	}
	if params.Description != nil {
		space.Description = *params.Description
	}
	if params.Avatar != nil {
		space.Avatar = *params.Avatar
	}
	if params.Status.IsValid() == nil {
		space.Status = params.Status
	}

	err = database.Save(space).Error

	return space, err
}

func (s *service) CreatePage(ctx context.Context, params *params.CreatePage) (*models.Page, error) {

	var (
		database = s.Database.WithContext(ctx)
		space    *models.Space
		page     *models.Page
	)

	if err := params.Status.IsValid(); err != nil {
		return nil, err
	}

	html, err := markdown.Parse(params.Body)
	if err != nil {
		return nil, err
	}

	err = database.Where("`id` = ?", params.SpaceID).First(&space).Error
	if err != nil {
		return nil, err
	}

	err = database.Transaction(func(tx *gorm.DB) error {

		page = &models.Page{
			SpaceID:      params.SpaceID,
			CreatorID:    params.CreatorID,
			ParentPageID: params.ParentID,
			PageID:       params.PageID,
			Lang:         params.Lang,
			Version:      params.Version,
			Status:       params.Status,
			Title:        params.Title,
			ShortTitle:   params.ShortTitle,
			Body:         params.Body,
			HTML:         html,
		}

		if len(page.ShortTitle) == 0 {
			page.ShortTitle = page.Title
		}

		if page.Lang == "" {
			page.Lang = space.Lang
		}

		if err := tx.Create(&page).Error; err != nil {
			return err
		}

		// TODO(m) add history version record
		return nil
	})

	return page, err
}

func (s *service) DescribePages(ctx context.Context, params *params.DescribePages) ([]*models.Page, error) {

	var (
		database = s.Database.WithContext(ctx)
		space    *models.Space
		pages    []*models.Page
	)

	err := database.Where("`id` = ?", params.SpaceID).First(&space).Error
	if err != nil {
		return nil, err
	}

	lang := params.Lang
	if lang == "" {
		lang = space.Lang
	}
	database = database.Where("`lang` = ?", lang)

	if len(params.Version) > 0 {
		database = database.Where("`version` = ?", params.Version)
	}

	err = database.Omit("html").Where("`space_id` = ?", space.ID).Find(&pages).Error

	return pages, err
}

func (s *service) DescribePage(ctx context.Context, params *params.DescribePage) (*models.Page, error) {

	var (
		database = s.Database.WithContext(ctx)
		space    *models.Space
		page     *models.Page
	)

	err := database.Where("`id` = ?", params.SpaceID).First(&space).Error
	if err != nil {
		return nil, err
	}

	lang := params.Lang
	if lang == "" {
		lang = space.Lang
	}

	if len(params.Version) > 0 {
		database = database.Where("`version` = ?", params.Version)
	}

	err = database.Where("`space_id` = ? AND `page_id` = ? AND `lang` = ?", space.ID, params.PageID, lang).First(&page).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = database.Where("`space_id` = ? AND `page_id` = ? AND `lang` = ?", space.ID, params.PageID, space.FallbackLang).First(&page).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
	}

	if err != nil {
		return nil, err
	}

	return page, nil
}

func (s *service) UpdatePage(ctx context.Context, params *params.UpdatePage) (*models.Page, error) {

	var (
		database = s.Database.WithContext(ctx)
		page     *models.Page
	)

	err := database.Where("`id` = ?", params.ID).First(&page).Error
	if err != nil {
		return nil, err
	}

	if err := params.Status.IsValid(); err != nil {
		return nil, err
	}

	if params.Lang != nil {
		page.Lang = *params.Lang
	}
	if params.Version != nil {
		page.Version = *params.Version
	}
	if params.Status != nil {
		page.Status = *params.Status
	}
	if params.Title != nil {
		page.Title = *params.Title
	}
	if params.ShortTitle != nil {
		page.ShortTitle = *params.ShortTitle
	}
	if params.Body != nil {
		html, err := markdown.Parse(*params.Body)
		if err != nil {
			return nil, err
		}
		page.Body = *params.Body
		page.HTML = html
	}

	err = database.Save(page).Error

	return page, err
}

func (s *service) Serach(ctx context.Context, params *params.Search) (*database.Pagination[*models.Page], error) {

	var (
		database   = s.Database.WithContext(ctx)
		pagination = &params.Pagination
		q          = strings.TrimSpace(params.Q)
	)

	if len(q) == 0 {
		return pagination, nil
	}

	like := fmt.Sprintf("%%%s%%", q)
	database = database.Where("`lang` = ? AND (`title` LIKE ? OR `body` LIKE ?)", params.Lang, like, like)

	if err := database.Model(&pagination.Items).Count(&pagination.Total).Error; err != nil {
		return nil, err
	}

	database = database.
		Scopes(pagination.Paginate()).
		Preload("Space")

	if err := database.Find(&pagination.Items).Error; err != nil {
		return nil, err
	}

	return pagination, nil
}
