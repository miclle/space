package spaces

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fox-gonic/fox/database"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/miclle/space/models"
	"github.com/miclle/space/spaces/params"
)

var spacer Service

func TestMain(m *testing.M) {

	db := filepath.Join(os.TempDir(), fmt.Sprintf("space_test_%d.db", time.Now().Unix()))
	log.Println(db)

	database, err := database.New(&database.Config{
		Dialect:  "sqlite",
		Database: db,
	})
	if err != nil {
		log.Fatalf("database init failed, err: %+v", err)
	}

	err = models.Migrate(database)
	if err != nil {
		log.Fatalf("migrate models failed, err: %+v", err)
	}

	spacer, err = NewService(database)
	if err != nil {
		log.Fatalf("new space service failed, err: %+v", err)
	}

	code := m.Run()

	if code == 0 {
		if err := os.Remove(db); err != nil {
			log.Printf("remove testing db: %s failed, err: %s", db, err.Error())
		}
	} else {
		log.Printf("The test failed and you need to manually delete the test database: %s\n", db)
	}

	// os.Exit(code)
}

func TestCreateSpace(t *testing.T) {
	assert := assert.New(t)

	space, err := spacer.CreateSpace(context.Background(), &params.CreateSpace{})
	assert.NotNil(err)
	assert.Nil(space)

	space, err = spacer.CreateSpace(context.Background(), &params.CreateSpace{
		Name:   "Website",
		Key:    "website",
		Status: models.SpaceStatusOnline,
		Lang:   "en-US",
	})
	assert.Nil(err)
	assert.NotNil(space)
	assert.Equal("Website", space.Name)
	assert.Equal("website", space.Key)
	assert.NotNil(space.Homepage)
	assert.Equal("Website", space.Homepage.Title)

	space, err = spacer.CreateSpace(context.Background(), &params.CreateSpace{
		Name:   "Website",
		Key:    "website",
		Status: models.SpaceStatusOnline,
		Lang:   "en-US",
	})
	assert.NotNil(err)
}

func TestDescribeSpaces(t *testing.T) {
	assert := assert.New(t)

	pagination, err := spacer.DescribeSpaces(context.Background(), &params.DescribeSpaces{})
	assert.Nil(err)
	assert.Len(pagination.Items, 1)

	for _, space := range pagination.Items {
		fmt.Printf("space.Homepage: %+v\n", space.Homepage)
	}

	pagination, err = spacer.DescribeSpaces(context.Background(), &params.DescribeSpaces{
		Q: "web",
	})
	assert.Nil(err)
	assert.Len(pagination.Items, 1)

	pagination, err = spacer.DescribeSpaces(context.Background(), &params.DescribeSpaces{
		Q: "web3",
	})
	assert.Nil(err)
	assert.Len(pagination.Items, 0)
}

func TestDescribeSpace(t *testing.T) {
	assert := assert.New(t)

	space, err := spacer.DescribeSpace(context.Background(), &params.DescribeSpace{})
	assert.NotNil(err)
	assert.Equal(gorm.ErrRecordNotFound, err)
	assert.Nil(space)

	space, err = spacer.DescribeSpace(context.Background(), &params.DescribeSpace{
		Key: "website",
	})
	assert.Nil(err)
	assert.NotNil(space)
	assert.NotNil(space.Homepage)
}

func TestUpdateSpace(t *testing.T) {
	assert := assert.New(t)

	space, err := spacer.UpdateSpace(context.Background(), &params.UpdateSpace{})
	assert.NotNil(err)
	assert.Nil(space)

	space, err = spacer.UpdateSpace(context.Background(), &params.UpdateSpace{
		Key:         "website",
		Description: lo.ToPtr("website space description"),
	})
	assert.Nil(err)
	assert.NotNil(space)
	assert.Equal("website space description", space.Description)
}

func TestCreatePage(t *testing.T) {
	// TODO(m)
}

func TestDescribePages(t *testing.T) {
	// TODO(m)
}

func TestDescribePage(t *testing.T) {
	// TODO(m)
}

func TestUpdatePage(t *testing.T) {
	// TODO(m)
}

func TestSerach(t *testing.T) {
	// TODO(m)
}
