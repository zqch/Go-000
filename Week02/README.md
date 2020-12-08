# 异常处理
我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

## Wrap 原则
- application 才选择 wrap error。基础库、公共组件等，只返回根错误值
- 如果 error 不打算处理，那就用上下文 wrap 起来抛上去
- error 只被 handle 一次，一旦处理了它就不是错误了，就返回 nil
- 和其他库协作时，其他库遵循第一条原则只返回根错误，因此自身需要 wrap
- 直接返回错误，只在程序顶部或是工作的 goroutine 顶部打日志

因此 dao 层遇到 `sql.ErrNoRows` 应该 wrap 返回给上层
```go
func GetUser(db *DB, id int) (User, error) {
	var user User
	// ctx, query...
	ret, err := db.Select(ctx, query)  // 查库，这里有可能返回 sql.ErrNoRows
	if err == nil {
		err = FillUser(&user, ret)  // 填充字段
	}
	return user, errors.Wrap(err)  // dao 层 error 包装再返回
}
```

service 层收到 dao 层返回的 error，通常不需要处理，既不需要包装也不需要打日志，直接返回给上层
```go
func GiveIphone(userId int) error {
	// db...
	user, err := GetUser(db, userId)
	if err != nil {
		return err  // service 层 error 直接返回
	}

	err = SaveIphoneRecord(user.ID)
	if err != nil {
		return err  // 直接返回
	}

	err = SendGiftNotify(user.Email)
	if err != nil {
		return nil  // 直接返回
	}

	return nil
}
```

api 层需要处理 error，api 层是程序的入口，统一在这里打日志
```go
func HandleIphone(w http.ResponseWriter, r *http.Request) {
	userId := r.Form.Get("user_id")
	err := GiveIphone(userId)
	if err != nil {
		// handle error，入口处一次打印详细堆栈信息
		log.Printf("%+v", err)
		WriteResponse(w, FAILED)
		return
	}
	WriteResponse(w, SUCCESS)
}
```
