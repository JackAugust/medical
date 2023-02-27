package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

const (
	MedicalRecordString   = "medical-record-key"
	OperationRecordString = "operation-record-key"
	AuditRecordString     = "audit-record-key"
)

// 写到链码里，只保留上链字段
type MedicalRecord struct {
	CaseNumber      string `json:"caseNumber"`      //病例号，上链
	Researcher      string `json:"researcher"`      //研究者，上链
	Organization    string `json:"organization"`    //机构，上链
	OperationType   string `json:"operationType"`   //操作类型
	OperationResult string `json:"operationResult"` //操作结果
	GatherTime      string `json:"gatherTime"`      //采集时间-1，上链
}

// 操作记录，记录在区块链上不可被修改
type OperationRecord struct {
	OperationRecordID string `json:"operationRecordId"` //操作记录ID
	UserID            string `json:"userId"`            //操作者用户
	OrganizationID    string `json:"organizationID"`    //用户组织ID
	DataType          string `json:"dataType"`          //记录类型：操作记录operation
	ActionType        string `json:"actionType"`        //操作类型：访问、修改、删除（Access、Update、Delete）
	MedicalRecordID   string `json:"medicalRecordId"`   //涉及到的医疗记录ID
	OriginalAuthorID  string `json:"originalAuthorId"`  //涉及到的医疗记录创建者
	PatientID         string `json:"patientId"`         //医疗记录涉及到的患者ID
	EntryMethod       string `json:"entryMethod"`       //操作记录的输入方式：自动输入
	IsSuccess         bool   `json:"isSuccess"`         //[fxb] 操作是否成功
	Time              string `json:"time"`              //时间
}

type OperationRecordArr struct {
	OperationRecord []OperationRecord
}

// 审计记录，记录在区块链上不可被修改
type AuditRecord struct {
	AuditRecordID  string `json:"auditRecordId"`  //审计记录ID
	UserID         string `json:"userId"`         //操作者用户（此处为管理员）
	OrganizationID string `json:"organizationID"` //用户组织ID
	DataType       string `json:"dataType"`       //记录类型：审计Audit
	EntryMethod    string `json:"entryMethod"`    //审计记录的输入方式：自动输入
	Time           string `json:"time"`           //时间
}

type SmartContract struct {
}

func (t *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println(" ==== Init ====")

	return shim.Success(nil)
}

func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// 获取用户意图
	fun, args := stub.GetFunctionAndParameters()

	if fun == "UploadMedicalRecord" {
		return t.UploadMedicalRecord(stub, args)
	} else if fun == "AccessMedicalRecord" {
		return t.AccessMedicalRecord(stub, args)
	} else if fun == "UpdateMedicalRecord" {
		return t.UpdateMedicalRecord(stub, args)
	} else if fun == "DeleteMedicalRecord" {
		return t.DeleteMedicalRecord(stub, args)
	} else if fun == "auditForAllLogs" {
		return t.auditForAllLogs(stub, args)
	} else if fun == "auditForTimeRange" {
		return t.auditForTimeRange(stub, args)
	} else if fun == "auditForUser" {
		return t.auditForUser(stub, args)
	} else if fun == "auditForOrganisation" {
		return t.auditForOrganisation(stub, args)
	} else if fun == "auditForMedicalRecord" {
		return t.auditForMedicalRecord(stub, args)
	} else if fun == "auditForOriginalAuthor" {
		return t.auditForOriginalAuthor(stub, args)
	} else if fun == "auditForPatient" {
		return t.auditForPatient(stub, args)
	}

	return shim.Error("指定的函数名称错误")

}

// 上传医疗记录
func (t *SmartContract) UploadMedicalRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	r, _ := stub.GetState(args[0])
	if r != nil {
		return shim.Error("医疗记录ID已经存在！")
	}
	m := MedicalRecord{
		CaseNumber:      args[0],
		Researcher:      args[1],
		Organization:    args[2],
		OperationType:   "Medical",
		OperationResult: args[3],
		GatherTime:      time.Now().Format("2006-01-02 15:04:05"),
	}
	b, err1 := json.Marshal(m)
	fmt.Println(err1)
	if err1 != nil {
		return shim.Error(err1.Error())
	}
	err2 := stub.PutState(m.CaseNumber, b)
	fmt.Println(err2)
	if err2 != nil {
		return shim.Error(err2.Error())
	}
	err3 := stub.SetEvent(args[4], []byte{})
	fmt.Println(err3)
	if err3 != nil {
		return shim.Error(err3.Error())
	}
	return shim.Success([]byte("医疗记录添加成功！"))
}

// 查看医疗记录
func (t *SmartContract) AccessMedicalRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	result0, err0 := stub.GetState(args[3])
	if err0 != nil {
		return shim.Error("医疗记录获取失败！")
	}
	if result0 == nil {
		return shim.Error("医疗记录ID获取失败！")
	}
	result1, _ := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{args[0]})
	defer result1.Close()
	if result1.HasNext() {
		return shim.Error("操作记录ID存在！")
	}
	var m MedicalRecord
	err2 := json.Unmarshal(result0, &m)
	if err2 != nil {
		return shim.Error("反序列化医疗记录失败！")
	}
	key, errK := stub.CreateCompositeKey(OperationRecordString, []string{args[0], args[1], args[2], args[3], m.Researcher, ""})
	if errK != nil {
		return shim.Error("组合键创建失败！")
	}
	o := OperationRecord{
		OperationRecordID: args[0],
		UserID:            args[1],
		OrganizationID:    args[2],
		DataType:          "Operation",
		ActionType:        "Access",
		MedicalRecordID:   args[3],
		OriginalAuthorID:  m.Researcher,
		PatientID:         "",
		EntryMethod:       "Auto",
		IsSuccess:         args[2] == m.Organization, //如果组织不一致则判断为False
		Time:              time.Now().Format("2006-01-02 15:04:05"),
	}
	//ActionType为Access
	op, err3 := json.Marshal(o)
	if err3 != nil {
		return shim.Error("操作记录加工失败！")
	}
	err4 := stub.PutState(key, op)
	if err4 != nil {
		return shim.Error("操作失败！")
	}
	err5 := stub.SetEvent(args[4], []byte{})
	if err5 != nil {
		return shim.Error(err5.Error())
	}

	//如果组织不一致则报错
	//if !o.IsSuccess {
	//	return shim.Error("错误操作！")
	//}
	return shim.Success(result0)
}

// 删除医疗记录
func (t *SmartContract) DeleteMedicalRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	err0 := stub.DelState(args[0])
	if err0 != nil {
		return shim.Error("医疗记录删除失败！")
	}
	err1 := stub.SetEvent(args[1], []byte{})
	if err1 != nil {
		return shim.Error(err1.Error())
	}
	return shim.Success([]byte("医疗记录删除成功！"))
}

// 对医疗记录进行修改
func (t *SmartContract) UpdateMedicalRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	result0, err0 := stub.GetState(args[0])
	if err0 != nil {
		return shim.Error("医疗记录获取失败！")
	}
	if result0 == nil {
		return shim.Error("医疗记录ID获取失败！")
	}
	var m MedicalRecord
	err1 := json.Unmarshal(result0, &m)
	if err1 != nil {
		return shim.Error("反序列化医疗记录失败！")
	}
	m.Researcher = args[1]
	m.Organization = args[2]
	m.OperationResult = args[3]
	newMedicalRecordAsBytes, err2 := json.Marshal(m)
	if err2 != nil {
		return shim.Error("医疗记录加工失败！")
	}
	err3 := stub.PutState(m.CaseNumber, newMedicalRecordAsBytes)
	if err3 != nil {
		return shim.Error("医疗记录修改失败！")
	}
	err4 := stub.SetEvent(args[4], []byte{})
	if err4 != nil {
		return shim.Error(err4.Error())
	}
	return shim.Success([]byte("医疗记录修改成功！"))
}

// 获取所有时间的操作记录
func (t *SmartContract) auditForAllLogs(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 4 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	resultsIterator, err0 := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{})
	if err0 != nil {
		return shim.Error("操作记录获取失败！")
	}
	result1, err1 := stub.GetState(args[0])
	if err1 == nil && result1 != nil {
		return shim.Error("审计记录存在！")
	}
	defer resultsIterator.Close()
	results := OperationRecordArr{}
	for resultsIterator.HasNext() {
		queryResult, err2 := resultsIterator.Next()
		if err2 != nil {
			return shim.Error("迭代失败！")
		}
		var o OperationRecord
		err3 := json.Unmarshal(queryResult.Value, &o)
		if err3 != nil {
			return shim.Error("反序列化操作记录失败！")
		}
		if reflect.TypeOf(o).Name() == "OperationRecord" && o.DataType == "Operation" && (o.ActionType == "Access" || o.ActionType == "Change" || o.ActionType == "Delete") {
			results.OperationRecord = append(results.OperationRecord, o)
		}
	}
	a := AuditRecord{
		AuditRecordID:  args[0],
		UserID:         args[1],
		OrganizationID: args[2],
		DataType:       "Audit",
		EntryMethod:    "Auto",
		Time:           time.Now().Format("2006-01-02 15:04:05"),
	}
	op, err4 := json.Marshal(a)
	if err4 != nil {
		return shim.Error("审计记录加工失败！")
	}
	err5 := stub.PutState(args[0], op)
	if err5 != nil {
		return shim.Error("审计失败！")
	}
	err6 := stub.SetEvent(args[3], []byte{})
	if err6 != nil {
		return shim.Error(err6.Error())
	}
	r, err7 := json.Marshal(results)
	if err7 != nil {
		return shim.Error("结果加工失败！")
	}
	return shim.Success(r)
}

// 获取指定开始结束时间的操作记录
func (t *SmartContract) auditForTimeRange(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 6 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	resultsIterator, err0 := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{})
	if err0 != nil {
		return shim.Error("操作记录获取失败！")
	}
	result1, err1 := stub.GetState(args[0])
	if err1 == nil && result1 != nil {
		return shim.Error("审计记录存在！")
	}
	defer resultsIterator.Close()
	results := OperationRecordArr{}
	for resultsIterator.HasNext() {
		queryResult, err2 := resultsIterator.Next()
		if err2 != nil {
			return shim.Error("迭代失败！")
		}
		var o OperationRecord
		err3 := json.Unmarshal(queryResult.Value, &o)
		if err3 != nil {
			return shim.Error("反序列化医疗记录失败！")
		}
		if reflect.TypeOf(o).Name() == "OperationRecord" && o.DataType == "Operation" && (o.ActionType == "Access" || o.ActionType == "Change" || o.ActionType == "Delete") {
			formatTime, errT := time.Parse("2006-01-02 15:04:05", o.Time)
			if errT != nil {
				return shim.Error("记录时间转换失败！")
			}
			startTime, errST := time.Parse("2006-01-02 15:04:05", args[3])
			if errST != nil {
				return shim.Error("开始时间转换失败！")
			}
			endTime, errET := time.Parse("2006-01-02 15:04:05", args[4])
			if errET != nil {
				return shim.Error("结束时间转换失败！")
			}
			if formatTime.Unix() >= startTime.Unix() && formatTime.Unix() <= endTime.Unix() {
				// arg[3]：开始时间
				// arg[4]：结束时间
				results.OperationRecord = append(results.OperationRecord, o)
			}
		}
	}
	a := AuditRecord{
		AuditRecordID:  args[0],
		UserID:         args[1],
		OrganizationID: args[2],
		DataType:       "Audit",
		EntryMethod:    "Auto",
		Time:           time.Now().Format("2006-01-02 15:04:05"),
	}
	op, err4 := json.Marshal(a)
	if err4 != nil {
		return shim.Error("审计记录加工失败！")
	}
	err5 := stub.PutState(args[0], op)
	if err5 != nil {
		return shim.Error("审计失败！")
	}
	err6 := stub.SetEvent(args[5], []byte{})
	if err6 != nil {
		return shim.Error(err6.Error())
	}
	r, err7 := json.Marshal(results)
	if err7 != nil {
		return shim.Error("结果加工失败！")
	}
	return shim.Success(r)
}

// 获取指定用户的操作记录
func (t *SmartContract) auditForUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	resultsIterator, err0 := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{})
	if err0 != nil {
		return shim.Error("操作记录获取失败！")
	}
	result1, err1 := stub.GetState(args[0])
	if err1 == nil && result1 != nil {
		return shim.Error("审计记录存在！")
	}
	defer resultsIterator.Close()
	results := OperationRecordArr{}
	for resultsIterator.HasNext() {
		queryResult, err2 := resultsIterator.Next()
		if err2 != nil {
			return shim.Error("迭代失败！")
		}
		var o OperationRecord
		err3 := json.Unmarshal(queryResult.Value, &o)
		if err3 != nil {
			return shim.Error("反序列化医疗记录失败！")
		}
		if reflect.TypeOf(o).Name() == "OperationRecord" && o.DataType == "Operation" && (o.ActionType == "Access" || o.ActionType == "Change" || o.ActionType == "Delete") {
			if o.UserID == args[3] {
				results.OperationRecord = append(results.OperationRecord, o)
			}
		}
	}
	a := AuditRecord{
		AuditRecordID:  args[0],
		UserID:         args[1],
		OrganizationID: args[2],
		DataType:       "Audit",
		EntryMethod:    "Auto",
		Time:           time.Now().Format("2006-01-02 15:04:05"),
	}
	op, err4 := json.Marshal(a)
	if err4 != nil {
		return shim.Error("审计记录加工失败！")
	}
	err5 := stub.PutState(args[0], op)
	if err5 != nil {
		return shim.Error("审计失败！")
	}
	err6 := stub.SetEvent(args[4], []byte{})
	if err6 != nil {
		return shim.Error(err6.Error())
	}
	r, err7 := json.Marshal(results)
	if err7 != nil {
		return shim.Error("结果加工失败！")
	}
	return shim.Success(r)
}

// 获取关于指定组织的操作记录
func (t *SmartContract) auditForOrganisation(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	resultsIterator, err0 := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{})
	if err0 != nil {
		return shim.Error("操作记录获取失败！")
	}
	result1, err1 := stub.GetState(args[0])
	if err1 == nil && result1 != nil {
		return shim.Error("审计记录存在！")
	}
	defer resultsIterator.Close()
	results := OperationRecordArr{}
	for resultsIterator.HasNext() {
		queryResult, err2 := resultsIterator.Next()
		if err2 != nil {
			return shim.Error("迭代失败！")
		}
		var o OperationRecord
		err3 := json.Unmarshal(queryResult.Value, &o)
		if err3 != nil {
			return shim.Error("反序列化医疗记录失败！")
		}
		if reflect.TypeOf(o).Name() == "OperationRecord" && o.DataType == "Operation" && (o.ActionType == "Access" || o.ActionType == "Change" || o.ActionType == "Delete") {
			if o.OrganizationID == args[3] {
				results.OperationRecord = append(results.OperationRecord, o)
			}
		}
	}
	a := AuditRecord{
		AuditRecordID:  args[0],
		UserID:         args[1],
		OrganizationID: args[2],
		DataType:       "Audit",
		EntryMethod:    "Auto",
		Time:           time.Now().Format("2006-01-02 15:04:05"),
	}
	op, err4 := json.Marshal(a)
	if err4 != nil {
		return shim.Error("审计记录加工失败！")
	}
	err5 := stub.PutState(args[0], op)
	if err5 != nil {
		return shim.Error("审计失败！")
	}
	err6 := stub.SetEvent(args[4], []byte{})
	if err6 != nil {
		return shim.Error(err6.Error())
	}
	r, err7 := json.Marshal(results)
	if err7 != nil {
		return shim.Error("结果加工失败！")
	}
	return shim.Success(r)
}

// 获取特定医疗记录的操作记录
func (t *SmartContract) auditForMedicalRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	resultsIterator, err0 := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{})
	if err0 != nil {
		return shim.Error("操作记录获取失败！")
	}
	result1, err1 := stub.GetState(args[0])
	if err1 == nil && result1 != nil {
		return shim.Error("审计记录存在！")
	}
	defer resultsIterator.Close()
	results := OperationRecordArr{}
	for resultsIterator.HasNext() {
		queryResult, err2 := resultsIterator.Next()
		if err2 != nil {
			return shim.Error("迭代失败！")
		}
		var o OperationRecord
		err3 := json.Unmarshal(queryResult.Value, &o)
		if err3 != nil {
			return shim.Error("反序列化医疗记录失败！")
		}
		if reflect.TypeOf(o).Name() == "OperationRecord" && o.DataType == "Operation" && (o.ActionType == "Access" || o.ActionType == "Change" || o.ActionType == "Delete") {
			if o.MedicalRecordID == args[3] {
				results.OperationRecord = append(results.OperationRecord, o)
			}
		}
	}
	a := AuditRecord{
		AuditRecordID:  args[0],
		UserID:         args[1],
		OrganizationID: args[2],
		DataType:       "Audit",
		EntryMethod:    "Auto",
		Time:           time.Now().Format("2006-01-02 15:04:05"),
	}
	op, err4 := json.Marshal(a)
	if err4 != nil {
		return shim.Error("审计记录加工失败！")
	}
	err5 := stub.PutState(args[0], op)
	if err5 != nil {
		return shim.Error("审计失败！")
	}
	err6 := stub.SetEvent(args[4], []byte{})
	if err6 != nil {
		return shim.Error(err6.Error())
	}
	r, err7 := json.Marshal(results)
	if err7 != nil {
		return shim.Error("结果加工失败！")
	}
	return shim.Success(r)
}

// 获取特定医疗记录创建者的操作记录
func (t *SmartContract) auditForOriginalAuthor(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	resultsIterator, err0 := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{})
	if err0 != nil {
		return shim.Error("操作记录获取失败！")
	}
	result1, err1 := stub.GetState(args[0])
	if err1 == nil && result1 != nil {
		return shim.Error("审计记录存在！")
	}
	defer resultsIterator.Close()
	results := OperationRecordArr{}
	for resultsIterator.HasNext() {
		queryResult, err2 := resultsIterator.Next()
		if err2 != nil {
			return shim.Error("迭代失败！")
		}
		var o OperationRecord
		err3 := json.Unmarshal(queryResult.Value, &o)
		if err3 != nil {
			return shim.Error("反序列化医疗记录失败！")
		}
		if reflect.TypeOf(o).Name() == "OperationRecord" && o.DataType == "Operation" && (o.ActionType == "Access" || o.ActionType == "Change" || o.ActionType == "Delete") {
			if o.OriginalAuthorID == args[3] {
				results.OperationRecord = append(results.OperationRecord, o)
			}
		}
	}
	a := AuditRecord{
		AuditRecordID:  args[0],
		UserID:         args[1],
		OrganizationID: args[2],
		DataType:       "Audit",
		EntryMethod:    "Auto",
		Time:           time.Now().Format("2006-01-02 15:04:05"),
	}
	op, err4 := json.Marshal(a)
	if err4 != nil {
		return shim.Error("审计记录加工失败！")
	}
	err5 := stub.PutState(args[0], op)
	if err5 != nil {
		return shim.Error("审计失败！")
	}
	err6 := stub.SetEvent(args[4], []byte{})
	if err6 != nil {
		return shim.Error(err6.Error())
	}
	r, err7 := json.Marshal(results)
	if err7 != nil {
		return shim.Error("结果加工失败！")
	}
	return shim.Success(r)
}

// 获取指定病人医疗记录的操作记录
func (t *SmartContract) auditForPatient(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	resultsIterator, err0 := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{})
	if err0 != nil {
		return shim.Error("操作记录获取失败！")
	}
	result1, err1 := stub.GetState(args[0])
	if err1 == nil && result1 != nil {
		return shim.Error("审计记录存在！")
	}
	defer resultsIterator.Close()
	results := OperationRecordArr{}
	for resultsIterator.HasNext() {
		queryResult, err2 := resultsIterator.Next()
		if err2 != nil {
			return shim.Error("迭代失败！")
		}
		var o OperationRecord
		err3 := json.Unmarshal(queryResult.Value, &o)
		if err3 != nil {
			return shim.Error("反序列化医疗记录失败！")
		}
		if reflect.TypeOf(o).Name() == "OperationRecord" && o.DataType == "Operation" && (o.ActionType == "Access" || o.ActionType == "Change" || o.ActionType == "Delete") {
			if o.PatientID == args[3] {
				results.OperationRecord = append(results.OperationRecord, o)
			}
		}
	}
	a := AuditRecord{
		AuditRecordID:  args[0],
		UserID:         args[1],
		OrganizationID: args[2],
		DataType:       "Audit",
		EntryMethod:    "Auto",
		Time:           time.Now().Format("2006-01-02 15:04:05"),
	}
	op, err4 := json.Marshal(a)
	if err4 != nil {
		return shim.Error("审计记录加工失败！")
	}
	err5 := stub.PutState(args[0], op)
	if err5 != nil {
		return shim.Error("审计失败！")
	}
	err6 := stub.SetEvent(args[4], []byte{})
	if err6 != nil {
		return shim.Error(err6.Error())
	}
	r, err7 := json.Marshal(results)
	if err7 != nil {
		return shim.Error("结果加工失败！")
	}
	return shim.Success(r)
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("启动链码时发生错误: %s", err)
	}
}
