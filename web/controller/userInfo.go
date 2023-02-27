package controller

import (
	"medical/abac"
	"medical/service"
    "database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"medical/web/mysql"
)

type Application struct {
	Setup *service.ServiceSetup
}

type User struct {
	LoginName string
	Password  string
	Identity  string
	//IsAdmin   string
}

type Data struct {
	CurrentUser User
	Msg         string
	Flag        bool
	Med         service.MedicalRecord
	History     bool
	Ops         service.OperationRecordArr
	Repo        service.AuditReport //*新增：审计报告格式
	AuditString string
	Policy      abac.Policy
	Table       []service.TableRow
}

var users []User

func init() {

	//admin := User{LoginName: "admin", Password: "123456", IsAdmin: "T"}
	//alice := User{LoginName: "ChainDesk", Password: "123456", IsAdmin: "T"}
	//bob := User{LoginName: "alice", Password: "123456", IsAdmin: "F"}
	//jack := User{LoginName: "bob", Password: "123456", IsAdmin: "F"}
	//
	//users = append(users, admin)
	//users = append(users, alice)
	//users = append(users, bob)
	//users = append(users, jack)

}



//search section

//导入二叉树自定义包 biTree

type btNode struct {
	Data           interface{}
	Lchild, Rchild *btNode
}

type biTree struct {
	root *btNode
}

type multi_branch_tree struct {
	searchtime int
	children   []*biTree
}

func Create_Multi_branch_tree(data [][]interface{}) *multi_branch_tree {

	multi_tree := &multi_branch_tree{}
	//二维数组的数据中每一组构成一个子树
	for m := 0; m < len(data); m++ {
		list := data[m]
		multi_tree.children = append(multi_tree.children, Create(list))
	}
	//初始搜中次数为0
	multi_tree.searchtime = 0
	return multi_tree
}

func Create(data interface{}) *biTree {
	var list []interface{}
	btree := &biTree{}
	switch data.(type) {
	case []interface{}:
		list = append(list, data.([]interface{})...)
	default:
		list = append(list, data)
	}
	if len(list) > 0 {
		btree.root = &btNode{Data: list[0]}
		for _, data := range list[1:] {
			btree.AppendNode(data)
		}
	}
	return btree
}

func (bt *biTree) Append(data interface{}) {
	var list []interface{}
	switch data.(type) {
	case []interface{}:
		list = append(list, data.([]interface{})...)
	default:
		list = append(list, data)
	}
	if len(list) > 0 {
		for _, data := range list {
			bt.AppendNode(data)
		}
	}
}

func (bt *biTree) AppendNode(data interface{}) {
	root := bt.root
	if root == nil {
		bt.root = &btNode{Data: data}
		return
	}
	Queue := []*btNode{root}
	for len(Queue) > 0 {
		cur := Queue[0]
		Queue = Queue[1:]
		if cur.Lchild != nil {
			Queue = append(Queue, cur.Lchild)
		} else {
			cur.Lchild = &btNode{Data: data}
			return
		}
		if cur.Rchild != nil {
			Queue = append(Queue, cur.Rchild)
		} else {
			cur.Rchild = &btNode{Data: data}
			break
		}
	}
}

// 广度优先搜索
func (bt *biTree) LeafNodeBFS_search(searchkey interface{}) interface{} {
	var res interface{}
	root := bt.root
	if root == nil {
		return res
	}
	Queue := []*btNode{root}
	for len(Queue) > 0 {
		cur := Queue[0]
		Queue = Queue[1:]
		if cur != nil && cur.Data == searchkey {
			res = root.Data
			return res
		}
		if cur.Lchild != nil {
			Queue = append(Queue, cur.Lchild)
		}
		if cur.Rchild != nil {
			Queue = append(Queue, cur.Rchild)
		}
	}
	return res
}



// 深度优先搜索
func (bt *biTree) LeafNodeDFS_search(search_key interface{}) interface{} {
	var res interface{}
	cur := bt.root
	Stack := []*btNode{}
	for cur != nil || len(Stack) > 0 {
		for cur != nil {
			// if cur.Lchild == nil && cur.Rchild == nil {
			// 	res = append(res, cur.Data)
			// }
			if cur.Data == search_key {
				res = bt.root.Data
				return res
			}
			Stack = append(Stack, cur)
			cur = cur.Lchild
		}
		if len(Stack) > 0 {
			cur = Stack[len(Stack)-1]
			Stack = Stack[:len(Stack)-1]
			cur = cur.Rchild
		}
	}
	return res
}



// 查询数据库全部数据
func Check_sqldata() [][]interface{} {
	sql_str := "select _CaseNumber,_Diseases,_SubjectMark,_Organization,_Diagnose,_GatherTime,_Researcher from base_info "
	mydb := mysql.Initdb()
	rows, err := mydb.Query(sql_str)
	var sql_data [][]interface{}
	if err == nil {
		var result [2][64]string
		columns, _ := rows.Columns()
		for i := range columns {
			result[0][i] = columns[i]
		}
		values := make([]sql.RawBytes, len(columns))
		//定义一个切片,元素类型是interface{} 接口
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			//把sql.RawBytes类型的地址存进去了
			scanArgs[i] = &values[i]
		}
		//获取字段值

		for rows.Next() {
			rows.Scan(scanArgs...)
			fmt.Print("\n")
			var temp []interface{}
			for i, col := range values {
				if i != len(columns)-1 {
					temp = append(temp, string(col))
				} else {
					temp = append(temp, string(col))
				}
			}
			sql_data = append(sql_data, temp)
		}
		rows.Close()

	} else {
		fmt.Print(err.Error())
	}

	//最终关闭数据库
	defer mydb.Close()
	return sql_data
}

// interface{}转string
func Strval(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

// 根据casenumber查询数据库数据
func Check_data(rs []interface{}) [][]interface{} {
	//根据casenumber搜索的总的数据库记录
	var data_sum [][]interface{}

	//rs中的每个casenumber都对应一条记录，存入data_sum中
	for _, i := range rs {
		sql_str := "select * from base_info where _CaseNumber='" + Strval(i)+"'"
		mydb := mysql.Initdb()
		rows, err := mydb.Query(sql_str)

		if err == nil {
			var result [2][64]string
			columns, _ := rows.Columns()
			for i := range columns {
				result[0][i] = columns[i]
			}
			values := make([]sql.RawBytes, len(columns))
			//定义一个切片,元素类型是interface{} 接口
			scanArgs := make([]interface{}, len(values))
			for i := range values {
				//把sql.RawBytes类型的地址存进去了
				scanArgs[i] = &values[i]
			}
			//获取字段值

			for rows.Next() {
				rows.Scan(scanArgs...)
				fmt.Print("\n")
				var temp []interface{}
				for i, col := range values {
					if i != len(columns)-1 {
						temp = append(temp, string(col))
					} else {
						temp = append(temp, string(col))
					}
				}
				data_sum = append(data_sum, temp)
			}
			rows.Close()

		} else {
			fmt.Print(err.Error())
		}

		//最终关闭数据库
		defer mydb.Close()

	}
	return data_sum
}

// 广度搜索，效率优先的搜索方式
func (bt *multi_branch_tree) LeafNodeBFS_key(data interface{}) []interface{} {
	var rs []interface{}
	for _, i := range bt.children {
		//每一个子树都要进行搜索
		search_reult := i.LeafNodeBFS_search(data)
		//如果没搜到就停止本课子树的搜索，进行下一棵树的搜索
		if search_reult == nil {
			continue
		}
		bt.searchtime++
		rs = append(rs, search_reult)
	}
	return rs
}

// 深度搜索，准确度优先的搜索方式
func (bt *multi_branch_tree) LeafNodeDFS_key(data interface{}) []interface{} {
	var rs []interface{}
	for _, i := range bt.children {
		search_result := i.LeafNodeDFS_search(data)
		if search_result == nil {
			continue
		}
		bt.searchtime++
		rs = append(rs, search_result)
	}
	return rs
}