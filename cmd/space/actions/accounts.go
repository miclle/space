package actions

import (
	"github.com/fox-gonic/fox/engine"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin/render"
	"github.com/miclle/space/accounts/params"
	"github.com/miclle/space/models"
)

// SignupArgs user sign up args
type SignupArgs struct {
	Login    string `json:"login"    binding:"required" validate:"required"`
	Email    string `json:"email"    binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=8,max=56"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
	Location string `json:"location"`
}

// Signup user sign up
// POST /accounts/signup
func (actions *Actions) Signup(c *engine.Context, args *SignupArgs) (*models.Account, error) {

	account, err := actions.Accounter.CreateAccount(c, &params.CreateAccount{
		Login:    args.Login,
		Email:    args.Email,
		Password: args.Password,
		Name:     args.Name,
		Bio:      args.Bio,
		Location: args.Location,
	})
	if err != nil {
		return nil, err
	}

	session := sessions.Default(c.Context)
	session.Set(SessionAccountKey, account.Login)
	if err := session.Save(); err != nil {
		c.Logger.Error("session.Save() failed, err: %+v", err)
		return nil, err
	}

	// TODO(m) Send email validations notification

	return account, nil
}

// ----------------------------------------------------------------------------

// SigninArgs user sign in args
type SigninArgs struct {
	Login    string `json:"login"    binding:"required" validate:"required"`
	Email    string `json:"email"    binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=8,max=56"`
}

// Signin user sign in
// POST /accounts/signin
func (actions *Actions) Signin(c *engine.Context, args *SigninArgs) (*models.Account, error) {

	account, err := actions.Accounter.AuthenticateAccount(c, &params.AuthenticateAccount{
		Login:    args.Login,
		Email:    args.Email,
		Password: args.Password,
		ClientIP: c.ClientIP(),
	})
	if err != nil {
		return nil, err
	}

	session := sessions.Default(c.Context)
	session.Set(SessionAccountKey, account.Login)
	if err := session.Save(); err != nil {
		c.Logger.Error("session.Save() failed, err: %+v", err)
		return nil, err
	}

	// TODO(m) Send login email notification

	return account, nil
}

// Logout user logout
func (actions *Actions) Logout(c *engine.Context) (res interface{}) {

	session := sessions.Default(c.Context)
	session.Clear()

	if err := session.Save(); err != nil {
		c.Logger.Error("session.Save() failed, err: %+v", err)
	}

	return render.Redirect{
		Code:     302,
		Location: "/",
	}
}
