package machine

import (
	"tool/errors"
	//"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

var dbu UserDB

func init() {
	var err error

	dbu, err = NewUserDB()
	if err != nil {
		if err == UserDataNotFound {
			fmt.Println(UserDataNotFound.Error())
			return
		}
		panic(err)
	}
}

func TestCheckUserByName(t *testing.T) {
	isRight, err := dbu.CheckAccount("test_consumer1")
	if err != nil {
		if err == UserDataNotFound {
			fmt.Println(UserDataNotFound.Error())
			return
		}
		t.Fatal(err)
	}
	fmt.Println(isRight)
}

func TestLoginCheck(t *testing.T) {
	isRight, id, err := dbu.LoginCheck("test_consumer2", "a123456A")
	if err != nil {
		if err == UserDataNotFound {
			fmt.Println(UserDataNotFound.Error())
			return
		}
		if err == UserWrongPass {
			fmt.Println("用户密码错误")
			return
		}
		t.Fatal(err)
	}
	fmt.Println(id, isRight)

}

func TestGetUserBaseInfoById(t *testing.T) {
	user, err := dbu.GetUserBaseInfoById(1)
	if err != nil {
		if UserDataNotFound.Equal(err) {
			fmt.Println(err)
			return
		}
		t.Fatal(err)
	}
	fmt.Println(user)

}

func TestGetUserBalanceInfoById(t *testing.T) {
	user, err := dbu.GetUserBalanceInfoById(1)
	if err != nil {
		if err == UserDataNotFound {
			fmt.Println(errors.As(err))
			return
		}
		t.Fatal(err)
	}
	fmt.Println(user)

}

func TestGetUserBaseInfoByAccount(t *testing.T) {
	user, err := dbu.GetUserBaseInfoByAccount("test_consumer1")
	if err != nil {
		if err == UserDataNotFound {
			fmt.Println(UserDataNotFound.Error())
			return
		}
		t.Fatal(err)
	}
	fmt.Println(user)

}

func TestGetUserBalanceInfoByAcount(t *testing.T) {
	user, err := dbu.GetUserBalanceInfoByAccount("test_consumer1")
	if err != nil {
		if err == UserDataNotFound {
			fmt.Println(UserDataNotFound.Error())
			return
		}
		t.Fatal(err)
	}
	fmt.Println(user)

}

func TestUserParentNode(t *testing.T) {
	user, err := dbu.GetParentNodeDirect(4)
	if err != nil {
		if err == UserDataNotFound {
			fmt.Println(UserDataNotFound.Error())
			return
		}
		t.Fatal(err)
	}
	fmt.Println(user)
}

func TestGetChildNodeInfo(t *testing.T) {
	total, err := dbu.GetChildNodeNum(100000, "", -1, "", -1, -1)
	if err != nil {
		if err == UserDataNotFound {
			fmt.Println(UserDataNotFound.Error())
			return
		}
		t.Fatal(err)
	}
	list, err := dbu.GetChildNodeInfo(100000, "", -1, "", -1, -1)
	if err != nil {
		if err == UserDataNotFound {
			fmt.Println(UserDataNotFound.Error())
			return
		}
		t.Fatal(err)
	}
	fmt.Println(list, total)
}

func TestGetChildNodeDirect(t *testing.T) {
	list, err := dbu.GetChildNodeDirect(2, "", 0)
	if err != nil {
		if err == UserDataNotFound {
			fmt.Println(UserDataNotFound.Error())
			return
		}
		t.Fatal(err)
	}
	fmt.Println(list)
}

func TestGetParentNodeInfo(t *testing.T) {
	list, err := dbu.GetParentNodeInfo(4, "ADMIN", 0)
	if err != nil {
		if err == UserDataNotFound {
			fmt.Println(UserDataNotFound.Error())
			return
		}
		t.Fatal(err)
	}
	fmt.Println(list)
}

/*
func TestCreateUser(t *testing.T){
	user:=&UserBaseInfo{
		//ParentId:100000,
		Mobile:"123132131",
		Account:"hazhao",
		RealName:"hazho",
		RoleCode:CONSUMER,
		IdCard:"12",
		BankCard:"12",
	}
	uid,err:=dbu.CreateUser(1,user,"123456")
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println(uid)
}
*/
func TestGetOperationRole(t *testing.T) {
	list := make([]string, 0)
	list = append(list, "111")
	str, err := json.Marshal(list)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(str))
	userRole, err := dbu.GetOperationRole(3, CreateRole)
	if err != nil {
		if err == RoleDataNotFound {
			fmt.Println(RoleDataNotFound.Error())
			return
		}
		t.Fatal(err)
	}
	fmt.Println(userRole)

}

func TestGetOperatePermission(t *testing.T) {
	roleA, err := dbu.GetOperatePermission(1, "sale_role", "")
	if err != nil {
		if err == RoleDataNotFound {
			fmt.Println(RoleDataNotFound.Error())
			return
		}
		t.Fatal(err)
	}
	fmt.Println("sale_role of ROOT is :", roleA)
	if roleA {
		t.Fatal(errors.New("鉴权程序出错！"))
	}

	roleB, err := dbu.GetOperatePermission(2, CreateRole, "ROOT")
	if err != nil {
		if err == RoleDataNotFound {
			fmt.Println(RoleDataNotFound.Error())
			return
		}
		t.Fatal(err)
	}
	fmt.Println("create_role of ADIMIN to ROOT is :", roleB)
	if roleB {
		t.Fatal(errors.New("鉴权程序出错！"))
	}

	roleC, err := dbu.GetOperatePermission(4, "sale_role", "")
	if err != nil {
		if err == RoleDataNotFound {
			fmt.Println(RoleDataNotFound.Error())
			return
		}
		t.Fatal(err)
	}
	fmt.Println("sale_role of CUNSUMER is :", roleC)
	if !roleC {
		t.Fatal(errors.New("鉴权程序出错！"))
	}

	roleD, err := dbu.GetOperatePermission(4, CreateRole, "CONSUMER")
	if err != nil {
		if err == RoleDataNotFound {
			fmt.Println(RoleDataNotFound.Error())
			return
		}
		t.Fatal(err)
	}
	fmt.Println("create_role of  CUNSUMER to CUNSUMER is :", roleD)
	if !roleD {
		t.Fatal(errors.New("鉴权程序出错！"))
	}

	roleE, err := dbu.GetOperatePermission(4, DistributeRole, "CONSUMER")
	if err != nil {
		if err == RoleDataNotFound {
			fmt.Println(RoleDataNotFound.Error())
			return
		}
		t.Fatal(err)
	}
	fmt.Println("distribute of  CUNSUMER to CUNSUMER is :", roleE)
	if !roleE {
		t.Fatal(errors.New("鉴权程序出错！"))
	}

}

func TestUpdateUserBaseInfo(t *testing.T) {
	user, err := dbu.GetUserBaseInfoById(1)
	if err != nil {
		if err == UserDataNotFound {
			fmt.Println(UserDataNotFound.Error())
			return
		}
		t.Fatal(err)
	}
	realNameSave := user.RealName
	user.RealName = "测试真名[" + time.Now().Format("01-02 15:04") + "]"
	if err := dbu.UpdateUserBaseInfo(user); err != nil {
		t.Fatal(err)
	}

	user.RealName = realNameSave
	if err := dbu.UpdateUserBaseInfo(user); err != nil {
		t.Fatal(err)
	}

	fmt.Println("TestUpdateUserBaseInfo ok")
}

func TestUpdateUserStatus(t *testing.T) {
	if err := dbu.UpdateUserStatus(123213, Frozen); err != nil {
		t.Fatal(err)
	}
	fmt.Println("TestUpdateUserStatus ok")
}

func TestApplyMoneyWithdrew(t *testing.T) {
	if err := dbu.ApplyMoneyWithdrew(4, "1"); err != nil {
		if err == UserNotEnoughMoney {
			fmt.Println(err.Error())
		} else {
			t.Fatal(err)
		}
	} else {
		fmt.Println("TestApplyMoneyWithdrew ok")
	}

}

func TestWithdrewMoneyForChild(t *testing.T) {
	if err := dbu.WithdrewMoneyForChild(4, 2, "1"); err != nil {
		if err == UserNotEnoughMoney {
			fmt.Println(err.Error())
		} else {
			t.Fatal(err)
		}
	} else {
		fmt.Println("TestWithdrewMoneyForChild ok")
	}
}

func TestCaculate(t *testing.T) {
	fmt.Println(10 % 9)
}

func TestCheckChildNode(t *testing.T) {
	is, err := dbu.CheckChildNodeInfo(100003, 3)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(is)
}


func TestUpdateUserPass(t *testing.T) {
	err := dbu.UpdateUserPass(3, "321", "asd")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("TestUpdateUserPass ok")
}

func TestHelpUpdateUserPass(t *testing.T) {
	err := dbu.HelpUpdatePass("12122", 4,"a123456A", "dsa")
	if err != nil {
		if(errors.Equal(err,UserNotEnoughAuth)){
			fmt.Println("权限不足")
			return
		}
		t.Fatal(err)
	}
	fmt.Println("TestUpdateUserPass ok")
}
