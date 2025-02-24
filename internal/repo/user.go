package repo

import (
	"context"
	"strings"

	"github.com/lantonster/askme/internal/data"
	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/pkg/checker"
	"github.com/lantonster/askme/pkg/errors"
	"github.com/lantonster/askme/pkg/log"
	"github.com/lantonster/askme/pkg/orm"
	"github.com/lantonster/askme/pkg/random"
	"github.com/lantonster/askme/pkg/reason"
	"github.com/mozillazg/go-pinyin"
	"gorm.io/gorm"
)

type UserRepo interface {
	// CreateUser 函数用于在数据库中创建新用户。
	CreateUser(c context.Context, user *model.User) error

	// GetUserByEmail 函数根据给定的电子邮件从数据库中获取用户信息。
	GetUserByEmail(c context.Context, email string) (user *model.User, exist bool, err error)

	// GetUserById 根据用户 ID 获取用户信息。
	GetUserById(c context.Context, userId int64) (*model.User, bool, error)

	// GetUserByUsername 函数根据给定的用户名从数据库中获取用户信息。
	GetUserByUsername(c context.Context, username string) (user *model.User, exist bool, err error)

	// GenerateUniqueUsername 将给定的用户名处理为唯一的有效用户名。
	GenerateUniqueUsername(c context.Context, username string) (string, error)

	// IncrRank 用于增加用户的积分。
	IncrRank(c context.Context, userId, currentRank, deltaRank int64) error

	// UpdateEmailStatus 更新用户的邮箱状态
	UpdateEmailStatus(c context.Context, userId int64, emailStatus string) error
}

type userRepo struct {
	*data.Data
}

func NewUserRepo(data *data.Data) UserRepo {
	return &userRepo{Data: data}
}

// CreateUser 函数用于在数据库中创建新用户。
//
// 参数:
//   - c: 上下文
//   - user: 要创建的用户对象
//
// 返回: 可能返回创建用户时发生的错误
func (r *userRepo) CreateUser(c context.Context, user *model.User) error {
	// 在数据库中创建用户
	if err := orm.Q.User.WithContext(c).Create(user); err != nil {
		// 如果是因为重复键导致的错误，返回特定错误
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.InternalServer(reason.UsernameDuplicate).WithError(err).WithStack()
		}
		// 其他数据库错误，返回通用的数据库错误
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	// 创建成功，返回 nil
	return nil
}

// GetUserByEmail 函数根据给定的电子邮件从数据库中获取用户信息。
//
// 参数:
//   - c: 上下文
//   - email: 要查找的用户的电子邮件
//
// 返回:
//   - *model.User: 用户对象，如果未找到则为 nil
//   - exist: 是否存在该用户
//   - error: 可能出现的错误
func (r *userRepo) GetUserByEmail(c context.Context, email string) (user *model.User, exist bool, err error) {
	users, err := orm.Q.User.WithContext(c).Where(orm.Q.User.Email.Eq(email)).Find()
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	if len(users) == 0 {
		return nil, false, nil
	}
	return users[0], true, nil
}

// GetUserById 根据用户 ID 获取用户信息。
//
// 参数:
//   - c: 上下文
//   - userId: 用户 ID
//
// 返回:
//   - *model.User: 用户对象，如果未找到则为 nil
//   - exist: 是否存在该用户
//   - error: 可能出现的错误
func (r *userRepo) GetUserById(c context.Context, userId int64) (*model.User, bool, error) {
	users, err := orm.Q.User.WithContext(c).Where(orm.Q.User.Id.Eq(userId)).Find()
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	if len(users) == 0 {
		return nil, false, nil
	}
	return users[0], true, nil
}

// GetUserByUsername 函数根据给定的用户名从数据库中获取用户信息。
//
// 参数:
//   - c: 上下文
//   - username: 要查找的用户名
//
// 返回:
//   - *model.User: 用户对象，如果未找到则为 nil
//   - exist: 是否存在该用户
//   - error: 可能出现的错误
func (r *userRepo) GetUserByUsername(c context.Context, username string) (user *model.User, exist bool, err error) {
	users, err := orm.Q.User.WithContext(c).Where(orm.Q.User.Username.Eq(username)).Find()
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	if len(users) == 0 {
		return nil, false, nil
	}
	return users[0], true, nil
}

// GenerateUniqueUsername 将给定的用户名处理为唯一的有效用户名。
//
// 参数:
//   - c: 上下文
//   - username: 原始用户名
//
// 返回:
//   - string: 处理后的唯一有效用户名
//   - error: 处理过程中可能出现的错误
func (r *userRepo) GenerateUniqueUsername(c context.Context, username string) (string, error) {
	displayName := username

	// 如果用户名包含中文，将其转换为拼音
	if checker.IsChinese(username) {
		username = strings.Join(pinyin.LazyConvert(username, nil), "")
	}

	// 替换空格为连字符，并转换为小写
	username = strings.ReplaceAll(username, " ", "-")
	username = strings.ToLower(username)

	// 检查处理后的用户名是否合法和是否为保留字
	if checker.IsInvalidUsername(username) {
		log.WithContext(c).Errorf("用户名 [%s] -> [%s] 不合法", displayName, username)
		return "", errors.BadRequest(reason.UsernameInvalid)
	}
	if checker.IsReservedUsername(username) {
		log.WithContext(c).Errorf("用户名 [%s] -> [%s] 为保留字", displayName, username)
		return "", errors.BadRequest(reason.UsernameInvalid)
	}

	// 不断生成后缀，直到得到未被使用的用户名
	for suffix := ""; ; {
		_, exist, err := r.GetUserByUsername(c, username+suffix)
		if err != nil {
			return "", err
		} else if !exist {
			return username + suffix, nil
		}

		suffix = random.UsernameSuffix()
	}
}

// IncrRank 增加用户的积分。
//
// 参数:
//   - c: 上下文
//   - userId: 用户 ID
//   - currentRank: 当前排名值
//   - deltaRank: 排名增加值
//
// 返回: 可能返回的错误，如果没有错误则返回 nil
func (r *userRepo) IncrRank(c context.Context, userId, currentRank, deltaRank int64) error {
	// 如果积分增加值为 0，直接返回
	if deltaRank == 0 {
		return nil
	}

	// 如果积分加值为负且会导致排名小于 1，调整增加值
	if deltaRank < 0 && currentRank+deltaRank < 1 {
		deltaRank = 1 - currentRank
	}

	// 执行数据库更新操作，如果出错则返回相应错误
	_, err := orm.Q.User.WithContext(c).Where(orm.Q.User.Id.Eq(userId)).UpdateSimple(orm.Q.User.Rank.Add(deltaRank))
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

// UpdateEmailStatus 更新用户的邮箱状态
func (r *userRepo) UpdateEmailStatus(c context.Context, userId int64, emailStatus string) error {
	_, err := orm.Q.User.WithContext(c).Where(orm.Q.User.Id.Eq(userId)).UpdateSimple(orm.Q.User.MailStatus.Value(emailStatus))
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}
