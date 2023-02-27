package main

import (
	"fmt"
	"medical/sdkInit"
	"medical/service"
	"medical/web"
	"medical/web/controller"
	"os"
)

// "medical/sdkInit"
// "medical/web"
// "medical/web/controller"
// "os"

const (
	cc_name    = "simplecc"
	cc_version = "1.0.0"
)

var SetPath string

func main() {
	// SetPath = "/workspace/github.com/medical/"
	SetPath = "/opt/gopath/src/github.com/hyperledger/fabric-samples/medical/"
	// SetPath = "/Users/monk/code/go/src/medical/"
	// init orgs information
	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: SetPath + "fixtures/channel-artifacts/Org1MSPanchors.tx",
		},
	}

	// init sdk env info
	info := sdkInit.SdkEnvInfo{
		ChannelID:        "mychannel",
		ChannelConfig:    SetPath + "fixtures/channel-artifacts/channel.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      cc_name,
		ChaincodePath:    SetPath + "chaincode/",
		ChaincodeVersion: cc_version,
	}

	// sdk setup
	// sdk, err := sdkInit.Setup("config_monk.yaml", &info)
	sdk, err := sdkInit.Setup("config.yaml", &info)
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}

	// create channel and join
	if err := sdkInit.CreateAndJoinChannel(&info); err != nil {
		fmt.Println(">> Create channel and join error:", err)
		os.Exit(-1)
	}

	// create chaincode lifecycle
	if err := sdkInit.CreateCCLifecycle(&info, 1, false, sdk); err != nil {
		fmt.Println(">> create chaincode lifecycle error: %v", err)
		os.Exit(-1)
	}

	// invoke chaincode set status
	fmt.Println(">> 通过链码外部服务设置链码状态......")
	serviceSetup, err := service.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk)
	if err != nil {
		fmt.Println()
		os.Exit(-1)
	}
	/*
		for i := 1; i <= 5; i++ {
			arr := [7]string{"m" + strconv.Itoa(i), "u" + strconv.Itoa(i), "p" + strconv.Itoa(i), "o" + strconv.Itoa(i), "arg", "www", "Manual"}
			msg, err := serviceSetup.UploadMed(arr[:])
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("信息发布成功, 交易编号为: " + msg)
			}
		}
	*/
	app := controller.Application{
		Setup: serviceSetup,
	}
	web.WebStart(app)
}

// func main() {

// 	// 读文件内容
// 	// fileContent, _ := file.Open()
// 	file, err := os.Open("/home/monk/test.csv")
// 	if err != nil {
// 		panic(err)
// 	}
// 	//func (f *File) Close() error {}
// 	defer file.Close()
// 	//func ReadAll(r io.Reader) ([]byte, error) {}
// 	r1 := csv.NewReader(file)
// 	content, err := r1.ReadAll()
// 	if err != nil {
// 		log.Fatalf("can not readall, err is %+v", err)
// 	}
// 	for _, row := range content[1:][:] {
// 		fmt.Println(row)
// 			// info, err := app.Setup.UploadMed(arr[:])

// 	}
// }
