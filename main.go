package main

import (
    "context"
    "github.com/chromedp/cdproto/target"
    "github.com/chromedp/chromedp"
    "log"
    "strings"
)

func main()  {
    // chromdp依赖context上限传递参数
    ctx, _ := chromedp.NewExecAllocator(
        context.Background(),
        
        // 以默认配置的数组为基础，覆写headless参数
        // 当然也可以根据自己的需要进行修改，这个flag是浏览器的设置
        append(
            chromedp.DefaultExecAllocatorOptions[:],
            chromedp.Flag("headless", false),
            chromedp.Flag("disable-gpu", false),
            chromedp.Flag("enable-automation", false),
            chromedp.Flag("disable-extensions", false),
            chromedp.UserDataDir("./data"),
            //chromedp.Flag("disable-sync", false),
        )...,
    )
    
    // 创建新的chromedp上下文对象，超时时间的设置不分先后
    // 注意第二个返回的参数是cancel()
    ctx, cancel := chromedp.NewContext(
        ctx,
        // 设置日志方法
        chromedp.WithLogf(log.Printf),
    )
    
    // 执行我们自定义的任务 - myTasks函数在第4步
    if err := chromedp.Run(ctx, myTasks()); err != nil {
        log.Fatal(err)
        return
    }
    
    
    chromedp.ListenTarget(ctx, func(ev interface{}) {
        //log.Println(reflect.TypeOf(ev))
        if tg, ok := ev.(*target.EventTargetInfoChanged); ok {
            //t := page.HandleJavaScriptDialog(false)
            log.Println(tg.TargetInfo.Type)
    
            if strings.HasPrefix(tg.TargetInfo.URL,"chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/notification.html#connect") {
                log.Println("metamask link")
                go func() {
                    sel := `#app-content > div > div.main-container-wrapper > div > div.permissions-connect-choose-account > div.permissions-connect-choose-account__footer-container > div.permissions-connect-choose-account__bottom-buttons > button.button.btn-primary`
                    if err := chromedp.Run(ctx, chromedp.Tasks{
                        chromedp.WaitVisible(sel),
                        chromedp.ActionFunc(func(ctx context.Context) error {
                            var val string
                            chromedp.Value(sel, &val)
                            
                            log.Println(22, val)
                            return nil
                        }),
                    
                    }); err != nil {
            
                    }
        
                    // ok
                }()
            }
            
            
        }
    })
    
    
    
    select {}
    cancel()
}

// 自定义任务
func myTasks() chromedp.Tasks {
    return chromedp.Tasks{
        chromedp.Navigate("chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/home.html"),
        chromedp.WaitVisible(`#password`),
        chromedp.SendKeys(`#password`, "11111111"),
        chromedp.Click(`#app-content > div > div.main-container-wrapper > div > div > button`, chromedp.ByID),
        chromedp.Navigate("https://v2.qkswap.io/"),
    }
}
