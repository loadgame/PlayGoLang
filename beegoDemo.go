// beegoDemo
package main

//var
//  x: TStream;
//begin
//  x := TMemoryStream.Create;
//  Memo1.Lines.SaveToStream(x);
//  IdHTTP1.Get('http://127.0.0.1:8080');
//  IdHTTP1.Post('http://127.0.0.1:8080', Memo1.Lines);
//  IdHTTP1.Put('http://127.0.0.1:8080', x);
//  IdHTTP1.Delete('http://127.0.0.1:8080');
//  x.Free;
//end;

import (
	"fmt"
	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) Get() {
	this.Ctx.WriteString("Hello User ")
	this.Ctx.WriteString(this.Ctx.Params[":username"])

}

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	mystruct := Mystruct{Id: 2, Msg: "[get ok]"}
	this.Data["json"] = &mystruct
	this.ServeJson()
	fmt.Print(mystruct)
	print(" get:")
}

type Mystruct struct {
	Id  int
	Msg string
}

func (this *MainController) Delete() {

	mystruct := Mystruct{Id: 2, Msg: "[Delete ok]"}
	this.Data["json"] = &mystruct
	this.ServeJson()
	fmt.Print(mystruct)
	print(" delete:")

}
func (this *MainController) Put() {
	jsoninfo := this.GetString("jsoninfo")

	this.Ctx.WriteString("hello world")
	print(" PUT:")
	print(jsoninfo)
}
func (this *MainController) Post() {
	jsoninfo := this.GetString("jsoninfo")
	this.Ctx.WriteString("hello world")
	print(" Post:")
	print(jsoninfo)
}
func main() {
	//beego.SetStaticPath("/images", "images") //静态路径
	beego.Router("/", &MainController{})
	beego.Router("/user/:username([\\w]+)", &UserController{})
	beego.Run()
}
