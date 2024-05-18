package core

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"micro-todolist/user/model"
	"micro-todolist/user/services"
)

// 将model.User序列化返回
func BuildUser(item model.User) *services.UserModel {
	userModel := services.UserModel{
		ID:        uint32(item.ID),
		UserName:  item.UserName,
		CreatedAt: item.CreatedAt.Unix(),
		UpdatedAt: item.UpdatedAt.Unix(),
	}
	return &userModel
}

// 登录
func (*UserService) UserLogin(ctx context.Context, req *services.UserRequest, resp *services.UserDetailResponse) error {
	var user model.User
	resp.Code = 200
	if err := model.DB.Where("user_name=?", req.UserName).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			resp.Code = 400 //前端请求出错
			err := errors.New("没有此账户")
			return err //right?
		}
		resp.Code = 500
		err := errors.New("服务器出错")
		return err //right?
	}
	if user.CheckPassword(req.Password) == false {
		resp.Code = 400
		err := errors.New("密码错误")
		return err //right?
	}
	resp.UserDetail = BuildUser(user)
	return nil //

}

// 注册
func (*UserService) UserRegister(ctx context.Context, req *services.UserRequest, resp *services.UserDetailResponse) error {
	if req.Password != req.PasswordConfirm {
		err := errors.New("两次输入密码不一致")
		return err
	}
	count := 0
	if err := model.DB.Model(&model.User{}).Where("user_name=?", req.UserName).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		err := errors.New("用户名已存在")
		return err
	}
	user := model.User{
		UserName: req.UserName,
	}
	//密码加密
	if err := user.SetPassword(req.Password); err != nil {
		return err
	}
	if err := model.DB.Count(&user).Error; err != nil {
		return err
	}
	resp.UserDetail = BuildUser(user)
	return nil
}
