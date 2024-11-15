# Parallel
golang 协程并行库，可以指定并发数量，可以当做简单的队列使用

## 使用 - usage
1. 普通用法
    ```go

    func TestParallel(t *testing.T) {
        // 最多 10 个协程同时运行
        p := parallel.NewParallel(10)
    
        p.Add(func() interface{} {
            return "执行了"
        })
        
        p.Add(func() interface{} {
            panic(errors.New("报错了"))
        })
    
        fmt.Println(p.Wait())
        //会输出 map[0:执行了 1:报错了]
    }
    ```
2. 简单的队列
    ```go

    // 测试优雅退出
    func TestParallelGracefulStop(t *testing.T) {
        p := parallel.NewParallel(2)
    
        go func() {
            i := 0
            for ; i < 10; i++ {
                (func(i int) {
                    if i >= 5 {
                        p.GracefulStop()
                    }
                    result := p.Add(func() interface{} {
                        time.Sleep(time.Second)
                        fmt.Printf("每隔1秒执行一次 %d \n", i)
                        return nil
                    })
    
                    fmt.Printf("添加结果: %v, i: %d\n", result, i)
                })(i)
            }
        }()
    
        fmt.Println(p.Listen())
    }
    ```

