package controller

type Controller struct {
	Auth AuthController
	User UserController
	Chat ChatController
}

func NewController(auth AuthController, user UserController, chat ChatController) *Controller {
	return &Controller{
		Auth: auth,
		User: user,
		Chat: chat,
	}
}
