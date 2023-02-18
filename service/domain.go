package service

import (
	"database/sql"
	"fmt"
	_ "hash/fnv"
	"medical/abac"
	"medical/sdkInit"
	"time"

	_ "github.com/go-sql-driver/mysql" //初始化
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

//写到服务层里面，所有字段用于和数据库对接（服务层调用数据库和链码）
type MedicalRecord struct {
	Groups         string `json:"groups"`         //分组
	SubjectMark    string `json:"subjectMark"`    //样本标识符
	Name           string `json:"name"`           //姓名
	NameInitials   string `json:"nameInitials"`   //姓名缩写
	CaseNumber     string `json:"caseNumber"`     //病例号，并上链
	Sex            string `json:"sex"`            //性别
	Nation         string `json:"nation"`         //民族
	Diseases       string `json:"diseases"`       //疾病种类
	MedicalHistory string `json:"medicalHistory"` //现病史
	NativePlace    string `json:"nativePlace"`    //现住地
	Diagnose       string `json:"diagnose"`       //诊断
	Researcher     string `json:"researcher"`     //研究者，并上链
	Organization   string `json:"organization"`   //机构，并上链
	Addition1      string `json:"addition1"`      //预留信息
	Addition2      string `json:"addition2"`      //预留信息
	Addition3      string `json:"addition3"`      //预留信息
	Status         string `json:"status"`         //状态
	EntryTime      string `json:"entryTime"`      //入组时间
	BaseTime       string `json:"baseTime"`       //基准时间
	GatherTime     string `json:"gatherTime"`     //采集时间-1，并上链
}

type OperationRecordArr struct {
	OperationRecord []OperationRecord
}

// 表格的每一行
type TableRow struct {
	FirstColumn   string
	SecondColumn  string
	ThirdColumn   string
	FourthColumn  string
	FifthColumn   string
	SixthColumn   string
	SeventhColumn string
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
	IsSuccess         bool   `json:"isSuccess"`         //*新增：操作是否成功，成功为True
	Time              string `json:"time"`              //时间
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

// *新增：审计报告，记录组织的失败率信息
type AuditReport struct {
	TargetOrg       string     `json:"targetOrg"`       //被审计组织ID
	CurrentCredit   float64    `json:"currentCredit"`   //组织现在的信誉值
	PreviousCredit  float64    `json:"previousCredit"`  //组织之前的信誉值
	CreditChange    string     `json:"creditChange"`    //组织信誉值变动情况：上升、下降、不变
	ReferenceRange  [2]float64 `json:"referenceRange"`  //参照区间
	TotalOperations int64      `json:"totalOperations"` //组织总操作次数
	FailOperations  int64      `json:"failOperations"`  //组织失败操作次数
	FailRate        float64    `json:"failRate"`        //组织失败操作率
	MaxFailRateUser string     `json:"maxFailRateUser"` //组织中失败率最高的用户ID
	MaxFailRate     float64    `json:"maxFailRate"`     //组织所有用户中最高的失败率
}

type ServiceSetup struct {
	ChaincodeID string
	Client      *channel.Client
}

//数据库配置
const (
	userName = "root"
	password = "root"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "itbtsql"
)

//Db数据库连接池
var DB *sql.DB

func regitserEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
	}
	return nil
}

func InitService(chaincodeID, channelID string, org *sdkInit.OrgInfo, sdk *fabsdk.FabricSDK) (*ServiceSetup, error) {
	handler := &ServiceSetup{
		ChaincodeID: chaincodeID,
	}
	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org.OrgUser), fabsdk.WithOrg(org.OrgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new channel client: %s", err)
	}
	handler.Client = client
	return handler, nil
}

//注意方法名大写，就是public
func InitDB() *sql.DB {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := userName + ":" + password + "@tcp(" + ip + ":" + port + ")/" + dbName + "?allowNativePasswords=true"
	fmt.Println(path)
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)
	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		fmt.Println("open database fail")
		fmt.Println(err)
	}
	fmt.Println("connnect success")
	return DB
}

// 用于生成sql查询语句
func SetSQL(sqlmap interface{}, queryString string) string {
	temp := sqlmap.(map[int]string)
	for k := range temp {
		fmt.Println(k, temp[k])
	}
	// "SELECT _SubjectMark From baseinfo WHERE _Groups='喹硫平'"
	// 可自行更改字段
	sql := "SELECT " + temp[1] + " From " + temp[0] + " WHERE " + temp[2] + "='" + queryString + "'"
	fmt.Println(sql)
	return sql
}

// 查询并返回数据
func queryDB(DB *sql.DB, Sql string) map[int]string {
	var subId string
	rows, err := DB.Query(Sql)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(subId)
	data := make(map[int]string)
	index := 0
	for rows.Next() {
		rows.Scan(&subId)
		data[index] = subId
		index = index + 1
		//fmt.Println(subId)
	}
	defer rows.Close()
	return data
}

func InsertDB(DB *sql.DB, data []string, casenumber string) bool {
	sqlString := "insert ignore into base_info(_Groups,_SubjectMark,_Name,_NameInitials,_CaseNumber,_Sex,_Nation,_Diseases,_MedicalHistory,_NativePlace,_Diagnose,_Researcher,_Organization,_Addition1,_Addition2,_Addition3,_Status,_EntryTime,_BaseTime,_GatherTime) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	r, err0 := DB.Exec(sqlString, data[0], data[1], data[2], data[3], casenumber, data[5], data[6], data[7], data[8], data[9], data[10], data[11], data[12], data[13], data[14], data[15], data[16], time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02 15:04:05"))
	if err0 != nil {
		fmt.Println("插入失败：", err0)
		return false
	}
	id, err1 := r.LastInsertId()
	if err1 != nil {
		fmt.Println("操作失败：", err1)
		return false
	}
	fmt.Println("插入成功：", id)
	//受影响的行数
	row_affect, err2 := r.RowsAffected()
	if err2 != nil {
		fmt.Println("受影响行数获取失败:", err2)
		return false
	}
	fmt.Println("受影响的行数：", row_affect)
	return true
}

//根据各种信息生成CaseNumber的函数，暂时没想好怎么写
func GenerateCaseNumber(data []string) string {
	//共有6个参数，分别为类型、返回编号、病种、所属人、机构、时间
	if len(data) != 6 {
		return ""
	}
	final_casenumber := "19260817"
	return final_casenumber //返回最终生成的CaseNumber
}

//将与区块链相关的数据先存到数据库中，在GenerateCaseNumber
func InsertChainBasedDataIntoDB(DB *sql.DB, data []string) bool {
	return true
}

func GeneratePolicy(DB *sql.DB, casenumber string) string {
	var g, s, d, r, o, t string
	rows := DB.QueryRow("select _Groups,_SubjectMark,_Diseases,_Researcher,_Organization,_GatherTime FROM base_info where _CaseNumber='" + casenumber + "'")
	rows.Scan(&g, &s, &d, &r, &o, &t)
	// result := make(map[int]string)
	// result = queryDB(DB, "select insti_name from insti_coop where _CaseNumber='"+casenumber+"'")
	// for _, v := range result {
	// 	o = o + " " + v
	// }
	attr := abac.Date2DB{
		Groups:       g,
		SubjectMark:  s,
		Diseases:     d,
		Researcher:   r,
		Organization: o,
	}
	policy := abac.AutoPolicyCreate(abac.Date2DB(attr), casenumber, t)
	var str string
	SQLString := "select * from policy where policy_id='" + casenumber + "'"
	err := DB.QueryRow(SQLString).Scan(&str)
	if err == sql.ErrNoRows { //没有结果
		SQLString3 := "insert into policy(policy_id,policy_data)values(?,?)"
		_, err := DB.Exec(SQLString3, casenumber, policy)
		if err != nil {
			fmt.Println("err")
		}
	}
	return policy
}

// func CheckAction(DB *sql.DB, casenumber string, action string) bool {

// 	var uid string
// 	SQLString := "select usertype from login where state='1'"
// 	err := DB.QueryRow(SQLString).Scan(&uid)
// 	if err == sql.ErrNoRows { //没有结果
// 		fmt.Println("err, can't find logined user")
// 		return false
// 	}
// 	var req abac.ABACRequest
// 	req.Obj = casenumber
// 	req.Op = action
// 	req.Sub = uid

// 	var user abac.Sub
// 	// user.UID = uid
// 	// var user_info string
// 	SQLString = "select * from user_type where user_id='" + uid + "'"
// 	err = DB.QueryRow(SQLString).Scan(&user.UID, &user.Role, &user.Org, &user.Group)
// 	if err == sql.ErrNoRows { //没有结果
// 		fmt.Println("err, can't find the user's information.")
// 		return false
// 	}
// 	var policy string
// 	SQLString = "select * from policy where policy_id='" + casenumber + "'"
// 	err = DB.QueryRow(SQLString).Scan(&policy)

// 	f := abac.CheckAccess(req, policy, user)
// 	return f
// }

func CheckAction(DB *sql.DB, casenumber string, action string) bool {

	var uid string
	SQLString := "select usertype from login where state='1'"
	err := DB.QueryRow(SQLString).Scan(&uid)
	if err != nil { //没有结果
		fmt.Println("err, can't find logined user")
		return false
	}
	var req abac.ABACRequest
	req.Obj = casenumber
	req.Op = action
	req.Sub = uid

	var user abac.Sub
	// user.UID = uid
	// var user_info string
	SQLString = "select * from user_type where user_id='" + uid + "'"
	err = DB.QueryRow(SQLString).Scan(&user.UID, &user.Role, &user.Org, &user.Group)
	if err == sql.ErrNoRows { //没有结果
		fmt.Println("err, can't find the user's information.")
		return false
	}
	// fmt.Println("user is ", user.UID, user.Role, user.Org, user.Group)

	var policy string
	SQLString = "select policy_data from policy where policy_id='" + casenumber + "'"
	err = DB.QueryRow(SQLString).Scan(&policy)
	// fmt.Println("policy is ", policy)
	f := abac.CheckAccess(req, policy, user)
	return f
}

func SelectDBSingle(DB *sql.DB, data []string) *MedicalRecord {
	sqlString := "select * from base_info where _CaseNumber = ?"
	stmt, _ := DB.Prepare(sqlString)
	row := stmt.QueryRow(data[3])
	if row == nil {
		fmt.Println("未获取到该记录:", data[3])
		return nil
	}
	m := &MedicalRecord{}
	err0 := row.Scan(&m.Groups, &m.SubjectMark, &m.Name, &m.NameInitials, &m.CaseNumber, &m.Sex, &m.Nation, &m.Diseases, &m.MedicalHistory, &m.NativePlace, &m.Diagnose, &m.Researcher, &m.Organization, &m.Addition1, &m.Addition2, &m.Addition3, &m.Status, &m.EntryTime, &m.BaseTime, &m.GatherTime)
	if err0 != nil {
		return nil
	}
	return m
}

func DeleteDB(DB *sql.DB, data []string) bool {
	sqlString := "delete  from base_info where _CaseNumber = ?"
	stmt, _ := DB.Prepare(sqlString)
	r, err0 := stmt.Exec(data[3])
	if err0 != nil {
		fmt.Println("删除记录失败：", err0)
		return false
	}
	row_affect, err1 := r.RowsAffected()
	if err1 != nil {
		fmt.Println("受影响行数获取失败:", err1)
		return false
	}
	fmt.Println("受影响的行数：", row_affect)
	return true
}

func InsertDB2Insti(DB *sql.DB, casenumber string) bool {
	fmt.Println(casenumber)
	var organization string
	SQLString := "select _Organization from base_info where _CaseNumber='" + casenumber + "'"
	err := DB.QueryRow(SQLString).Scan(&organization)
	if err == sql.ErrNoRows {
		return false
	} else {
		fmt.Println(organization)
		SQLString3 := "insert into insti_coop(_CaseNumber,insti_name)values(?,?)"
		_, err := DB.Exec(SQLString3, casenumber, organization)
		if err != nil {
			fmt.Println("err")
			return false
		} else {
			return true
		}
	}
}

func UpdateDB(DB *sql.DB, data []string) bool {
	sqlString := "update base_info set _Groups = ?, _SubjectMark = ?, _Name = ?, _NameInitials = ?, _Sex = ?, _Nation = ?, _Diseases = ?, _MedicalHistory = ?, _NativePlace = ?, _Diagnose = ?, _Researcher = ?, _Organization = ?, _Addition1 = ?, _Addition2 = ?, _Addition3 = ?, _Status = ?, _EntryTime = ?, _BaseTime = ? where _CaseNumber = ?;"
	r, err0 := DB.Exec(sqlString, data[0], data[1], data[2], data[3], data[5], data[6], data[7], data[8], data[9], data[10], data[11], data[12], data[13], data[14], data[15], data[16], time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02"), data[4])
	if err0 != nil {
		fmt.Println("修改失败：", err0)
		return false
	}
	//受影响的行数
	row_affect, err1 := r.RowsAffected()
	if err1 != nil {
		fmt.Println("受影响行数获取失败:", err1)
		return false
	}
	fmt.Println("受影响的行数：", row_affect)
	return true
}

// 定义返回的时间格式
var returnTimeType string

func ReturnTime(t int32) string {
	if t == 1 {
		returnTimeType = "2006-01-02"
	} else if t == 2 {
		returnTimeType = "2006-01-02 15:03:04"
	}
	now := time.Now()
	seconds := time.Unix(now.Unix(), 0)
	timeString := seconds.Format(returnTimeType)
	// 确定输出格式为string
	// fmt.Printf("%T\n", timeString)
	// 输出 时间戳，完整时间，最后返回内容
	// fmt.Println(now.Unix(), seconds, timeString)
	return timeString
}

// 返回hash——这个先放弃
/*
func hash(a Date2DB, t string) uint32 {
	baseString := []string{a.Groups, a.SubjectMark, a.Diseases, a.Researcher, a.Organization, ReturnTime(2)}
	s := strings.Join(baseString, "")
	h := fnv.New32a()
	h.Write([]byte(s))
	fmt.Println(h.Sum32())
	return h.Sum32()
}
*/
