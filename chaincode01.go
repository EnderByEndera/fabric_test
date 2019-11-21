package main

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"strconv"
	_ "strings"
	"time"
	//提供中间层来读取和修改账本

	pb "github.com/hyperledger/fabric-protos-go/peer"
	//封装执行结果的响应信息
)

type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Reduce {
	//实现链码初始化或升级时的处理逻辑
	//编码时可灵活使用stub中的API
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("初始化的参数只能为2个，分别代表名称和状态数据")
	}
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error("在保存状态时出现错误")
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fun, args := stub.GetFunctionAndParameters()
	//获取用户传递的函数名称及参数
	if expectedLen, ok := paramLength[fun]; !ok {

		return shim.Error("Undefined function")

	} else if expectedLen != len(args) {
		errStr, _ := paramLengthError[fun]
		return shim.Error(errStr)

	}

	var result string
	var err error

	switch fun {
	case "get":
		result, err = get(stub, args)

	case "add":
		result, err = add(stub, args)

	case "reduce":
		result, err = reduce(stub, args)

	case "create":
		result, err = create(stub, args)

	case "delete":
		result, err = delete(stub, args)

	default:
		return shim.Error("操作失败！")
	}

	return shim.Success([]byte(result))
}

//查询资产
func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	result, err := stub.GetState(args[0])

	if err != nil {
		return "", fmt.Errorf("获取%s数据发生%s错误", args[0], err)
	}

	if result == nil {
		return "", fmt.Errorf("根据%s没有获取到相应的数据", args[0])
	}
	return string(result), nil
}

//增加资产
func add(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	valueTemp, err := stub.GetState(args[0])

	if err != nil {
		return "", fmt.Errorf("没有获取%s数据因为%s", args[0], err)
	}
	intArgs1, err := strconv.Atoi(args[1])

	if err != nil {
		return "", fmt.Errorf("发生%s错误", err)
	}

	intValueTemp, err := strconv.Atoi(string(valueTemp))
	if err != nil {
		return "", fmt.Errorf("发生%s错误", err)
	}

	err = stub.PutState(args[0], []byte(strconv.Itoa(intArgs1+intValueTemp)))

	if err != nil {
		return "", fmt.Errorf("操作%s数据失败，发生%s错误", args[0], err)
	}

	return fmt.Sprintf("成功添加%s数据，余额是：%d", args[0], intValueTemp+intArgs1), nil
}

//减少资产
func reduce(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	valueTemp, err := stub.GetState(args[0])

	if err != nil {
		return "", fmt.Errorf("操作%s数据失败，发生%s错误", args[0], err)
	}
	intArgs1, err := strconv.Atoi(args[1])

	if err != nil {
		return "", fmt.Errorf("发生%s错误，操作失败", err)
	}
	intValueTemp, err := strconv.Atoi(string(valueTemp))

	if err != nil {
		return "", fmt.Errorf("发生%s错误，操作失败", err)
	}

	if intArgs1 > intValueTemp {
		return "", fmt.Errorf("%s'的账户余额不够减", args[0])
	}

	err = stub.PutState(args[0], []byte(strconv.Itoa(intValueTemp-intArgs1)))

	if err != nil {
		return "", fmt.Errorf("操作%s数据失败，发生%s错误", args[0], err)
	}
	return fmt.Sprintf("操作成功数据：%s，余额是%d", args[0], intValueTemp-intArgs1), nil
}

//创造资产
func create(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	var name []byte
	name, err := stub.GetState(args[0])
	if name != nil {
		return "", fmt.Errorf(fmt.Sprintf("账户已存在！"))
	}

	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("操作%s数据失败，发生%s错误", args[0], err))
	}
	err = stub.PutState(args[0], []byte(args[1]))

	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("操作%s数据失败，发生%s错误", args[0], err))
	}
	return fmt.Sprintf("成功创建%s账户!", args[0]), nil
}

//删除资产
func delete(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	err := stub.DelState(args[0])

	if err != nil {
		return "", fmt.Errorf("操作%s数据失败，发生%s错误", args[0], err)
	}
	return fmt.Sprintf("成功删除数据： %s", args[0]), nil

}

//密钥
func createHistoryKey(stub shim.ChaincodeStubInterface, args []string, first string) (string, error) {
	FormatTime, err := stub.GetTxTimestamp()

	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Get transaction timestamp failed!"))
	}

	tm := time.Unix(FormatTime.Seconds, 0)
	if first == "out" {
		historyKey, err := stub.CreateCompositeKey(first, []string{
			args[0], "->", args[1],
			"\t", stub.GetTxID(),
			"\t", tm.Format("Mon Jan 2 15:04:05 +0800 UTC 2006"),
		})

		if err != nil {
			return "", fmt.Errorf("发生%s错误，创建密钥失败", err)
		}

		err = stub.PutState(historyKey, []byte(args[2]))
		if err != nil {
			return "", fmt.Errorf("发生%s错误，操作失败", err)
		}
	} else if first == "in" {
		historyKey, err := stub.CreateCompositeKey(first, []string{
			args[1], "<-", args[0],
			"\t", stub.GetTxID(),
			"\t", tm.Format("Mon Jan 2 15:04:05 +0800 UTC 2006"),
		})

		if err != nil {
			return "", fmt.Errorf("发生%s错误，操作失败", err)
		}
		err = stub.PutState(historyKey, []byte(args[2]))

		if err != nil {
			return "", fmt.Errorf("发生%s错误，操作失败", err)
		}
	}
	return fmt.Sprintf("Insert records success!"), nil
}
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Println("启动SimpleChaincode时发生错误：%s", err)
	}
}
