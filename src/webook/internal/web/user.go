package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"net/http"
	"xiaoweishu/webook/internal/domain"
	"xiaoweishu/webook/internal/service"
)

type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp // 预编译
	svc         *service.UserService
}

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		// 在引入三方库之后，要加一个option regexp.None
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:         svc,
	}
}
func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/users/signup", h.SignUp)
	// 若又很多个相同前缀，可以
	ug := server.Group("/users")
	// POST /users/login
	ug.POST("/login", h.Login)
	// POST /users/profile
	ug.GET("/profile", h.Profile)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignupReq struct {
		Email           string `json: "email"`
		Password        string `json: "password"`
		ConfirmPassword string `json: "confirmPassword"`
	}

	var req SignupReq
	// Bind根据content-type解决数据
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 原本正则匹配方法 -> 进阶预编译
	//isEmail, err := regexp.Match(emailRegexPattern, []byte(req.Email))
	// 预编译 -> go自带正则较弱，这里引用别的库
	//isEmail := h.emailExp.Match([]byte(req.Email))
	isEmail, err := h.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入密码不一致")
	}
	if !isEmail {
		ctx.String(http.StatusOK, "邮箱不正确")
		return
	}

	isPassword, err := h.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码格式不正确，包含数字，特殊字符且长度大于8位")
		return
	}
	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	ctx.String(http.StatusOK, "注册成功")

}

func (h *UserHandler) Login(ctx *gin.Context) {

}

func (h *UserHandler) Profile(ctx *gin.Context) {

}
