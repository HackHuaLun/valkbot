package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"log"
	"path/filepath"
	"strings"
	"time"
)

type Item struct {
	XPath string
	Title string
}

var items []Item

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

	// 启动游戏
	//if err := chromedp.Run(ctx, LaunchGame()); err != nil {
	//    log.Fatal(err)
	//    return
	//}

	// 获取信息
	for {

		if err := chromedp.Run(ctx, GetInfo()); err != nil {
			log.Fatal(err)
			return
		}
		//time.Sleep(1* time.Second)
		if err := chromedp.Run(ctx, chromedp.Click(items[0].XPath)); err != nil {
			log.Fatal(err)
			return
		}

		log.Println(items)
		time.Sleep(time.Second * 5)
	}

	//time.Sleep(2 * time.Second)
	//targets, _ := chromedp.Targets(ctx)
	//for _, tg := range targets {
	//    if strings.HasPrefix(tg.URL, "chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/notification.html") {
	//        newCtx, _ := chromedp.NewContext(ctx, chromedp.WithTargetID(tg.TargetID))
	//        sel := `#app-content`
	//        var tit string
	//        var val string
	//        if err := chromedp.Run(newCtx, chromedp.Tasks{
	//            chromedp.Text(sel, &val),
	//            chromedp.Title(&tit),
	//        }); err != nil {
	//            log.Fatalln(err)
	//        }
	//
	//        log.Println(tit)
	//    }
	//}

	chromedp.ListenTarget(ctx, func(ev interface{}) {

		//if reflect.TypeOf(ev).String() != "*cdproto.Message" {
		//    log.Println(reflect.TypeOf(ev))
		//}

		if tg, ok := ev.(*target.EventTargetInfoChanged); ok {
			//t := page.HandleJavaScriptDialog(false)
			log.Println(tg.TargetInfo.URL)

			if strings.HasPrefix(tg.TargetInfo.URL, "chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/notification.html#connect") {
				//log.Println("metamask link")
				//sel := `#app-content > div > div.main-container-wrapper > div > div.permissions-connect-choose-account > div.permissions-connect-choose-account__footer-container > div.permissions-connect-choose-account__bottom-buttons > button.button.btn-primary`
				//var tit string
				var tit string
				var text string

				newCtx, _ := chromedp.NewContext(ctx, chromedp.WithTargetID(tg.TargetInfo.TargetID))
				if err := chromedp.Run(newCtx, chromedp.Tasks{
					//chromedp.Title(&tit),
					chromedp.WaitVisible(`//button[text()='下一步']`),
					chromedp.Click(`//button[text()='下一步']`, chromedp.BySearch),
				}); err != nil {

				}

				log.Println(22, tit, text)
			}

		}
	})

	select {}
	cancel()
}

// LaunchGame 启动游戏
func LaunchGame() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/home.html"),
		chromedp.WaitVisible(`#password`),
		chromedp.SendKeys(`#password`, "11111111"),
		chromedp.Sleep(time.Second),
		chromedp.Click(`#app-content > div > div.main-container-wrapper > div > div > button`, chromedp.ByID),
		chromedp.Navigate("https://game.playvalkyr.io/fighters"),
		chromedp.WaitVisible(`#app > div.warning > div > button.el-button.el-button--block.el-button--success.el-button--large`),
		chromedp.Click(`#app > div.warning > div > button.el-button.el-button--block.el-button--success.el-button--large`, chromedp.ByQuery),
	}
}

func GetInfo() chromedp.Tasks {
	var nodes []*cdp.Node
	//var nodeId []cdp.NodeID
	return chromedp.Tasks{
		chromedp.Navigate("https://v2ex.com/"),
		chromedp.WaitVisible(`#Main`),
		chromedp.Nodes(`#Main .item_title`, &nodes, chromedp.BySearch),
		chromedp.ActionFunc(func(c context.Context) error {

			for i, node := range nodes {
				html, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(c)
				if err != nil {
					log.Println(err)
					continue
				}

				doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
				if err != nil {
					log.Println(err)
					continue
				}

				items = append(items, Item{
					XPath: node.FullXPath(),
					Title: doc.Find(`.topic-link`).Text(),
				})
				log.Println(i, node.FullXPathByID(), doc.Find(`.topic-link`).Text())
			}

			return nil
		}),
	}
}

func GetNodeInfo(i int) chromedp.Tasks {

	var text string
	return chromedp.Tasks{chromedp.WaitVisible(`#characters`),
		chromedp.Text(fmt.Sprintf(`#characters div.character-info .level:eq(%d)`, i), &text, chromedp.BySearch),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println(text)
			return nil
		}),
	}
}
