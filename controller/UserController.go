package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"kaiyuan10nian/common"
	"kaiyuan10nian/model"
	"kaiyuan10nian/response"
	"math/rand"
	"net/http"
	"time"
)

func Register (ctx *gin.Context) {
	DB := common.GetDB()
	//获取参数
	name := ctx.PostForm("name")
	mobile := ctx.PostForm("mobile")
	password := ctx.PostForm("password")
	//数据验证
	if len(mobile) != 11 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"密码不能少于6位")
		return
	}
	if len(name) == 0 {
		name = RandomString(10)
	}

	//判断手机号是否存在
	if isTelephoneExist(DB,mobile){
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"用户已经存在")
		return
	}
	//创建用户
	hasedPassword,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx,http.StatusUnprocessableEntity,500,nil,"加密错误")
		return
	}
	newUser := model.User{
		Name: name,
		Mobile: mobile,
		Password: string(hasedPassword),
	}
	DB.Create(&newUser)
	//返回结果

	response.Success(ctx,nil,"注册成功")
}
//判断手机号是否存在
func isTelephoneExist(db *gorm.DB, mobile string) bool {
	var user model.User
	db.Where("mobile = ?",mobile).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
//生成随机10个字符
func RandomString(n int) string {
	var letters = []byte("qwertyuioplkjhgfdsazxcvbnmMNBVCXZASDFGHJKLPOIUYTREWQ")
	result := make([]byte,n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}