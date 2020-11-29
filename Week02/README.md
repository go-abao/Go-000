### Week2学习要点： ERROR的错误处理注意事项

1. 可以使用`errors.Cause`获取根因来与`sentinel error`进行判定

2. go1.13新特性 `errors.Is(与sentinel error进行判断)` 和 `errors.As(用处：断言)` 的使用, 不要再直接与sentinel error进行等值判断
   - `errors.Is` 会与sentinel error进行等值判断, 如果判断没过, 则尝试通过Unwrap获取error的根因
   - `errors.As` 可以将err对象转换为自定义的error类型

3. Error只处理1一次, 打日志(最顶层)或者向上返回, 不要操作两次, 打印日志可以通过使用`%+v`谓词把堆栈信息打印出来
   - 如果不处理error, 那么要使用Wrap(f)添加一些上下文信息向上返回（不需要把整个response打印出来, 因为它的内容太多）
   - 如果处理了error(打日志 降级 or 其他逻辑), 那么不要再将error向上返回

4. 首次产生错误的时候对error进行warp 
   - `首次`是指在业务代码与标准库, 第三方库等交互以及自己的代码产生错误时. 
   - 自己的代码中返回错误信息使用`errors.New`或`errors.Errorf`

5. 如果在自己的代码中调用其他函数, 要直接返回, 不要再次进行warp

6. 可以使用github.com/pkg/errors包来代替标准库的errors包, 该包可以通过`errors.Wrap(f)`和`errors.WithMessage`来添加上下文信息

7.  `errors.Wrap(f)`保存了error的堆栈信息, `errors.WithMessage`不保存堆栈信息

8. kit(基础库)或标准库不应该对error进行Wrap, 只能返回根因, 业务层代码可以进行Wrap

9. error的Unwrap方法可以将根因返回

10. 清理mod缓存可以使用 `go clean --modcache`

### Week02 作业题目：

问：我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个error抛给上层？为什么，应该怎么做请写出代码？

答：自身业务逻辑出错不应该，详情见代码。