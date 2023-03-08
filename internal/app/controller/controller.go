package controller

type Controller struct {
	Auth AuthController
	User UserController
}

func NewController(auth AuthController, user UserController) *Controller {
	return &Controller{
		Auth: auth,
		User: user,
	}
}
