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
	"strconv"
	"strings"
	"time"
)

type Item struct {
	XPath string
	Id string
	Lv string
	CP int
	SP int
}


type Account struct {
	XPath string
	Name string
	Amount string
}

var items []Item
var accounts []Account
var accountLoading bool

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

	// 注意第二个返回的参数是cancel()
	ctx, cancel := chromedp.NewContext(
		ctx,
		// 设置日志方法
		chromedp.WithLogf(log.Printf),
	)

	// 加载账户
	if err := chromedp.Run(ctx, Home()); err != nil {
		log.Fatal(err)
		return
	}
	
	for _, account := range accounts {
		log.Println(fmt.Sprintf("账户: %s, 余额: %s", account.Name, account.Amount))
	}

	// 默认1账户
	ChangeAccount(ctx, accounts[1])
	
	go func(c context.Context) {
		// 获取信息
		for {
			time.Sleep(time.Second * 3)
			if accountLoading {
				continue
			}

			log.Println("刷新信息:")
			if err := chromedp.Run(c, GetInfo()); err != nil {
				log.Fatal(err)
				return
			}
			
			for _, item := range items {
				log.Println(fmt.Sprintf("%s, CP: %d, SP: %d", item.Lv, item.CP, item.SP))
			}
			
			ChangeAccount(ctx, accounts[0])
		}
	}(ctx)
	
	
	select {}
	
	time.Sleep(time.Second * 1)
	item := items[1]

	if err := chromedp.Run(ctx, SelectHero(item.Id)); err != nil {
		log.Fatal(err)
		return
	}

	if err := chromedp.Run(ctx, Attach(item.Id)); err != nil {
		log.Fatal(err)
		return
	}

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
		chromedp.Navigate("https://game.playvalkyr.io/fighters"),
		chromedp.WaitVisible(`#app > div.warning`),
		chromedp.Click(`#app div.warning button.el-button.el-button--block.el-button--success.el-button--large`, chromedp.ByQuery),
		chromedp.WaitVisible(`#characters`),
	}
}

func GetInfo() chromedp.Tasks {
	var nodes []*cdp.Node
	return chromedp.Tasks{
		chromedp.Nodes(`#characters .el-tabs__item`, &nodes, chromedp.BySearch),
		chromedp.ActionFunc(func(c context.Context) error {

			for _, node := range nodes {
				
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

				cp := doc.Find(`.cp`).Text()
				cp = strings.ReplaceAll(cp, "CP:", "")
				cp = strings.ReplaceAll(cp, ",", "")
				cp = strings.TrimSpace(cp)
				cpInt,_ := strconv.Atoi(cp)
				
				sp := doc.Find(`.pool .el-progress__text`).Text()
				sp = strings.Split(sp, "/")[0]
				spInt,_ := strconv.Atoi(sp)
				
				if cpInt > 0 {
					items = append(items, Item{
						XPath: node.FullXPath(),
						Id: node.AttributeValue("aria-controls"),
						Lv:    doc.Find(`.level`).Text(),
						CP:    cpInt,
						SP:    spInt,
					})
				}
			}

			return nil
		}),
	}
}

// SelectHero 选择英雄
func SelectHero(id string) chromedp.Tasks {
	sel := fmt.Sprintf(`#characters div[aria-controls=%s] div.character-info`, id)
	return chromedp.Tasks{
		chromedp.Click(sel, chromedp.ByID),
	}
}

// Attach 攻击
func Attach(id string) chromedp.Tasks {
	sel := fmt.Sprintf(`#characters div[id=%s] div.text-xl.cp > button`, id)
	return chromedp.Tasks{
		chromedp.Click(sel, chromedp.ByID),
	}
}

func Home() chromedp.Tasks {
	var nodes []*cdp.Node
	return chromedp.Tasks{
		chromedp.Sleep(time.Second),
		chromedp.Navigate("chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/home.html"),
		chromedp.Sleep(time.Second),
		chromedp.WaitNotPresent(`#app-content .loading-overlay`),
		chromedp.SendKeys(`#password`, "11111111"),
		chromedp.Click(`#app-content > div > div.main-container-wrapper > div > div > button`, chromedp.ByID),
		
		chromedp.WaitVisible(`#app-content div.account-menu__icon`, chromedp.ByID),
		chromedp.Click(`#app-content div.account-menu__icon`, chromedp.ByID),
		
		chromedp.WaitVisible(`.account-menu`),
		chromedp.Nodes(`.account-menu .account-menu__account`, &nodes, chromedp.BySearch),
		
		chromedp.ActionFunc(func(c context.Context) error {
			
			for _, node := range nodes {
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
				
				name := doc.Find(`.account-menu__name`).Text()
				amount := doc.Find(`.currency-display-component__text`).Text()
				
				accounts = append(accounts, Account{
					XPath:  node.FullXPath(),
					Name:   name,
					Amount: amount,
				})
			}
			
			return nil
		}),
	}
}

// ChangeAccount 攻击
func ChangeAccount(c context.Context, account Account) {
	accountLoading = true
	log.Println("选择账户: "+ account.Name)
	
	if err := chromedp.Run(c, chromedp.Tasks{
		chromedp.Navigate("chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/home.html"),
		chromedp.WaitVisible(`#app-content div.account-menu__icon`, chromedp.ByID),
		chromedp.Click(`#app-content div.account-menu__icon`, chromedp.ByID),
		chromedp.WaitVisible(`.account-menu`),
		chromedp.Click(account.XPath),
		LaunchGame(),
		chromedp.ActionFunc(func(ctx context.Context) error {
			accountLoading = false
			return nil
		}),
	}); err != nil {
		log.Fatal(err)
		return
	}
}