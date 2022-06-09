/**
 * @Author: wenliangzhang
 * @Description:
 * @File: errno
 * @Version: 1.0.0
 * @Date: 2022/4/22 2:09 PM
 */
package common

const (
	ErrOk      = 0
	ErrUnknown = 109999 // 未知错误

	// 服务调用
	ErrServiceIllegalHost        = 100001 // 非法host
	ErrServiceCallFail           = 100002 // 调用服务失败
	ErrServiceCheckSignAkError   = 100003 // sign ak校验失败
	ErrServiceCheckSignSignError = 100004 // sign校验失败
	ErrServiceCheck              = 100005 // sign时间校验失败
	ErrServicePackConvertFail    = 100006 // 参数解析失败
	ErrServiceMethodNoExist      = 100007 // 方法不存在
	ErrServiceCheckCallChain     = 100008 // 调用链出错

	// 权限相关
	ErrTokenError      = 10001  //Token异常，请重新登录
	ErrIllegalSource   = 100009 // 来源非法
	ErrOnlyInner       = 100010 // 只允许内网访问
	ErrIllegalPerm     = 100011 // 权限不足
	ErrResourceNoExist = 100012 // 资源不存在
	ErrParamError      = 101001 // 参数错误
	ErrCheckFail       = 101002 // 校验失败
	ErrActionControl   = 101003 // 粒度控制
	ErrUiInitFail      = 101004 // ui基类初始化失败
	ErrApiLocation     = 101005 // 接口重定向
	ErrAntiBot         = 101006 // 反爬取
	ErrAntiCaptcha     = 101007 // 二次验证

	// 用户相关
	ErrUserNotLogin  = 102001 // 用户未登录
	ErrUserBlackList = 102002 // 黑名单用户
	ErrUserAntiBlack = 102003 // 反作弊判黑

	//Dao相关
	ErrDaoInstanceNameRequired = 103001 // 缺少实例名
	ErrDaoCreateMethodRequired = 103002 // 缺少创建方法
	ErrDaoInvalidCreateMethod  = 103003 // 创建方法无效
	ErrDaoInvalidInitMethod    = 103004 // 初始化方法无效
	ErrDaoInvalidHookMethod    = 103005 // 钩子方法无效
	ErrDaoInitFail             = 103006 // 初始化失败
	ErrDaoConfError            = 103007 // 配置错误
	ErrDaoStatusError          = 103008 // 状态错误
	ErrDaoInvalidContext       = 103009 // 上下文非法
	ErrDaoCreateFail           = 103010 // 创建失败

	// 数据库相关
	ErrDbConstructFail = 104001  // db名字或表名为空
	ErrDbConnFail      = 104002  // db连接失败
	ErrDbQueryFail     = 104003  // db操作失败
	ErrDbParamError    = 104004  // db参数错误
	ErrRedisCallFail   = 104010  // redis调用失败
	ErrMqConfError     = 104021  // mq没有找到conf
	ErrMqCallFail      = 1040222 // mq操作失败
)
