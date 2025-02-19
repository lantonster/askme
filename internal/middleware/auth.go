package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lantonster/askme/internal/model"
)

var ctxUUIDKey = "ctxUuidKey"

// GetUserIsAdminModerator 检查 Gin 上下文中的用户是否为管理员或版主。
//
// 参数:
//   - c: Gin 上下文
//
// 返回:
//   - bool: 如果用户是管理员或版主则返回 true，否则返回 false
func GetUserIsAdminModerator(c *gin.Context) bool {
	// 从上下文中获取用户信息，如果用户信息不存在，返回 false
	userInfo, exist := c.Get(ctxUUIDKey)
	if !exist {
		return false
	}

	// 将获取到的用户信息转换为指定的用户信息类型，如果转换不成功，返回 false
	user, ok := userInfo.(*model.UserInfo)
	if !ok {
		return false
	}

	// 如果用户的角色 ID 是管理员或版主的 ID，返回 true
	if user.RoleId == model.RoleIdAdminID || user.RoleId == model.RoleIdModeratorID {
		return true
	}

	return false
}
