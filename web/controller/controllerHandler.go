package controller

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"medical/abac"
	"medical/service"
	"net/http"
	"path"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

var cuser User
var data Data

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "login.html", nil)
}

// index
func (app *Application) Index(w http.ResponseWriter, r *http.Request) {
	loginName := r.FormValue("loginName")
	password := r.FormValue("password")
	fmt.Println("the loginname is ", loginName, " and the password is ", password)
	result, _ := app.Setup.UserLogin(loginName, password)
	var flag bool
	flag = result
	data.Flag = false
	data.CurrentUser.Identity = ""
	data.CurrentUser.Password = ""
	if flag {
		// 登录成功
		// TODO: 这里路由有问题，改一下。

		// 添加对用户信息的展示 BY Jack 20230218
		result, _ := app.Setup.UserLoginInfo()
		data.CurrentUser.LoginName = result[0]
		data.CurrentUser.Identity = result[1]

		ShowView(w, r, "index.html", data)
		// app.Index(w, r)
	} else {
		// 登录失败
		data.Flag = true
		data.CurrentUser.LoginName = loginName
		ShowView(w, r, "login.html", data)
		// app.LoginView(w, r)
	}
}

// 用户登录
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------调用controllerhandle login-----------------")
	loginName := r.FormValue("loginName")
	password := r.FormValue("password")
	fmt.Println("the loginname is ", loginName, " and the password is ", password)
	result, _ := app.Setup.UserLogin(loginName, password)
	var flag bool
	flag = result
	data.Flag = false
	data.CurrentUser.Identity = ""
	data.CurrentUser.Password = ""
	if flag {
		// 登录成功
		// TODO: 这里路由有问题，改一下。

		// 添加对用户信息的展示 BY Jack 20230218
		result, _ := app.Setup.UserLoginInfo()
		data.CurrentUser.LoginName = result[0]
		data.CurrentUser.Identity = result[1]

		ShowView(w, r, "index.html", data)
		// app.Index(w, r)
	} else {
		// 登录失败
		data.Flag = true
		data.CurrentUser.LoginName = loginName
		ShowView(w, r, "login.html", data)
		// app.LoginView(w, r)
	}
}

// 用户登出
func (app *Application) LoginOut(w http.ResponseWriter, r *http.Request) {
	cuser = User{}
	result, _ := app.Setup.UserLoginOut()
	if result {
		ShowView(w, r, "login.html", nil)
	}
}

// 忘记密码? Forgotpassword
func (app *Application) Forgotpassword(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "forgotpassword.html", nil)
}

// 更改密码 Forgotpassword
func (app *Application) Changepassword(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "changepassword.html", nil)
}

// 注册 Register
func (app *Application) Register(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "register.html", nil)
}

// 身份验证 Register
func (app *Application) Verify(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "verify.html", nil)
}

// 添加机构 Register
func (app *Application) Addinstitution(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "addinstitution.html", nil)
}

/*
	添加菜单栏跳转
	By Jack
	02-17
*/
// 00-简单搜索展示 SimpleSearch
func (app *Application) SimpleSearch(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "menu/SearchDisplay.html", data)
}

// 00-高级搜索展示 显示页面为: AdvancedSearch.html
func (app *Application) AdvancedSearch(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "menu/AdvancedSearch.html", data)
}

// 01-队列信息展示 显示页面为: QueueDisplay.html
func (app *Application) QueueDisplay(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "menu/QueueDisplay.html", data)
}

// 01-区块信息展示 显示页面为: BlockDisplay.html
func (app *Application) BlockDisplay(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "menu/BlockDisplay.html", data)
}

// 01-本地存储详情 显示页面为: LocalStorage.html
func (app *Application) LocalStorage(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "menu/LocalStorage.html", data)
}

// 02-医疗数据上传 显示页面为:02医疗数据上传.html
func (app *Application) UploadMed(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------调用controllerhandle UploadMed-----------------")
	data.CurrentUser = cuser
	data.Flag = true
	data.Msg = ""
	// TODO：这里要用go解析文件然后将文件内容存到数据库中，然后在baseinfo里加一个数据本体，之后按格式插入即可
	// TODO：这里还没做，暂时先空着------>在做ing
	// TODO：subject目前由后端生成为当前时间戳
	// subjectmark := strconv.FormatInt(time.Now().Unix(), 10)
	// datafiles := r.FormValue("check")
	// fmt.Println("datafiles is set to: ", datafiles)
	if r.Method == "GET" {
		ShowView(w, r, "02医疗数据上传.html", data)
		fmt.Println("----zhanshiyemian")
	} else if r.Method == "POST" {
		// ShowView(w, r, "02医疗数据上传.html", data)

		fmt.Println("解析文件")
		icon, file, err := r.FormFile("datafiles")
		if err != nil {
			//返回错误
			// c.JSON(500,gin.H{"error":"get file info fail "+err.Error()})
			fmt.Println("err is ", err.Error())
		}
		defer icon.Close()
		fmt.Println("上传文件名:", file.Filename)
		//限制上传文件类型
		var FileAllow map[string]bool = map[string]bool{
			".xlsx": true,
			".xls":  true,
			".csv":  true,
		}
		ext := strings.ToLower(path.Ext(file.Filename))
		if _, ok := FileAllow[ext]; !ok {
			fmt.Println("文件后缀名不符合上传要求")
			return
		}
		// 读文件内容
		fileContent, _ := file.Open()
		// TODO: 这里仅用csv测试了，后期测一下exl
		r1 := csv.NewReader(fileContent)
		content, err := r1.ReadAll()
		if err != nil {
			log.Fatalf("can not readall, err is %+v", err)
		}
		for _, row := range content[1:][:] {
			fmt.Println(row)
			// info, err := app.Setup.UploadMed(arr[:])
			info, err := app.Setup.UploadMed(row[:])
			fmt.Println("info is ", info)

			transactionID := strings.Split(info, "-")[0]
			policy := strings.Split(info, "=")[1]
			fmt.Println("policy is ", policy)

			if err != nil {
				data.Msg = err.Error()
			} else {
				var p abac.Policy
				err = json.Unmarshal([]byte(policy), &p)
				data.Msg = "信息添加成功:" + transactionID
				data.Policy = p
			}
			fmt.Println("上传数据后生成的策略为：", data.Policy)
		}
		app.DataUpload(w, r)

		// byteContainer, err := ioutil.ReadAll(fileContent)
		// fmt.Println("接收到内容为=====", byteContainer)
		// mapContainer := make(map[string]interface{})
		// err = json.Unmarshal(byteContainer, &mapContainer)
		// fmt.Println("mapContainer is ", mapContainer)
		// stringContainer := string(byteContainer)
		// fmt.Println("stringContainer is ", stringContainer)
		// container := strings.Split(stringContainer, "")
		// for i := 0; i < len(container); i++ {
		// 	// fmt.Println("i is ", i)
		// 	// fmt.Println("and container is ", container[i])
		// 	// arr := strings.Split(container[i], ",")
		// 	// arr := [17]string{}
		// 	// fmt.Println("datafiles is ", arr)
		// 	// if arr[1] != "" {
		// 	// 	// info, err := app.Setup.UploadMed(arr[:])
		// 	// 	fmt.Println("info is ", info)
		// 	// 	transactionID := strings.Split(info, "-")[0]
		// 	// 	policy := strings.Split(info, "=")[1]
		// 	// 	fmt.Println("policy is ", policy)

		// 	// 	if err != nil {
		// 	// 		data.Msg = err.Error()
		// 	// 	} else {
		// 	// 		var p abac.Policy
		// 	// 		err = json.Unmarshal([]byte(policy), &p)
		// 	// 		data.Msg = "信息添加成功:" + transactionID
		// 	// 		data.Policy = p
		// 	// 	}
		// 	// 	fmt.Println("上传数据后生成的策略为：", data.Policy)
		// 	// }
		// }
		// app.DataUpload(w, r)

	}
	// ShowView(w, r, "02医疗数据上传.html", data)

}

// 02-数据上传 显示页面为:数据上传.html
func (app *Application) DataUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("进入数据上传页面")
	data.CurrentUser = cuser
	fmt.Println("数据为：", data.Msg)
	ShowView(w, r, "数据上传.html", data)
}

// 02-医疗数据管理 显示页面为: 02医疗数据管理.html
func (app *Application) ManageMed(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------调用controllerhandle ManageMed-----------------")
	// TODO: 如果用户未登录，则跳转至登陆界面
	// if cuser.LoginName == "" {
	// 	ShowView(w, r, "login.html", nil)
	// 	return
	// }
	data.Flag = true
	data.Msg = ""
	// tabledata, err := app.Setup.QueryAllMed()
	tabledata, err := app.Setup.QueryAllMed()

	// tabledata_bytes, _ := json.Marshal(tabledata)
	// tabledata_str := string(tabledata_bytes)
	if err != nil {
		data.Msg = err.Error()
	} else {
		// fmt.Println("info is ", tabledata_str)
		data.Table = tabledata
	}
	result, _ := app.Setup.UserLoginInfo()
	data.CurrentUser.LoginName = result[0]
	data.CurrentUser.Identity = result[1]
	ShowView(w, r, "02医疗数据管理.html", data)
}

// 02-访问控制管理 显示页面为: AccessControlManagement.html
func (app *Application) AccessControlManagement(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------调用AccessControlManagement-----------------")
	// TODO: 如果用户未登录，则跳转至登陆界面
	result, _ := app.Setup.UserLoginInfo()
	data.CurrentUser.LoginName = result[0]
	data.CurrentUser.Identity = result[1]
	data.Flag = true

	// 获取表格数据
	tabledata, err := app.Setup.QueryAllMed()
	if err != nil {
		data.Msg = err.Error()
	} else {
		// fmt.Println("info is ", tabledata_str)
		data.Table = tabledata
	}

	dataId := r.FormValue("dataId")
	dataName := r.FormValue("dataName")
	fmt.Println("dataId is ", dataId, "and the dataName is ", dataName)
	if dataId != "" {
		data.Msg = dataId + "/" + dataName
		app.UpdatePolicy(w, r)
	} else {
		fmt.Println("---------------刷新页面-----------------")
		ShowView(w, r, "menu/AccessControlManagement.html", data)

	}
}

// 02-访问策略生成 显示页面为02访问策略生成.html
func (app *Application) UpdatePolicy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------调用controllerhandle UpdatePolicy-----------------")
	result, _ := app.Setup.UserLoginInfo()
	data.CurrentUser.LoginName = result[0]
	data.CurrentUser.Identity = result[1]
	data.Flag = true
	fmt.Println("data.Msg is ", data.Msg)
	if data.Msg != "" {
		dataId := strings.Split(strings.Split(data.Msg, "/")[0], ",")
		dataName := strings.Split(strings.Split(data.Msg, "/")[1], ",")
		fmt.Println("dataId is ", dataId, "and the dataName is ", dataName)
	}

	arr := [18]string{r.FormValue("startTime"), r.FormValue("endTime"), r.FormValue("adminRead"), r.FormValue("adminWrite"), r.FormValue("adminDelete"), r.FormValue("u1Read"), r.FormValue("u1Write"), r.FormValue("u1Delete"), r.FormValue("u2Read"), r.FormValue("u2Write"), r.FormValue("u2Delete"), r.FormValue("u3Read"), r.FormValue("u3Write"), r.FormValue("u3Delete"), r.FormValue("firstOrg"), r.FormValue("secondOrg"), r.FormValue("thirdOrg"), r.FormValue("forthOrg")}
	fmt.Println(" the arr is ", arr)
	ShowView(w, r, "02访问策略生成.html", data)
	// fmt.Println("------------02访问策略生成.html--------------------")

}

// 02-数据加密共享 显示页面为: EncryDataShared.html
func (app *Application) EncryDataShared(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "menu/EncryDataShared.html", data)
}

// 03-医疗数据溯源 显示页面为: MedicalDataTraceability.html
func (app *Application) MedicalDataTraceability(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "menu/MedicalDataTraceability.html", data)
}

// 03-医疗数据审计 显示页面为: MedicalDataAudit.html
func (app *Application) MedicalDataAudit(w http.ResponseWriter, r *http.Request) {
	result, _ := app.Setup.UserLoginInfo()
	data.CurrentUser.LoginName = result[0]
	data.CurrentUser.Identity = result[1]
	data.Msg = ""
	data.Flag = false
	data.History = false
	data.AuditString = "AuditReport"
	ShowView(w, r, "menu/MedicalDataAudit.html", data)
}

// 04-搜索结果展示 显示页面为: SearchDisplay.html
func (app *Application) SearchDisplay(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "menu/SearchDisplay.html", data)
}

// 04-访问记录展示 显示页面为: AccessRecordDisplay.html
func (app *Application) AccessRecordDisplay(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "menu/AccessRecordDisplay.html", data)
}

// 04-操作记录展示 显示页面为: OperationRecordDisplay.html
func (app *Application) OperationRecordDisplay(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "menu/OperationRecordDisplay.html", data)
}

// 05-用户信息更正 显示页面为: ChangeUserInfo.html
func (app *Application) ChangeUserInfo(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "menu/ChangeUserInfo.html", data)
}

// 05-用户信息验证 显示页面为: VerifyUserInfo.html
func (app *Application) VerifyUserInfo(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "menu/VerifyUserInfo.html", data)
}

// 05-用户信息展示 显示页面为: DisplayUserInfo.html
func (app *Application) DisplayUserInfo(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "menu/DisplayUserInfo.html", data)
}

// END

// 暂时弃用 By Jack 02-17

func (app *Application) OperateMed(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	ShowView(w, r, "operateMed.html", data)
}

func (app *Application) MedicalDataTrace(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = true
	ShowView(w, r, "03医疗数据溯源.html", data)
}

func (app *Application) AccessMed(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	ShowView(w, r, "accessMed.html", data)
}

func (app *Application) AccessMedResult(w http.ResponseWriter, r *http.Request) {
	arr := [4]string{r.FormValue("operationRecordID"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("medicalRecordID")}
	var result []byte
	var err error
	if data.History {
		// result, err = app.Setup.GetMedHistory(arr[:])
	} else {
		result, err = app.Setup.OperateMed(arr[:])
		fmt.Println("the result is:", result)
	}
	var med = service.MedicalRecord{}
	json.Unmarshal(result, &med)
	data.Med = med
	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}
	if reflect.DeepEqual(med, service.MedicalRecord{}) {
		if data.History {
			ShowView(w, r, "accessMedHistory.html", data)
		} else {
			ShowView(w, r, "accessMed.html", data)
		}
	} else {
		ShowView(w, r, "accessMedResult.html", data)
	}
}

func (app *Application) DeleteMed(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data.CurrentUser = cuser
	data.Flag = true
	data.Msg = ""
	cs := r.Form["cases[]"]
	fmt.Println(cs)
	for _, c := range cs {
		transactionID, err := app.Setup.DeleteMed(c)
		if err != nil {
			data.Msg = err.Error()
		} else {
			data.Msg = "信息删除成功:" + transactionID
		}
	}
	ShowView(w, r, "02医疗数据管理.html", data)
}

func (app *Application) UpdateMed(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Flag = true
	data.Msg = ""
	arr := [17]string{r.FormValue("groups"), r.FormValue("subjectMark"), r.FormValue("name"), r.FormValue("nameInitials"), r.FormValue("caseNumber"), r.FormValue("sex"), r.FormValue("nation"), r.FormValue("diseases"), r.FormValue("medicalHistory"), r.FormValue("nativePlace"), r.FormValue("diagnose"), cuser.LoginName, r.FormValue("organization"), r.FormValue("addition1"), r.FormValue("addition2"), r.FormValue("addition3"), r.FormValue("status")}
	transactionID, err := app.Setup.UpdateMed(arr[:])

	if err != nil {
		data.Msg = err.Error()
	} else {
		data.Msg = "信息修改成功:" + transactionID
	}
	ShowView(w, r, "updateMed.html", data)
}

func (app *Application) AuditMed(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	ShowView(w, r, "auditMed.html", data)
}

func (app *Application) AccessMedHistory(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = true
	ShowView(w, r, "accessMedHistory.html", data)
}

func (app *Application) AuditAllRecords(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	data.AuditString = "AuditAll"
	ShowView(w, r, "auditAllRecords.html", data)
}

func (app *Application) AuditTimeRangeStartEnd(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	data.AuditString = "AuditTimeRangeStartEnd"
	ShowView(w, r, "auditTimeRangeStartEnd.html", data)
}

func (app *Application) AuditByUser(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	data.AuditString = "AuditByUser"
	ShowView(w, r, "auditByUser.html", data)
}

func (app *Application) AuditByPatient(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	data.AuditString = "AuditByPatient"
	ShowView(w, r, "auditByPatient.html", data)
}

func (app *Application) AuditByOrganisation(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	data.AuditString = "AuditByOrganisation"
	ShowView(w, r, "auditByOrganisation.html", data)
}

func (app *Application) AuditByMedicalRecord(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	data.AuditString = "AuditByMedicalRecord"
	ShowView(w, r, "auditByMedicalRecord.html", data)
}

func (app *Application) AuditByOriginalAuthor(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	data.AuditString = "AuditByOriginalAuthor"
	ShowView(w, r, "auditByOriginalAuthor.html", data)
}

func (app *Application) AuditResult(w http.ResponseWriter, r *http.Request) {
	if data.AuditString == "AuditAll" {
		arr := [3]string{r.FormValue("auditRecordID"), cuser.LoginName, r.FormValue("organisationID")}
		result, err := app.Setup.AuditAll(arr[:])
		var ops = service.OperationRecordArr{}
		json.Unmarshal(result, &ops)
		data.Ops = ops
		data.CurrentUser = cuser
		data.Msg = ""
		data.Flag = false
		if err != nil {
			data.Msg = err.Error()
			data.Flag = true
		}
		if reflect.DeepEqual(data.Ops, service.OperationRecordArr{}) {
			ShowView(w, r, "auditAllRecords.html", nil)
		} else {
			ShowView(w, r, "auditResult.html", data)
		}
	} else if data.AuditString == "AuditTimeRangeStartEnd" {
		arr := [5]string{r.FormValue("auditRecordID"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("startTimePoint"), r.FormValue("endTimePoint")}
		result, err := app.Setup.AuditTimeRange(arr[:])
		var ops = service.OperationRecordArr{}
		json.Unmarshal(result, &ops)
		data.Ops = ops
		data.CurrentUser = cuser
		data.Msg = ""
		data.Flag = false
		if err != nil {
			data.Msg = err.Error()
			data.Flag = true
		}
		if reflect.DeepEqual(data.Ops, service.OperationRecordArr{}) {
			ShowView(w, r, "auditTimeRangeStartEnd.html", nil)
		} else {
			ShowView(w, r, "auditResult.html", data)
		}
	} else if data.AuditString == "AuditByMedicalRecord" {
		arr := [4]string{r.FormValue("auditRecordID"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("str")}
		result, err := app.Setup.AuditMedicalRecord(arr[:])
		var ops = service.OperationRecordArr{}
		json.Unmarshal(result, &ops)
		data.Ops = ops
		data.CurrentUser = cuser
		data.Msg = ""
		data.Flag = false
		if err != nil {
			data.Msg = err.Error()
			data.Flag = true
		}
		if reflect.DeepEqual(data.Ops, service.OperationRecordArr{}) {
			ShowView(w, r, "auditByMedicalRecord.html", nil)
		} else {
			ShowView(w, r, "auditResult.html", data)
		}
	} else if data.AuditString == "AuditByUser" {
		arr := [4]string{r.FormValue("auditRecordID"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("str")}
		result, err := app.Setup.AuditUser(arr[:])
		var ops = service.OperationRecordArr{}
		json.Unmarshal(result, &ops)
		data.Ops = ops
		data.CurrentUser = cuser
		data.Msg = ""
		data.Flag = false
		if err != nil {
			data.Msg = err.Error()
			data.Flag = true
		}
		if reflect.DeepEqual(data.Ops, service.OperationRecordArr{}) {
			ShowView(w, r, "auditByUser.html", nil)
		} else {
			ShowView(w, r, "auditResult.html", data)
		}
	} else if data.AuditString == "AuditByPatient" {
		arr := [4]string{r.FormValue("auditRecordID"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("str")}
		result, err := app.Setup.AuditPatient(arr[:])
		var ops = service.OperationRecordArr{}
		json.Unmarshal(result, &ops)
		data.Ops = ops
		data.CurrentUser = cuser
		data.Msg = ""
		data.Flag = false
		if err != nil {
			data.Msg = err.Error()
			data.Flag = true
		}
		if reflect.DeepEqual(data.Ops, service.OperationRecordArr{}) {
			ShowView(w, r, "auditByPatient.html", nil)
		} else {
			ShowView(w, r, "auditResult.html", data)
		}
	} else if data.AuditString == "AuditByOrganisation" {
		arr := [4]string{r.FormValue("auditRecordID"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("str")}
		result, err := app.Setup.AuditOrganisation(arr[:])
		var ops = service.OperationRecordArr{}
		json.Unmarshal(result, &ops)
		data.Ops = ops
		data.CurrentUser = cuser
		data.Msg = ""
		data.Flag = false
		if err != nil {
			data.Msg = err.Error()
			data.Flag = true
		}
		if reflect.DeepEqual(data.Ops, service.OperationRecordArr{}) {
			ShowView(w, r, "auditByOrganisation.html", nil)
		} else {
			ShowView(w, r, "auditResult.html", data)
		}
	} else if data.AuditString == "AuditByOriginalAuthor" {
		arr := [4]string{r.FormValue("auditRecordID"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("str")}
		result, err := app.Setup.AuditOriginalAuthor(arr[:])
		var ops = service.OperationRecordArr{}
		json.Unmarshal(result, &ops)
		data.Ops = ops
		data.CurrentUser = cuser
		data.Msg = ""
		data.Flag = false
		if err != nil {
			data.Msg = err.Error()
			data.Flag = true
		}
		if reflect.DeepEqual(data.Ops, service.OperationRecordArr{}) {
			ShowView(w, r, "auditByOriginalAuthor.html", nil)
		} else {
			ShowView(w, r, "auditResult.html", data)
		}
	} else {
		ShowView(w, r, "auditMed.html", nil)
	}
}

// *新增：调用两次查询函数，返回指定时间段、指定组织的审计报告
func (app *Application) AuditReportResult(w http.ResponseWriter, r *http.Request) {
	if data.AuditString == "AuditReport" {
		arr := []string{r.FormValue("auditRecordID0"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("startTime"), r.FormValue("endTime")}
		result0, err0 := app.Setup.AuditTimeRange(arr[:])
		var ops0 = service.OperationRecordArr{}
		json.Unmarshal(result0, &ops0)
		arr = []string{r.FormValue("auditRecordID0") + "000", cuser.LoginName, r.FormValue("organisationID"), r.FormValue("auditOrg")}
		result1, err1 := app.Setup.AuditOrganisation(arr[:])
		var ops1 = service.OperationRecordArr{}
		json.Unmarshal(result1, &ops1)

		ops := intersection(ops0, ops1)
		data.Ops = ops
		data.CurrentUser = cuser
		data.Msg = ""
		data.Flag = false
		if err0 != nil {
			data.Msg = err0.Error()
			data.Flag = true
		}
		if err1 != nil {
			data.Msg = err1.Error()
			data.Flag = true
		}

		var repo = service.AuditReport{}
		repo.TargetOrg = r.FormValue("auditOrg")

		// 组织操作信息
		total := 0
		fail := 0

		// m0：每个用户的成功操作数
		// m1：每个用户的失败操作数
		m0 := make(map[string]int, 0)
		m1 := make(map[string]int, 0)

		for _, op := range ops.OperationRecord {
			total++
			m0[op.UserID]++
			if !op.IsSuccess {
				fail++
				m1[op.UserID]++
			}
		}

		curFailRate := float64(0)
		maxFailRate := float64(0)

		for user, v := range m0 {
			curFailRate = float64(m1[user]) / float64(v)
			if curFailRate > maxFailRate {
				maxFailRate = curFailRate
				repo.MaxFailRateUser = user
			}
		}
		repo.MaxFailRate = maxFailRate

		repo.TotalOperations = int64(total)
		repo.FailOperations = int64(fail)
		repo.FailRate = float64(fail) / float64(total)

		//动态区间实现
		//filePath := "./web/controller/0.txt"
		//file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0666)
		//if err != nil {
		//fmt.Println("文件打开失败", err)
		//}
		//及时关闭file句柄
		//defer file.Close()
		//读原来文件的内容，并且显示在终端
		//reader := bufio.NewReader(file)

		// Credit[0] Credit[1]：成功率区间
		// Credit[2]：组织信誉值
		//var Credit [3]float64

		//for i := 0; i < 3; i++ {
		//str, err := reader.ReadString('\n')

		//if err == io.EOF {
		//break
		//
		//str = strings.Replace(str, "\r", "", -1)
		//str = strings.Replace(str, "\n", "", -1)
		//sc, err := strconv.ParseFloat(str, 64)
		//if err != nil {
		//fmt.Println("error in string to float64", err)
		//}
		//Credit[i] = sc
		//}
		var (
			userName  string = "root"
			password  string = "root"
			ipAddrees string = "127.0.0.1"
			port      int    = 3306
			dbName    string = "itbtsql"
		)

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?allowNativePasswords=true", userName, password, ipAddrees, port, dbName)
		Db, err := sqlx.Open("mysql", dsn)
		if err != nil {
			fmt.Printf("mysql connect failed, detail is [%v]", err.Error())
		} else {
			fmt.Printf("mysql connect success!\n")
		}
		defer Db.Close()

		var intv0 float64
		var intv1 float64
		var Credit float64
		rows, err := Db.Query("select intv0, intv1, Credit from credit_table where TargetOrg=" + repo.TargetOrg)
		if err != nil {
			fmt.Println("select failed:", err)
		}
		for rows.Next() {
			rows.Scan(&intv0, &intv1, &Credit)
			fmt.Println(intv0, intv1, Credit)
		}
		defer rows.Close()

		// 成功率区间
		var intv [2]float64
		intv[0] = intv0
		intv[1] = intv1
		repo.ReferenceRange = intv
		repo.PreviousCredit = Credit

		//成功率高出区间，则信誉值上升，区间变化
		//成功率低出区间，则信誉值降低，区间变化
		//成功率处于区间里，信誉值不变，区间不变
		if 1-repo.FailRate > intv[1] {
			repo.CreditChange = "上升"
			Credit = (Credit + 1 - repo.FailRate) / 2
			intv0 = math.Min(Credit, 1-repo.FailRate)
			intv1 = math.Max(Credit, 1-repo.FailRate)
		} else if 1-repo.FailRate < intv[0] {
			repo.CreditChange = "下降"
			Credit = (Credit + 1 - repo.FailRate) / 2
			intv0 = math.Min(Credit, 1-repo.FailRate)
			intv1 = math.Max(Credit, 1-repo.FailRate)
		} else {
			repo.CreditChange = "不变"
		}
		repo.CurrentCredit = Credit

		//os.Truncate("./web/controller/0.txt", 0)
		//写入文件时，使用带缓存的 *Writer
		//write := bufio.NewWriter(file)
		//for i := 0; i < 3; i++ {
		//str := strconv.FormatFloat(Credit[i], 'f', 10, 64)
		//write.WriteString(str + "\r\n")
		//}
		//Flush将缓存的文件真正写入到文件中
		//write.Flush()

		sql := "update credit_table set intv0=?, intv1=?, Credit=? where TargetOrg=?"
		result, err := Db.Exec(sql, intv0, intv1, Credit, repo.TargetOrg)
		if err != nil {
			fmt.Println("update failed:", err)
		}
		row, err := result.RowsAffected()
		if err != nil {
			fmt.Println("row failed:", err)
		}
		fmt.Println("update success:", row)

		data.Repo = repo

		if reflect.DeepEqual(ops, service.OperationRecordArr{}) {
			ShowView(w, r, "auditReportResult.html", data)
			//ShowView(w, r, "03医疗数据审计.html", data)
		} else {
			ShowView(w, r, "auditReportResult.html", data)
		}
	} else {
		ShowView(w, r, "auditMed.html", nil)
	}
}

// *新增：同时依照时间段和组织ID审计
func (app *Application) AuditReportByTimeRangeAndOrg(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	data.AuditString = "AuditReport"
	ShowView(w, r, "03医疗数据审计.html", data)
}

// *新增：返回两个OperationRecordArr的交集
func intersection(nums1 service.OperationRecordArr, nums2 service.OperationRecordArr) service.OperationRecordArr {

	m := make(map[string]int, 0)

	for _, v := range nums1.OperationRecord {
		m[v.OperationRecordID] += 1
	}

	count := 0 //记录新数组长度
	for _, v := range nums2.OperationRecord {
		if m[v.OperationRecordID] > 0 {
			m[v.OperationRecordID] = 0
			nums1.OperationRecord[count] = v
			count++
		}
	}

	return service.OperationRecordArr{OperationRecord: nums1.OperationRecord[:count]}

}

// choose search method section
func (app *Application) Choose_search_method(w http.ResponseWriter, r *http.Request) {
	select_method := r.FormValue("search_select")
	search_key := r.FormValue("search_key")
	fmt.Print("搜索方法：")
	fmt.Println(select_method)
	fmt.Print("搜索关键字：")
	fmt.Println(search_key)
	method_chose := Strval(select_method)

	//different search based on the posted method above

	if method_chose == "深度搜索" {

		//数据库中的测试数据转变成二维的数组并显示出来
		data_sql := Check_sqldata()

		//将数据库中的数据变为一棵树
		data_sql_tree := Create_Multi_branch_tree(data_sql)
		fmt.Println("深度搜索相关记录：")
		rsc := data_sql_tree.LeafNodeDFS_key(search_key) //...代表使用语法糖进行打散操作，append的是data_sql_tree.LeafNodeDFS_key("天天")中的元素，不是这个数组本身
		fmt.Print("_CaseNumber：")
		fmt.Println(rsc...)
		fmt.Println("查询的对应的记录：")
		rs_sum2 := Check_data(rsc)
		for _, i := range rs_sum2 {
			fmt.Println(i...)
		}

	} else if method_chose == "广度搜索" {

		//数据库中的测试数据转变成二维的数组并显示出来
		data_sql := Check_sqldata()

		//将数据库中的数据变为一棵树
		data_sql_tree := Create_Multi_branch_tree(data_sql)
		fmt.Println("广度搜索相关记录：")
		rs := data_sql_tree.LeafNodeBFS_key(search_key)
		fmt.Print("_CaseNumber：")
		fmt.Println(rs...)
		fmt.Println("查询的对应的记录：")
		rs_sum1 := Check_data(rs)
		for _, i := range rs_sum1 {
			fmt.Println(i...)
		}

	}

	ShowView(w, r, "index.html", nil)

}
