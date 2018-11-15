package bcrypt

import (
    "github.com/jameskeane/bcrypt"
)

// 加密
func BcryptHash(oriPwd string) (pwd, salt string) {
    salt, _ = bcrypt.Salt(10)
    pwd, _ = bcrypt.Hash(oriPwd, salt)
    return  
}

// 验证
func BcryptMatch(oriPwd, encodePwd string) bool {
    return bcrypt.Match(oriPwd, encodePwd)
}
