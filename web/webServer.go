package web

import (
	"fmt"
	"medical/web/controller"
	"net/http"
)

func WebStart(app controller.Application) {

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 指定路由信息(匹配请求)
	http.HandleFunc("/", app.LoginView)

	// 登录
	http.HandleFunc("/login", app.Login)
	http.HandleFunc("/loginout", app.LoginOut)
	// 注册
	http.HandleFunc("/register", app.Register)
	// 身份验证
	http.HandleFunc("/verify", app.Verify)
	// 机构添加
	http.HandleFunc("/addinstitution", app.Addinstitution)
	// 忘记密码
	http.HandleFunc("/forgotpassword", app.Forgotpassword)
	// 更改密码
	http.HandleFunc("/changepassword", app.Changepassword)
	// 主页面
	http.HandleFunc("/index", app.Index)

	// 菜单栏跳转
	// 00-简单搜索展示 - 对应页面为 04-搜索结果展示 显示页面为: SearchDisplay.html
	http.HandleFunc("/简单搜索展示", app.SimpleSearch)
	// 00-高级搜索展示 显示页面为: AdvancedSearch.html
	http.HandleFunc("/高级搜索展示", app.AdvancedSearch)
	// 01-队列信息展示 显示页面为: QueueDisplay.html
	http.HandleFunc("/队列信息展示", app.QueueDisplay)
	// 01-区块信息展示 显示页面为: BlockDisplay.html
	http.HandleFunc("/区块信息展示", app.BlockDisplay)
	// 01-本地存储详情 显示页面为: LocalStorage.html
	http.HandleFunc("/本地存储详情", app.LocalStorage)

	// 02-医疗数据管理 显示页面为: MedicalDataManagement.html
	// 用原来版本  app.ManageMed by monk 02-17
	http.HandleFunc("/医疗数据管理", app.ManageMed)
	// http.HandleFunc("/医疗数据管理", app.MedicalDataManagement)

	// 02-访问控制管理 显示页面为: AccessControlManagement.html
	http.HandleFunc("/访问控制管理", app.AccessControlManagement)
	// 02-数据加密共享 显示页面为: EncryDataShared.html
	http.HandleFunc("/数据加密共享", app.EncryDataShared)
	// 03-医疗数据溯源 显示页面为: MedicalDataTraceability.html
	http.HandleFunc("/医疗数据溯源", app.MedicalDataTraceability)
	// 03-医疗数据审计 显示页面为: MedicalDataAudit.html
	http.HandleFunc("/医疗数据审计", app.MedicalDataAudit)
	// 04-搜索结果展示 显示页面为: SearchDisplay.html
	http.HandleFunc("/搜索结果展示", app.SearchDisplay)
	// 04-访问记录展示 显示页面为: AccessRecordDisplay.html
	http.HandleFunc("/访问记录展示", app.AccessRecordDisplay)
	// 04-操作记录展示 显示页面为: OperationRecordDisplay.html
	http.HandleFunc("/操作记录展示", app.OperationRecordDisplay)
	// 05-用户信息更正 显示页面为: ChangeUserInfo.html
	http.HandleFunc("/用户信息更正", app.ChangeUserInfo)
	// 05-用户信息验证 显示页面为: VerifyUserInfo.html
	http.HandleFunc("/用户信息验证", app.VerifyUserInfo)
	// 05-用户信息展示 显示页面为: DisplayUserInfo.html
	http.HandleFunc("/用户信息展示", app.DisplayUserInfo)

	// 02访问策略生成.html 其实是更新
	http.HandleFunc("/访问策略生成", app.UpdatePolicy)

	// 02医疗数据上传.html
	http.HandleFunc("/uploadMed", app.UploadMed)

	http.HandleFunc("/operateMed", app.OperateMed)
	http.HandleFunc("/auditMed", app.AuditMed)

	http.HandleFunc("/accessMed", app.AccessMed)
	http.HandleFunc("/accessMedHistory", app.AccessMedHistory)
	http.HandleFunc("/accessMedResult", app.AccessMedResult)
	http.HandleFunc("/updateMed", app.UpdateMed)
	http.HandleFunc("/deleteMed", app.DeleteMed)
	http.HandleFunc("/medicalDataTrace", app.MedicalDataTrace)

	http.HandleFunc("/auditResult", app.AuditResult)
	http.HandleFunc("/auditReportResult", app.AuditReportResult)
	http.HandleFunc("/03医疗数据审计", app.AuditReportByTimeRangeAndOrg)
	http.HandleFunc("/auditAllRecords", app.AuditAllRecords)
	http.HandleFunc("/auditTimeRangeStartEnd", app.AuditTimeRangeStartEnd)
	http.HandleFunc("/auditByUser", app.AuditByUser)
	http.HandleFunc("/auditByOrganisation", app.AuditByOrganisation)
	http.HandleFunc("/auditByMedicalRecord", app.AuditByMedicalRecord)
	http.HandleFunc("/auditByOriginalAuthor", app.AuditByOriginalAuthor)
	http.HandleFunc("/auditByPatient", app.AuditByPatient)

	// 数据上传.html
	http.HandleFunc("/dataUpload", app.DataUpload)

    //主页面中选择搜索方式
	http.HandleFunc("/search_method", app.Choose_search_method)

	fmt.Println("启动Web服务, 监听端口号为: 8088")
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		fmt.Printf("Web服务启动失败: %v", err)
	}

}
