package app

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func Run(version string) error {
	err := run(version)
	if isFormCancelled(err) || errors.Is(err, flag.ErrHelp) {
		return nil
	}
	return err
}

func run(version string) error {
	flags := flag.NewFlagSet("badminton", flag.ContinueOnError)
	showVersion := flags.Bool("version", false, "显示版本")
	testDingTalk := flags.Bool("testdingtalk", false, "测试钉钉通知")
	testWindows := flags.Bool("testwin", false, "测试 Windows 弹窗")
	err := flags.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("不支持位置参数，请使用 --help 查看可用参数")
	}
	if *showVersion {
		fmt.Printf("BitBMTGrab %s\n", version)
		return nil
	}
	if *testDingTalk && *testWindows {
		return fmt.Errorf("每次只测试一种通知")
	}

	path, err := configPath()
	if err != nil {
		return err
	}
	settings, err := loadSettings(path)
	missing := errors.Is(err, os.ErrNotExist)
	if err != nil && !missing {
		return err
	}
	if *testDingTalk {
		err = sendDingTalkTestNotice(settings)
		if err != nil {
			return fmt.Errorf("钉钉通知测试失败: %w", err)
		}
		fmt.Println("钉钉通知测试发送成功")
		return nil
	}
	if *testWindows {
		err = showWindowsTestNotice(settings)
		if err != nil {
			return fmt.Errorf("Windows 通知测试失败: %w", err)
		}
		fmt.Println("Windows 通知测试显示成功")
		return nil
	}
	return runMenu(version, path, &settings, missing)
}
