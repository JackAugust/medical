package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"medical/sqlaction"
	"strconv"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *ServiceSetup) UploadMed(args []string) (string, error) {
	DB := InitDB()
	var groups, owner, usertype, org, disease, casenumer, policy, mess string
	// owner := sqlaction.GetUserLogin(DB, "select username from login where state='1'") //确定用户名
	rows := DB.QueryRow("select username,usertype from login where state='1'")
	rows.Scan(&owner, &usertype)
	rows = DB.QueryRow("select user_insti, user_disease from user_type where user_id='" + usertype + "'")
	rows.Scan(&org, &disease)
	// TODO: Groups由哪里获得
	groups = "xxx"
	arr := [17]string{}
	// args[0]是subject，1是txt
	arr[0] = groups
	arr[1] = args[0]
	arr[7] = disease
	arr[11] = owner
	arr[12] = org
	policy = ""
	casenumer = sqlaction.GetCaseNumber(arr[:])
	// TODO: 这里数据添加，数据不全，后面问一下
	if !InsertDB(DB, arr[:], casenumer) {
		return "", fmt.Errorf("数据库插入不成功！")
	} else {
		if InsertDB2Insti(DB, casenumer) {
			policy = GeneratePolicy(DB, casenumer)
			fmt.Println(policy)
		}
	}
	eventID := "eventUploadMed"
	resultStr := "Success"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "UploadMedicalRecord", Args: [][]byte{[]byte(casenumer), []byte(args[11]), []byte(args[12]), []byte(resultStr), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return "", err0
	}

	err1 := eventResult(notifier, eventID)
	fmt.Println(err1)
	if err1 != nil {
		return "", err1
	}

	mess = string(respone.TransactionID)[0:6] + " : " + casenumer + "-policy = " + policy
	//return policy, nil
	//return string(respone.TransactionID), nil
	return mess, nil
}

func (t *ServiceSetup) QueryAllMed() ([]TableRow, error) {
	DB := InitDB()
	var owner string
	//TODO:确定用户名，这里需要每个函数都调用么？有全局变量吗？
	rows := DB.QueryRow("select username from login where state='1'")
	// rows = DB.QueryRow("select username from login where state='1'")
	rows.Scan(&owner)
	fmt.Println("owner is ", owner)
	SQLString1 := "select _SubjectMark from base_info where _Researcher='" + owner + "'"
	subjectMark_list := queryDB(DB, SQLString1)
	// fmt.Println("subjectMark is:", subjectMark_list)
	SQLString2 := "select _CaseNumber from base_info where _Researcher='" + owner + "'"
	caseNumber_list := queryDB(DB, SQLString2)

	// fmt.Println("caseNumber is:", caseNumber_list)
	// SQLString3 := "select _CaseNumber from base_info where _Researcher='" + owner + "'"
	// intro := queryDB(DB, SQLString2)
	// fmt.Println("caseNumber is:", caseNumber_list)
	// TODO：这里需要上链吗？
	// mess_map := make(map[int]interface{})
	// var firstColumn []string
	// var secondColumn []string
	var tabledata []TableRow
	for i := 0; i < len(subjectMark_list); i++ {

		var tablerow TableRow
		// 第一列序号
		tablerow.FirstColumn = strconv.Itoa(i + 1)
		// 第二列subjectmark
		tablerow.SecondColumn = subjectMark_list[i]
		// 第三列caseNumber
		tablerow.ThirdColumn = caseNumber_list[i]
		// 第四列返回结果
		tablerow.FourthColumn = "成功"
		// 第五列备注
		tablerow.FifthColumn = "无"
		// 第六列策略已生成
		tablerow.SixthColumn = "已生成"
		tabledata = append(tabledata, tablerow)
	}
	return tabledata, nil

}

// func (t *ServiceSetup) AllData(user string) (string, error) {

// }
func (t *ServiceSetup) OperateMed(args []string) ([]byte, error) {
	if len(args) != 4 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
	}
	DB := InitDB()
	casenumer := args[0]
	if !CheckAction(DB, casenumer, "r") {
		return nil, fmt.Errorf("权限不足，无法操作")
	}
	eventID := "eventAccessMed"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "AccessMedicalRecord", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return respone.Payload, err0
	}

	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return []byte{0x00}, err1
	}

	mp := SelectDBSingle(DB, args)
	if mp == nil {
		return []byte{0x00}, fmt.Errorf("数据库查询不成功！")
	}
	m := *mp
	b, err2 := json.Marshal(m)
	if err2 != nil {
		return []byte{0x00}, err2
	}
	return b, nil
}

func (t *ServiceSetup) DeleteMed(args []string) (string, error) {
	if len(args) != 4 {
		return "", fmt.Errorf("给定的参数个数不符合要求！")
	}
	DB := InitDB()
	// DB := InitDB()
	casenumer := args[0]
	if !CheckAction(DB, casenumer, "d") {
		return "", fmt.Errorf("权限不足，无法操作")
	}
	if !DeleteDB(DB, args) {
		return "", fmt.Errorf("删除数据不成功！")
	}
	eventID := "eventDeleteMed"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "DeleteMedicalRecord", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return "", err0
	}

	err1 := eventResult(notifier, eventID)
	fmt.Println(err1)
	if err1 != nil {
		return "", err1
	}
	return string(respone.TransactionID), nil
}

func (t *ServiceSetup) UpdateMed(args []string) (string, error) {
	if len(args) != 17 {
		return "", fmt.Errorf("给定的参数个数不符合要求！")
	}
	DB := InitDB()
	casenumer := args[0]
	if !CheckAction(DB, casenumer, "w") {
		return "", fmt.Errorf("权限不足，无法操作")
	}
	if !UpdateDB(DB, args) {
		return "", fmt.Errorf("数据库修改不成功！")
	}
	eventID := "eventUpdateMed"
	resultStr := "Success"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "UpdateMedicalRecord", Args: [][]byte{[]byte(args[4]), []byte(args[11]), []byte(args[12]), []byte(resultStr), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return "", err0
	}

	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return "", err1
	}

	return string(respone.TransactionID), nil
}

func (t *ServiceSetup) UserLogin(username string, password string) (bool, error) {
	DB := InitDB()
	SQLString := "select username from login"
	Username := make(map[int]string)
	Username = queryDB(DB, SQLString)

	for _, user := range Username {
		if user == username {
			var str string
			SQLString2 := "select password from login where username='" + user + "'"
			err := DB.QueryRow(SQLString2).Scan(&str)
			if err != sql.ErrNoRows && str == password {
				SQLString3 := "UPDATE login SET state= '1' where username='" + user + "'"
				_, err := DB.Exec(SQLString3)
				if err != nil {
					return false, err
				}
				return true, nil
			}
		}
	}
	return false, nil
}

func (t *ServiceSetup) UserLoginInfo() (map[int]string, error) {
	Result := make(map[int]string)
	DB := InitDB()
	var str string
	SQLString1 := "select usertype from login where state ='1'"
	err := DB.QueryRow(SQLString1).Scan(&str)
	if err != sql.ErrNoRows {
		SQLString2 := "select * from user_type where user_id ='" + str + "'"
		Result = queryDB(DB, SQLString2)
	}
	return Result, nil
}

func (t *ServiceSetup) UserLoginOut() (bool, error) {
	DB := InitDB()
	Userinfo := sqlaction.GetUserLogin(DB, "select username from login where state='1'") //确定用户名
	SQLString3 := "UPDATE login SET state= '0' where username='" + Userinfo + "'"
	_, err := DB.Exec(SQLString3)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (t *ServiceSetup) AuditAll(args []string) ([]byte, error) {
	if len(args) != 3 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
	}
	eventID := "eventAuditAll"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "auditForAllLogs", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return []byte{0x00}, err0
	}
	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return []byte{0x00}, err1
	}
	return respone.Payload, nil
}
func (t *ServiceSetup) AuditTimeRange(args []string) ([]byte, error) {
	if len(args) != 5 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
	}
	eventID := "eventAuditTimeRange"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "auditForTimeRange", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return []byte{0x00}, err0
	}

	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return []byte{0x00}, err1
	}

	return respone.Payload, nil
}
func (t *ServiceSetup) AuditUser(args []string) ([]byte, error) {
	if len(args) != 4 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
	}
	eventID := "eventAuditUser"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "auditForUser", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return []byte{0x00}, err0
	}

	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return []byte{0x00}, err1
	}

	return respone.Payload, nil
}
func (t *ServiceSetup) AuditOrganisation(args []string) ([]byte, error) {
	if len(args) != 4 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
	}
	eventID := "eventAuditOrganisation"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "auditForOrganisation", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return []byte{0x00}, err0
	}

	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return []byte{0x00}, err1
	}

	return respone.Payload, nil
}
func (t *ServiceSetup) AuditMedicalRecord(args []string) ([]byte, error) {
	if len(args) != 4 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
	}
	eventID := "eventAuditMedicalRecord"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "auditForMedicalRecord", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return []byte{0x00}, err0
	}

	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return []byte{0x00}, err1
	}

	return respone.Payload, nil
}
func (t *ServiceSetup) AuditOriginalAuthor(args []string) ([]byte, error) {
	if len(args) != 4 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
	}
	eventID := "eventAuditOriginalAuthor"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "auditForOriginalAuthor", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return []byte{0x00}, err0
	}

	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return []byte{0x00}, err1
	}

	return respone.Payload, nil
}
func (t *ServiceSetup) AuditPatient(args []string) ([]byte, error) {
	if len(args) != 4 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
	}
	eventID := "eventAuditPatient"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "auditForPatient", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return []byte{0x00}, err0
	}

	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return []byte{0x00}, err1
	}

	return respone.Payload, nil
}
