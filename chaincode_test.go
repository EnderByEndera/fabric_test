package chaincode

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"testing"
)

func TestChaincode(t *testing.T) {
	cc := new(BankChaincode)
	stub := shimtest.NewMockStub("test", cc)

	//初始化，输出Success to initialize!
	res := stub.MockInit("1", [][]byte{[]byte("init")})
	fmt.Println("Init result:", string(res.Payload))

	//注册lihongyao，输出 Success to register Lihongyao!
	res = stub.MockInvoke("1", [][]byte{[]byte("register"), []byte("Lihongyao"), []byte("123456"), []byte("No")})
	fmt.Println("register Lihongyao result:", string(res.Payload))

	//再次注册lihongyao，应该输出到error，无结果显示
	res = stub.MockInvoke("1", [][]byte{[]byte("register"), []byte("Lihongyao"), []byte("123456"), []byte("No")})
	fmt.Println("register Lihongyao result:", string(res.Payload))

	//登出lihongyao，应该输出到error，无结果显示
	res = stub.MockInvoke("1", [][]byte{[]byte("loginOut"), []byte("Lihongyao"), []byte("123456")})
	fmt.Println("log out user Lihongyao result:", string(res.Payload))

	//登录lihongyao，输出Log in success!!
	res = stub.MockInvoke("1", [][]byte{[]byte("login"), []byte("Lihongyao"), []byte("123456")})
	fmt.Println("log in user Lihongyao result:", string(res.Payload))

	//再次登录lihongyao，应该输出到error，无结果显示
	res = stub.MockInvoke("1", [][]byte{[]byte("login"), []byte("Lihongyao"), []byte("123456")})
	fmt.Println("log in user Lihongyao result:", string(res.Payload))

	//登出lihongyao，输出Log out success!!
	res = stub.MockInvoke("1", [][]byte{[]byte("loginOut"), []byte("Lihongyao"), []byte("123456")})
	fmt.Println("log out user Lihongyao result:", string(res.Payload))

	//修改lihongyao的密码，输出 Alter the password success!!
	res = stub.MockInvoke("1", [][]byte{[]byte("alterPasswd"), []byte("Lihongyao"), []byte("123456"), []byte("654321")})
	fmt.Println("log out user Lihongyao result:", string(res.Payload))

	//修改lihongyao的密码，使用错误的旧密码，应该输出到error，无结果显示
	res = stub.MockInvoke("1", [][]byte{[]byte("alterPasswd"), []byte("Lihongyao"), []byte("123456"), []byte("654321")})
	fmt.Println("log out user Lihongyao result:", string(res.Payload))

	//登录lihongyao，使用错误的密码，应该输出到error，无结果显示
	res = stub.MockInvoke("1", [][]byte{[]byte("login"), []byte("Lihongyao"), []byte("123456")})
	fmt.Println("log in user Lihongyao result:", string(res.Payload))

	//登录lihongyao，输出Log in success!!
	res = stub.MockInvoke("1", [][]byte{[]byte("login"), []byte("Lihongyao"), []byte("654321")})
	fmt.Println("log in user Lihongyao result:", string(res.Payload))

	//登出lihongyao，使用错误的密码，应该输出到error，无结果显示
	res = stub.MockInvoke("1", [][]byte{[]byte("loginOut"), []byte("Lihongyao"), []byte("123456")})
	fmt.Println("log out user Lihongyao result:", string(res.Payload))

	//登出lihongyao，输出Log out success!!
	res = stub.MockInvoke("1", [][]byte{[]byte("loginOut"), []byte("Lihongyao"), []byte("654321")})
	fmt.Println("log out user Lihongyao result:", string(res.Payload))

	//再次登出lihongyao，应该输出到error，无结果显示
	res = stub.MockInvoke("1", [][]byte{[]byte("loginOut"), []byte("Lihongyao"), []byte("654321")})
	fmt.Println("log out user Lihongyao result:", string(res.Payload))

	//创建lihongyao账户，输出Create account: Lihongyao@Account  is success!
	res = stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("Lihongyao@Account"), []byte("0")})
	fmt.Println("create lihongyao's account result: ", string(res.Payload))

	//向lihongyao账户存入10，输出 Add is success! Account: Lihongyao@Account; Remaining balance is: 10
	res = stub.MockInvoke("1", [][]byte{[]byte("add"), []byte("Lihongyao@Account"), []byte("10")})
	fmt.Println("add lihongyao's account result: ", string(res.Payload))

	//从lihongyao账户取钱，输出 Reduce is success! Account: Lihongyao@Account; Remaining balance is: 0
	res = stub.MockInvoke("1", [][]byte{[]byte("reduce"), []byte("Lihongyao@Account"), []byte("10")})
	fmt.Println("reduce lihongyao's account result: ", string(res.Payload))

	//创建wangyiwen账户，输出Create account: Wangyiwen@Account  is success!
	res = stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("Wangyiwen@Account"), []byte("0")})
	fmt.Println("create wangyiwen's account result: ", string(res.Payload))

	//查询wangyiwen账户的余额，余额为0 ，输出 Account: Wangyiwen@Account; Balance: 0
	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Wangyiwen@Account")})
	fmt.Println("get wangyiwen's account result: ", string(res.Payload))

	//创建wangyiwen账户，应该输出到error，无结果显示
	res = stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("Wangyiwen@Account"), []byte("10")})
	fmt.Println("create wangyiwen's account result: ", string(res.Payload))

	//lihongyao向wangyiwen转账10，交易失败，应该输出到error，无结果显示
	res = stub.MockInvoke("1", [][]byte{[]byte("transfer"), []byte("Lihongyao@Account"), []byte("Wangyiwen@Account"), []byte("10")})
	fmt.Println("transfer lihongyao wangyiwen result: ", string(res.Payload))

	//查询wangyiwen账户的余额，应该为0，输出Account: Wangyiwen@Account; Balance: 0
	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Wangyiwen@Account")})
	fmt.Println("get wangyiwen's account result: ", string(res.Payload))

	//查询lihongyao账户的余额，应该为0，输出 Account: Lihongyao@Account; Balance: 0
	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Lihongyao@Account")})
	fmt.Println("get lihongyao's account result: ", string(res.Payload))

	//向lihongyao账户存钱，100，输出Add is success! Account: Lihongyao@Account; Remaining balance is: 100
	res = stub.MockInvoke("1", [][]byte{[]byte("add"), []byte("Lihongyao@Account"), []byte("100")})
	fmt.Println("add lihongyao's account result: ", string(res.Payload))

	//lihongyao向wangyiwen转账10，输出Transfer is success!
	res = stub.MockInvoke("1", [][]byte{[]byte("transfer"), []byte("Lihongyao@Account"), []byte("Wangyiwen@Account"), []byte("10")})
	fmt.Println("transfer lihongyao wangyiwen result: ", string(res.Payload))

	//查询wangyiwen余额，应该为10，输出Account: Wangyiwen@Account; Balance: 10
	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Wangyiwen@Account")})
	fmt.Println("get wangyiwen's account result: ", string(res.Payload))

	//查询lihongyao余额，应该为90，输出  Account: Lihongyao@Account; Balance: 90
	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Lihongyao@Account")})
	fmt.Println("get lihongyao's account result: ", string(res.Payload))

	//查询wangyiwen收入的交易记录
	res = stub.MockInvoke("1", [][]byte{[]byte("query"), []byte("in"), []byte("Wangyiwen@Account")})
	fmt.Println("query In wangyiwen's account result: ", string(res.Payload))

	//查询lihongyao支出的交易记录
	res = stub.MockInvoke("1", [][]byte{[]byte("query"), []byte("out"), []byte("Lihongyao@Account")})
	fmt.Println("query Out lihongyao's account result: ", string(res.Payload))

	//创建changwen的账户
	res = stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("Changwen@Account"), []byte("0")})
	fmt.Println("create changwen's account result: ", string(res.Payload))

	//删除changwen的账户
	res = stub.MockInvoke("1", [][]byte{[]byte("delete"), []byte("Changwen@Account")})
	fmt.Println("delete changwen's account result: ", string(res.Payload))

	//查询changwen的余额，应该输出到error，无结果显示
	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Changwen@Account")})
	fmt.Println("get changwen's account result: ", string(res.Payload))

	//回滚lihongyao和wangyiwen的事务id为1的交易
	res = stub.MockInvoke("1", [][]byte{[]byte("RollBack"), []byte("Lihongyao@Account"), []byte("Wangyiwen@Account"), []byte("1")})
	fmt.Println("RollBack lihongyao wangyiwen 1 result: ", string(res.Payload))

	//查询wangyiwen的余额，应该是0， Account: Wangyiwen@Account; Balance: 0,可实际输出 Account: Wangyiwen@Account; Balance: 10
	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Wangyiwen@Account")})
	fmt.Println("get wangyiwen's account result: ", string(res.Payload))

	//查询lihongyao的余额，应该是100，输出 Account: Lihongyao@Account; Balance: 100，可实际输出 Account: Wangyiwen@Account; Balance: 90
	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Lihongyao@Account")})
	fmt.Println("get lihongyao's account result: ", string(res.Payload))

	//回滚lihongyao和changwen的事务id为1的交易，应该输出到error，无结果显示
	res = stub.MockInvoke("1", [][]byte{[]byte("RollBack"), []byte("Lihongyao@Account"), []byte("Changwen@Account"), []byte("1")})
	fmt.Println("RollBack lihongyao changwen 1 result: ", string(res.Payload))
}

func BenchmarkCreateGetDelete(b *testing.B) {
	cc := new(BankChaincode)
	stub := shimtest.NewMockStub("test", cc)

	for i := 0; i < b.N; i++ {
		stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("Lihongyao"), []byte("0")})
		stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Lihongyao")})
		stub.MockInvoke("1", [][]byte{[]byte("delete"), []byte("Lihongyao")})
	}
}

func BenchmarkCreateTransferQueryRollBack(b *testing.B) {
	cc := new(BankChaincode)
	stub := shimtest.NewMockStub("test", cc)

	for i := 0; i < b.N; i++ {
		stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("wangyiwen"), []byte("0")})
		stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("lihongyao"), []byte("100000")})
		stub.MockInvoke("1", [][]byte{[]byte("transfer"), []byte("lihongyao"), []byte("wangyiwen"), []byte("1")})
		stub.MockInvoke("1", [][]byte{[]byte("query"), []byte("out"), []byte("lihongyao")})
		stub.MockInvoke("1", [][]byte{[]byte("query"), []byte("in"), []byte("wangyiwen")})
		stub.MockInvoke("1", [][]byte{[]byte("RollBack"), []byte("lihongyao"), []byte("wangyiwen"), []byte("1")})
	}
}

func BenchmarkCreateAddReduceDelete(b *testing.B) {
	cc := new(BankChaincode)
	stub := shimtest.NewMockStub("test", cc)

	for i := 0; i < b.N; i++ {
		stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("wangyiwen"), []byte("0")})
		stub.MockInvoke("1", [][]byte{[]byte("add"), []byte("wangyiwen"), []byte("1")})
		stub.MockInvoke("1", [][]byte{[]byte("reduce"), []byte("wangyiwen"), []byte("1")})
		stub.MockInvoke("1", [][]byte{[]byte("delete"), []byte("wangyiwen")})
	}
}
