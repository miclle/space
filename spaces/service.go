package spaces

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/fox-gonic/fox/database"
	"github.com/fox-gonic/fox/database/nestedset"
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
		}

		if err := nestedset.Create(tx, page, nil); err != nil {
			return err
		}

		content := &models.PageContent{
			SpaceID:   space.ID,
			CreatorID: params.CreatorID,
			PageID:    page.ID,
			Lang:      space.Lang,
			// Version:   space.Version, // TODO(m) space default version
			Status:     models.PageStatusPublished,
			Title:      space.Name,
			ShortTitle: space.Name,
			Body:       space.Description,
			HTML:       space.Description,
		}

		err = tx.Create(content).Error
		if err != nil {
			return err
		}

		err = tx.Model(space).Update("homepage_id", page.ID).Error
		if err != nil {
			return err
		}

		space.Homepage = page

		return nil
	})

	if err != nil {
		return nil, err
	}

	return space, nil
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

	{
		// get homepage content
		if params.Lang != "" {
			database = database.Joins("HomepageContent", s.Database.Select("title", "short_title").Where(&models.PageContent{Lang: params.Lang}))
		} else {
			database = database.Joins("HomepageContent", s.Database.Select("title", "short_title").Where("`HomepageContent`.`lang` = `spaces`.`lang`"))
		}

		// get fallback homepage content
		database = database.Joins("HomepageFallbackContent", s.Database.Select("title", "short_title").Where("`HomepageFallbackContent`.`lang` = `spaces`.`fallback_lang`"))
	}

	// Pagination
	database = database.Scopes(pagination.Paginate())

	if err := database.Find(&pagination.Items).Error; err != nil {
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
	err = database.Where("`space_id` = ? AND `id` = ?", space.ID, space.HomepageID).First(&space.Homepage).Error
	if err != nil {
		return nil, err
	}

	// find space homepage content
	err = database.Where("`page_id` = ? AND `lang` = ?", space.HomepageID, lang).First(&space.Homepage.Content).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) && lang != space.FallbackLang {
			err = database.Where("`page_id` = ? AND `lang` = ?", space.HomepageID, space.FallbackLang).First(&space.Homepage.Content).Error
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
	if params.Multilingual != nil {
		space.Multilingual = *params.Multilingual
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
		parent   *models.Page
		page     *models.Page
		content  *models.PageContent
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

	if params.ParentID > 0 {
		if err := database.Where("`id` = ?", params.ParentID).Find(&parent).Error; err != nil {
			return nil, err
		}
	}

	err = database.Transaction(func(tx *gorm.DB) error {

		page = &models.Page{
			SpaceID:  space.ID,
			ParentID: sql.NullInt64{Valid: true, Int64: params.ParentID},
		}

		if parent != nil {
			err = nestedset.Create(tx, page, parent)
		} else {
			err = nestedset.Create(tx, page, nil)
		}

		if err != nil {
			return err
		}

		content = &models.PageContent{
			SpaceID:    space.ID,
			CreatorID:  params.CreatorID,
			PageID:     page.ID,
			Lang:       space.Lang,
			Version:    params.Version,
			Status:     params.Status,
			Title:      params.Title,
			ShortTitle: params.ShortTitle,
			Body:       params.Body,
			HTML:       html,
		}

		if len(content.ShortTitle) == 0 {
			content.ShortTitle = content.Title
		}

		if content.Lang == "" {
			content.Lang = space.Lang
		}

		if err := tx.Create(&content).Error; err != nil {
			return err
		}

		// TODO(m) add history version record
		return nil
	})

	page.Content = content

	return page, err
}

func (s *service) DescribePages(ctx context.Context, params *params.DescribePages) ([]*models.Page, error) {

	var (
		database = s.Database.WithContext(ctx)
		space    *models.Space
	)

	if err := database.Where("`id` = ?", params.SpaceID).First(&space).Error; err != nil {
		return nil, err
	}

	lang := params.Lang
	if lang == "" {
		lang = space.Lang
	}

	// TODO(m)
	// if len(params.Version) > 0 {
	// 	database = database.Where("`version` = ?", params.Version)
	// }

	var pages models.Pages

	db := database.Joins("Content", database.Omit("body", "html").Where(&models.PageContent{Lang: lang}))

	if len(space.FallbackLang) > 0 && lang != space.FallbackLang {
		db = db.Joins("FallbackContent", database.Omit("body", "html").Where(&models.PageContent{Lang: space.FallbackLang}))
	}

	if params.Depth > 0 {
		db = db.Where("`space_pages`.`depth` <= ?", params.Depth)
	}

	if params.ParentID != nil {
		db = db.Where("`space_pages`.`parent_id` = ?", *params.ParentID)
	}

	err := db.Where("`space_pages`.`space_id` = ?", space.ID).Order("`lft` ASC").Find(&pages).Error
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		page.Space = space
	}

	return pages.Build(), nil
}

func (s *service) DescribePage(ctx context.Context, params *params.DescribePage) (*models.Page, error) {

	var (
		database = s.Database.WithContext(ctx)
		space    *models.Space
		page     *models.Page
	)

	if err := database.Where("`id` = ?", params.SpaceID).First(&space).Error; err != nil {
		return nil, err
	}

	lang := params.Lang
	if lang == "" {
		lang = space.Lang
	}

	if len(params.Version) > 0 {
		database = database.Where("`version` = ?", params.Version)
	}

	database = database.Joins("Content", database.Where(&models.PageContent{Lang: lang}))

	if len(space.FallbackLang) > 0 && lang != space.FallbackLang {
		database = database.Joins("FallbackContent", database.Where(&models.PageContent{Lang: space.FallbackLang}))
	}

	err := database.
		InstanceSet("query", &models.PageQuery{Lang: lang}).
		Where("`space_pages`.`space_id` = ? AND `page_id` = ?", space.ID, params.PageID).
		Find(&page).Error

	if err != nil {
		return nil, err
	}

	page.Space = space

	return page, nil
}

func (s *service) UpdatePage(ctx context.Context, params *params.UpdatePage) (*models.Page, error) {

	var (
		database = s.Database.WithContext(ctx)
		page     *models.Page
	)

	err := database.Where("`id` = ?", params.ID).Preload("Space").First(&page).Error
	if err != nil {
		return nil, err
	}

	if err := params.Status.IsValid(); err != nil {
		return nil, err
	}

	db := database.Where("`page_id` = ?", page.ID)

	if params.Lang != nil {
		db = db.Where("`lang` = ?", *params.Lang)
	}
	if params.Version != nil {
		db = db.Where("`version` = ?", *params.Version)
	}

	err = db.First(&page.Content).Error
	if err != nil {
		return nil, err
	}

	if params.Status != nil {
		page.Content.Status = *params.Status
	}
	if params.Title != nil {
		page.Content.Title = *params.Title
	}
	if params.ShortTitle != nil {
		page.Content.ShortTitle = *params.ShortTitle
	}
	if params.Body != nil {
		html, err := markdown.Parse(*params.Body)
		if err != nil {
			return nil, err
		}
		page.Content.Body = *params.Body
		page.Content.HTML = html
	}

	err = database.Save(page.Content).Error

	return page, err
}

func (s *service) Serach(ctx context.Context, params *params.Search) (*database.Pagination[*models.Page], error) {

	var (
		database   = s.Database.WithContext(ctx)
		pagination = &params.Pagination
		q          = strings.TrimSpace(params.Q)
		contents   []*models.PageContent
	)

	if len(q) == 0 {
		return pagination, nil
	}

	like := fmt.Sprintf("%%%s%%", q)
	database = database.Where("`lang` = ? AND (`title` LIKE ? OR `body` LIKE ?)", params.Lang, like, like)
	if err := database.Model(&contents).Count(&pagination.Total).Error; err != nil {
		return nil, err
	}

	database = database.
		Scopes(pagination.Paginate()).
		Preload("Space").
		Preload("Page", func(db *gorm.DB) *gorm.DB {
			return db.InstanceSet("query", &models.PageQuery{Lang: params.Lang})
		})

	if err := database.Find(&contents).Error; err != nil {
		return nil, err
	}

	for _, c := range contents {
		c.Page.Content = c
		pagination.Items = append(pagination.Items, c.Page)
	}

	return pagination, nil
}
