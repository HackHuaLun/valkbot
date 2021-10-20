package main

import (
    "context"
    "github.com/chromedp/chromedp"
    "log"
    "path/filepath"
    "strings"
    "time"
)

func main() {
    // industry fence famous level unknown hungry about chief divide wait critic hockey
    path, _ := filepath.Abs("./data")
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
            chromedp.UserDataDir(path),
            chromedp.Flag("disable-sync", false),
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
    
    time.Sleep(2 * time.Second)
    targets, _ := chromedp.Targets(ctx)
    for _, tg := range targets {
        if strings.HasPrefix(tg.URL, "chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/notification.html") {
            newCtx, _ := chromedp.NewContext(ctx, chromedp.WithTargetID(tg.TargetID))
            sel := `#app-content`
            var tit string
            var val string
            if err := chromedp.Run(newCtx, chromedp.Tasks{
                chromedp.ActionFunc(func(c context.Context) error {
                  chromedp.Text(sel, &val)
                  log.Println(22, val)
                  return nil
                }),
    
                chromedp.Text(sel, &val),
                chromedp.Title(&tit),
            }); err != nil {
                log.Fatalln(err)
            }
            
            log.Println(tit)
        }
    }
    
    //chromedp.ListenTarget(ctx, func(ev interface{}) {
    //
    //   //if reflect.TypeOf(ev).String() != "*cdproto.Message" {
    //   //    log.Println(reflect.TypeOf(ev))
    //   //}
    //
    //   if tg, ok := ev.(*target.EventTargetInfoChanged); ok {
    //       //t := page.HandleJavaScriptDialog(false)
    //       log.Println(tg.TargetInfo.URL)
    //
    //       if strings.HasPrefix(tg.TargetInfo.URL,"chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/notification.html#connect") {
    //           //log.Println("metamask link")
    //           //sel := `#app-content > div > div.main-container-wrapper > div > div.permissions-connect-choose-account > div.permissions-connect-choose-account__footer-container > div.permissions-connect-choose-account__bottom-buttons > button.button.btn-primary`
    //           //var tit string
    //           if err := chromedp.Run(ctx, chromedp.Tasks{
    //               chromedp.Sleep(time.Second*2),
    //               chromedp.ActionFunc(func(c context.Context) error {
    //                   log.Println(target.CloseTarget(tg.TargetInfo.TargetID).Do(ctx))
    //                   return nil
    //               }),
    //               //chromedp.WaitVisible(sel),
    //               //chromedp.Sleep(time.Second*2),
    //               chromedp.ActionFunc(func(c context.Context) error {
    //                  log.Println(33, tg.TargetInfo.Attached)
    //                  return nil
    //               }),
    //               //chromedp.Reload(),
    //               //chromedp.ActionFunc(func(c context.Context) error {
    //               //    var val string
    //               //    chromedp.Value(sel, &val)
    //               //
    //               //    log.Println(22, val, tit)
    //               //    return nil
    //               //}),
    //
    //           }); err != nil {
    //
    //           }
    //       }
    //
    //
    //   }
    //})
    
    select {}
    cancel()
}

// 自定义任务
func myTasks() chromedp.Tasks {
    return chromedp.Tasks{
        chromedp.Navigate("chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/home.html"),
        chromedp.WaitVisible(`#password`),
        chromedp.SendKeys(`#password`, "11111111"),
        chromedp.Sleep(time.Second),
        chromedp.Click(`#app-content > div > div.main-container-wrapper > div > div > button`, chromedp.ByID),
        chromedp.Navigate("https://v2.qkswap.io/"),
    }
}
