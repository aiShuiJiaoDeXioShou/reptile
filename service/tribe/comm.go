package tribe

import (
	"github.com/gin-gonic/gin"
	"log"
	"reptile/comm/res"
	"reptile/dao/hotdao/pkmhot"
	"reptile/dao/tribedao"
	"reptile/dao/userdao"
	"strconv"
)

type TribeService struct {
	TribeDao *tribedao.TribeDao
	router   *gin.Engine
}

func NewTribeService(router *gin.Engine) *TribeService {
	return &TribeService{
		TribeDao: tribedao.NewTribeDao(),
		router:   router,
	}
}

func (t *TribeService) StartTribeService() {
	tribe := t.router.Group("/tribe/api")
	{
		// 热门榜单
		tribe.GET("/hot/:page", t.GetHotTribe)
		// 推荐榜单
		tribe.GET("/recommend/", t.GetRecommendTribe)
		// 好友列表
		tribe.GET("/friend/", t.GetFriendList)
		// 好友近期动态
		tribe.GET("/friend/recent/:id", t.GetFriendRecent)
		// 关注一个用户
		tribe.POST("/follow/:id", t.FollowUser)
		// 取消关注一个用户
		tribe.DELETE("/follow/:id", t.UnFollowUser)
	}
}

// 获取热门榜单
func (t *TribeService) GetHotTribe(ctx *gin.Context) {
	// 从路径里获取page
	page := ctx.Param("page")
	// 将page转换为int类型
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(200, res.FormatError("page参数错误"))
		return
	}
	// 获取热门榜单
	tribeList, err := t.TribeDao.GetTribeByTypeAndPage(pkmhot.HotType.HOT, pageInt)
	if err != nil {
		log.Println(err)
		ctx.JSON(200, res.FormatError("获取热门榜单失败"))
		return
	}
	// 没有发生异常则获取榜单成功
	ctx.JSON(200, res.Ok(tribeList))
}

// 获取推荐榜单
func (t *TribeService) GetRecommendTribe(ctx *gin.Context) {
	// 从数据库里面获取推荐榜单
	tribeList, err := t.TribeDao.GetTribeByType(pkmhot.HotType.RECOMMEND)
	if err != nil {
		ctx.JSON(200, res.FormatError("获取推荐榜单失败"))
		return
	}
	// 没有发生异常则获取榜单成功
	ctx.JSON(200, res.Ok(tribeList))
}

// 获取好友列表
func (t *TribeService) GetFriendList(ctx *gin.Context) {
	// 从数据库里面获取好友列表
	// 获取当前的用户对象
	value, exists := ctx.Get("current_user")
	if !exists {
		ctx.JSON(200, res.FormatError("评论失败,请先登录!"))
		return
	}
	user := value.(*userdao.User)
	// 获取好友列表
	friendList, err := userdao.GetFriends(user.ID)
	if err != nil {
		ctx.JSON(200, res.FormatError("获取好友列表失败"))
		return
	}
	// 没有发生异常则获取榜单成功
	ctx.JSON(200, res.Ok(friendList))
}

// 获取好友近期动态
func (t *TribeService) GetFriendRecent(ctx *gin.Context) {
	// 从路径里获取id
	id := ctx.Param("id")
	// 将id转化为uint类型
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.JSON(200, res.FormatError("id参数错误"))
		return
	}
	tribeList, err := userdao.GetRecentDynamic(uint(idUint))
	if err != nil {
		ctx.JSON(200, res.FormatError("获取好友近期动态失败"))
		return
	}
	// 没有发生异常则获取榜单成功
	ctx.JSON(200, res.Ok(tribeList))
}

// 关注一个用户
func (t *TribeService) FollowUser(ctx *gin.Context) {
	// 从路径里获取id
	id := ctx.Param("id")
	// 将id转化为uint类型
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.JSON(200, res.FormatError("id参数错误"))
		return
	}
	// 获取当前的用户对象
	value, exists := ctx.Get("current_user")
	if !exists {
		ctx.JSON(200, res.FormatError("关注失败,请先登录!"))
		return
	}
	user := value.(*userdao.User)
	// 根据id查询路径id用户
	user2, err := userdao.FindUserById(uint(idUint))
	if err != nil {
		ctx.JSON(200, res.FormatError("关注失败,用户不存在!"))
		return
	}

	// 关注用户
	err = userdao.Follow(&userdao.Friend{
		SelfID:      user.ID,
		FriendID:    user2.ID,
		SelfName:    user.Name,
		FriendName:  user2.Name,
		FriendImage: user2.HeadPortrait,
	})

	if err != nil {
		ctx.JSON(200, res.FormatError("关注失败"))
		return
	}
	// 没有发生异常则关注成功
	ctx.JSON(200, res.Ok("关注成功"))
}

// 取关
func (t *TribeService) UnFollowUser(ctx *gin.Context) {
	// 从路径里获取id
	id := ctx.Param("id")
	// 获取当前的用户对象
	value, exists := ctx.Get("current_user")
	if !exists {
		ctx.JSON(200, res.FormatError("取消关注失败,请先登录!"))
		return
	}
	user := value.(*userdao.User)
	// 调用dao层方法取消关注用户
	err := userdao.Unfollow(user.ID, id)
	if err != nil {
		ctx.JSON(200, res.FormatError("取消关注失败"))
		return
	}
	// 没有发生异常则取消关注成功
	ctx.JSON(200, res.Ok("取消关注成功"))
}
