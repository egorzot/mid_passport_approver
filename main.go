package main

import (
	"context"
	"fmt"
	"github.com/2captcha/2captcha-go"
	"github.com/chromedp/chromedp"
	"github.com/go-co-op/gocron"
	"log"
	"math/rand"
	"os"
	"passport/config"
	"strings"
	"time"
)

var twoCaptchaApiClient *api2captcha.Client

func main() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(86400 + 60*5).Seconds().Do(func() {
		makeMyPassport()
	})

	s.StartBlocking()
}

func makeMyPassport() {
	fmt.Println("Process started")

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", true))...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://q.midpass.ru/"),
		chromedp.Sleep(5*time.Second),
		chromedp.WaitVisible(`.SelInp`),
		chromedp.SetValue(`select.SelInp`, config.GetCountry(), chromedp.ByQuery),
		chromedp.Sleep(randomDuration()),
		chromedp.SetValue(`.wrap:nth-child(3) > .register_form > select.SelInp`, config.GetInstitution(), chromedp.ByQuery),
		chromedp.Sleep(randomDuration()),
		chromedp.SendKeys("#Email", config.GetLogin(), chromedp.ByID),
		chromedp.Sleep(randomDuration()),
		chromedp.SendKeys("#Password", config.GetPassword(), chromedp.ByID),
		chromedp.Sleep(randomDuration()),
	); err != nil && !strings.Contains(err.Error(), "net::ERR_ABORTED") {
		handleError(ctx, err)
		return
	}

	twoCaptchaApiClient = api2captcha.NewClient(config.GetApiKey())

	code := resolveCptch(ctx, "//img[@id='imgCaptcha']")
	if err := chromedp.Run(ctx,
		chromedp.SendKeys("#Captcha", code, chromedp.ByID),
		chromedp.Sleep(randomDuration()),
		chromedp.Click(".registerForm:nth-child(16) > button", chromedp.ByQuery),
		chromedp.Sleep(10*time.Second),
		chromedp.Click("#topNav > li:nth-child(3)", chromedp.ByQuery),
		chromedp.Sleep(randomDuration()),
		chromedp.Click("#datagrid-row-r1-1-0 > td > div > input", chromedp.ByQuery),
		chromedp.Sleep(randomDuration()),
		chromedp.Click("#confirmAppointments", chromedp.ByQuery),
		chromedp.Sleep(5*time.Second),
		chromedp.WaitVisible(`#captchaValue`),
	); err != nil && !strings.Contains(err.Error(), "net::ERR_ABORTED") {
		handleError(ctx, err)
		return
	}

	code2 := resolveCptch(ctx, "//img[@id='imgCaptcha']")

	if err := chromedp.Run(ctx,

		chromedp.SendKeys("#captchaValue", code2, chromedp.ByID),
		chromedp.Sleep(randomDuration()),
		chromedp.Click(".dialog-button > a", chromedp.ByQuery),
		chromedp.Sleep(randomDuration()),
	); err != nil && !strings.Contains(err.Error(), "net::ERR_ABORTED") {
		handleError(ctx, err)
		return
	}
	time.Sleep(5 * time.Second)
}

func resolveCptch(ctx context.Context, querySelector string) string {
	var buf []byte
	if err := chromedp.Run(ctx, chromedp.Screenshot(querySelector, &buf, chromedp.NodeVisible)); err != nil {
		log.Fatal(err)
	}

	fileName := "cptch" + time.Now().Format(time.RFC850) + ".jpg"

	if err := os.WriteFile("cptch/"+fileName, buf, 0o644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove("cptch/" + fileName)

	capt := api2captcha.Normal{
		File:          "cptch/" + fileName,
		CaseSensitive: true,
	}

	code, err := twoCaptchaApiClient.Solve(capt.ToRequest())
	if err != nil {
		if err == api2captcha.ErrTimeout {
			log.Fatal("Timeout")
		} else if err == api2captcha.ErrApi {
			log.Fatal("API error")
		} else if err == api2captcha.ErrNetwork {
			log.Fatal("Network error")
		} else {
			log.Fatal(err)
		}
	}

	fmt.Printf("Captcha code is " + code)

	return code
}

func randomDuration() time.Duration {
	rand.Seed(time.Now().UnixNano())
	min := 2
	max := 5
	return (time.Duration(rand.Intn(max-min+1) + min)) * time.Second
}

func handleError(ctx context.Context, error error) {
	var buf []byte
	if err := chromedp.Run(ctx, chromedp.FullScreenshot(&buf, 90)); err != nil {
		log.Fatal(err)
	}

	filePath := "full" + time.Now().Format(time.RFC850) + ".png"

	if err := os.WriteFile("screenshots/"+filePath, buf, 0o644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Error: " + error.Error() + ". The full screenshot was saved to the " + filePath)
}
