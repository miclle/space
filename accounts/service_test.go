package accounts

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fox-gonic/fox/database"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/miclle/space/accounts/params"
	"github.com/miclle/space/models"
)

var accounter Service

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

	accounter, err = NewService(database)
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

func TestCreateAccount(t *testing.T) {
	assert := assert.New(t)

	account, err := accounter.CreateAccount(context.Background(), &params.CreateAccount{})
	assert.NotNil(err)
	assert.Nil(account)

	account, err = accounter.CreateAccount(context.Background(), &params.CreateAccount{
		Login:    "lisa",
		Email:    "lisa@domain.local",
		Password: "lisa_Password",
		Name:     "Lisa",
	})
	assert.Nil(err)
	assert.NotNil(account)
	assert.Equal("lisa", account.Login)
	assert.Equal("lisa@domain.local", account.Email)
	assert.Equal("Lisa", account.Name)

	account, err = accounter.CreateAccount(context.Background(), &params.CreateAccount{
		Login:    "lisa",
		Email:    "lisa@domain.local",
		Password: "lisa_Password",
		Name:     "Lisa",
	})
	assert.NotNil(err)
}

func TestDescribeAccounts(t *testing.T) {
	assert := assert.New(t)

	pagination, err := accounter.DescribeAccounts(context.Background(), &params.DescribeAccounts{})
	assert.Nil(err)
	assert.Len(pagination.Items, 1)

	for _, account := range pagination.Items {
		fmt.Printf("account: %+v\n", account)
	}

	pagination, err = accounter.DescribeAccounts(context.Background(), &params.DescribeAccounts{
		Q: "Lisa",
	})
	assert.Nil(err)
	assert.Len(pagination.Items, 1)

	pagination, err = accounter.DescribeAccounts(context.Background(), &params.DescribeAccounts{
		Q: "Lili",
	})
	assert.Nil(err)
	assert.Len(pagination.Items, 0)
}

func TestDescribeAccount(t *testing.T) {
	assert := assert.New(t)

	account, err := accounter.DescribeAccount(context.Background(), &params.DescribeAccount{})
	assert.NotNil(err)
	assert.Equal(params.ErrDescribeAccountParamsInvalid, err)
	assert.Nil(account)

	account, err = accounter.DescribeAccount(context.Background(), &params.DescribeAccount{
		Login: "lili",
	})
	assert.NotNil(err)
	assert.Equal(gorm.ErrRecordNotFound, err)
	assert.Nil(account)

	account, err = accounter.DescribeAccount(context.Background(), &params.DescribeAccount{
		Login: "lisa",
	})
	assert.Nil(err)
	assert.NotNil(account)
}

func TestAuthenticateAccount(t *testing.T) {
	assert := assert.New(t)

	account, err := accounter.AuthenticateAccount(context.Background(), &params.AuthenticateAccount{
		Login:    "lili",
		Password: "lili_Password",
		ClientIP: "",
	})
	assert.NotNil(err)
	assert.Equal(errors.New(http.StatusText(http.StatusUnauthorized)), err)
	assert.Nil(account)

	account, err = accounter.AuthenticateAccount(context.Background(), &params.AuthenticateAccount{
		Login:    "lisa",
		Password: "unauthorized",
		ClientIP: "",
	})
	assert.NotNil(err)
	assert.Equal(errors.New(http.StatusText(http.StatusUnauthorized)), err)
	assert.Nil(account)

	account, err = accounter.AuthenticateAccount(context.Background(), &params.AuthenticateAccount{
		Login:    "lisa",
		Password: "lisa_Password",
		ClientIP: "",
	})
	assert.Nil(err)
	assert.NotNil(account)
}

func TestUpdateAccount(t *testing.T) {
	assert := assert.New(t)

	account, err := accounter.UpdateAccount(context.Background(), &params.UpdateAccount{})
	assert.NotNil(err)
	assert.Nil(account)

	account, err = accounter.UpdateAccount(context.Background(), &params.UpdateAccount{
		Login: "lisa",
		Name:  lo.ToPtr("Mona Lisa"),
	})
	assert.Nil(err)
	assert.NotNil(account)
	assert.Equal("Mona Lisa", account.Name)
}

func TestCreateUnlockToken(t *testing.T) {
	assert := assert.New(t)

	token, err := accounter.CreateUnlockToken(context.Background(), &params.CreateUnlockToken{
		Login: "lisa",
	})
	assert.Nil(err)
	assert.NotEmpty(token)
}

func TestUnlock(t *testing.T) {
	// TODO(m)
}

func TestUpdatePassword(t *testing.T) {
	// TODO(m)
}

func TestCreateResetPasswordToken(t *testing.T) {
	// TODO(m)
}

func TestResetPassword(t *testing.T) {
	// TODO(m)
}
