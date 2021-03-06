package service

import (
	"errors"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"strings"
)

// 登录校验
func (s *MysqlService) LoginCheck(user *models.SysUser) (*models.SysUser, error) {
	var u models.SysUser
	// 查询用户及其角色
	err := s.tx.Preload("Role").Where("username = ?", user.Username).First(&u).Error
	if err != nil {
		return nil, errors.New(response.LoginCheckErrorMsg)
	}
	// 校验密码
	if ok := utils.ComparePwd(user.Password, u.Password); !ok {
		return nil, errors.New(response.LoginCheckErrorMsg)
	}
	return &u, err
}

// 获取用户
func (s *MysqlService) GetUsers(req *request.UserListRequestStruct) ([]models.SysUser, error) {
	var err error
	list := make([]models.SysUser, 0)
	db := global.Mysql.
		Model(models.SysUser{}).
		Order("created_at DESC")
	// 非超级管理员
	if *req.CurrentRole.Sort != models.SysRoleSuperAdminSort {
		roleIds, err := s.GetRoleIdsBySort(*req.CurrentRole.Sort)
		if err != nil {
			return list, err
		}
		db = db.Where("role_id IN (?)", roleIds)
	}
	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	mobile := strings.TrimSpace(req.Mobile)
	if mobile != "" {
		db = db.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	if req.Status != nil {
		if *req.Status > 0 {
			db = db.Where("status = ?", 1)
		} else {
			db = db.Where("status = ?", 0)
		}
	}
	// 查询条数
	err = db.Count(&req.PageInfo.Total).Error
	if err == nil {
		if req.PageInfo.NoPagination {
			// 不使用分页
			err = db.Find(&list).Error
		} else {
			// 获取分页参数
			limit, offset := req.GetLimit()
			err = db.Limit(limit).Offset(offset).Find(&list).Error
		}
	}
	return list, err
}

// 获取单个用户
func (s *MysqlService) GetUserById(id uint) (models.SysUser, error) {
	var user models.SysUser
	var err error
	err = s.tx.Preload("Role").
		Where("id = ?", id).
		// 状态为正常
		Where("status = ?", models.SysUserStatusNormal).
		First(&user).Error
	return user, err
}

// 获取多个用户
func (s *MysqlService) GetUsersByIds(ids []uint) ([]models.SysUser, error) {
	var users []models.SysUser
	var err error
	err = s.tx.Preload("Role").Where("id IN (?)", ids).Find(&users).Error
	return users, err
}
