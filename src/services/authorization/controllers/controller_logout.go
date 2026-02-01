package controllers

import "context"

func (c *controller) LogoutUser(ctx context.Context, request *RequestLogoutService) (response *ResponseLogoutService, err error) {
	c.rdb.Del(ctx, request.UserID)

	return &ResponseLogoutService{}, nil
}