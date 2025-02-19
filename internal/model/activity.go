package model

const (
	ActivityTypeUserActivated = ConfigKeyUserActivated // 用户激活
)

type Activity struct {
	Id               int64  // id
	CreatedAt        int64  // 创建时间
	UpdatedAt        int64  // 更新时间
	CancelledAt      int64  // 取消时间
	UserId           int64  // 用户 id
	TriggerUserId    int64  // 触发用户 id
	ObjectId         int64  // 对象 id
	OriginalObjectId int64  // 原始对象 id
	RevisionId       int64  // 修订 id
	Rank             int64  // 排序
	Type             string // 类型
}

func (Activity) TableName() string {
	return "activity"
}
