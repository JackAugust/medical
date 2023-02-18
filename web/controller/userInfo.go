package controller

import (
	"medical/abac"
	"medical/service"
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
