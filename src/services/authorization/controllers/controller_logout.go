package controllers

import "context"

func (c *controller) LogoutUser(ctx context.Context, request *RequestLogoutService) (response *ResponseLogoutService, err error) {
	if err := c.rdb.Del(ctx, request.UserID).Err(); err != nil {
		return nil, err
	}

	return &ResponseLogoutService{}, nil
}