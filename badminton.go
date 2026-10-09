package main

import (
	"badminton/api"
	"badminton/config"
	"badminton/models"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	maxBookingTargets = 5
	checkInterval     = 5 * time.Second
)

type bookingTarget struct {
	VenueName    string
	VenueFieldID string
	Scene        string
}

// formatAvailableVenues 格式化并输出场地可用信息
func formatAvailableVenues(venues map[string]*models.Venue, hours []*models.SportEventsHour, booked *models.SportScheduleBooked) {
	var venueKeys []string
	for key := range venues {
		venueKeys = append(venueKeys, key)
	}
	sort.Strings(venueKeys)

	timeSlots := make(map[string]map[string]bool)
	for _, hour := range hours {
		timeSlot := fmt.Sprintf("%s-%s", hour.BegintimeText, hour.EndtimeText)
		timeSlots[timeSlot] = make(map[string]bool)
		for _, venueKey := range venueKeys {
			// 构造场地-时段键
			bookingKey := fmt.Sprintf("%s-%d", venueKey, hour.Id)
			// 检查是否被预订（0表示未预订）
			timeSlots[timeSlot][venueKey] = booked.Data[bookingKey] == 0
		}
	}

	fmt.Println("\n场地预订情况:")
	var timeSlotKeys []string
	for timeSlot := range timeSlots {
		timeSlotKeys = append(timeSlotKeys, timeSlot)
	}
	sort.Strings(timeSlotKeys)

	fmt.Printf("%-15s", "时间段")
	for _, venueKey := range venueKeys {
		fmt.Printf("%-8s", venues[venueKey].Name)
	}
	fmt.Println()

	for _, timeSlot := range timeSlotKeys {
		fmt.Printf("%-15s", timeSlot)
		for _, venueKey := range venueKeys {
			if timeSlots[timeSlot][venueKey] {
				fmt.Printf("%-8s", "可用")
			} else {
				fmt.Printf("%-8s", "已订")
			}
		}
		fmt.Println()
	}
}

// getWeekday 根据日期字符串获取星期几（1-7）
func getWeekday(dateStr string) (int, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, fmt.Errorf("日期格式错误: %v", err)
	}

	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return weekday, nil
}

func normalizeBookingDate(input string, defaultYear int) (string, error) {
	if input == "" {
		return "", fmt.Errorf("日期不能为空")
	}

	if strings.Contains(input, "-") {
		date, err := time.Parse("2006-01-02", input)
		if err != nil {
			return "", fmt.Errorf("日期格式错误: %v", err)
		}
		return date.Format("2006-01-02"), nil
	}

	dateInput := input
	if defaultYear > 0 && len(input) == 4 {
		dateInput = fmt.Sprintf("%04d%s", defaultYear, input)
	}
	if len(dateInput) != 8 {
		if defaultYear > 0 {
			return "", fmt.Errorf("请按 MMDD 或 YYYYMMDD 输入日期")
		}
		return "", fmt.Errorf("请按 YYYYMMDD 输入日期")
	}

	date, err := time.Parse("20060102", dateInput)
	if err != nil {
		return "", fmt.Errorf("日期格式错误: %v", err)
	}
	return date.Format("2006-01-02"), nil
}

// findVenueIDByName 根据完整名称或名称末尾的编号查找场地ID
func findVenueIDByName(venues map[string]*models.Venue, name string) (string, error) {
	prefix, number, _, err := splitVenueNameNumber(name)
	if err == nil && prefix == "" {
		var matchedID string
		var matchedNames []string
		for id, venue := range venues {
			_, venueNumber, _, err := splitVenueNameNumber(venue.Name)
			if err == nil && venueNumber == number {
				matchedID = id
				matchedNames = append(matchedNames, venue.Name)
			}
		}
		switch len(matchedNames) {
		case 0:
			return "", fmt.Errorf("未找到编号为 %s 的场地", name)
		case 1:
			return matchedID, nil
		default:
			sort.Strings(matchedNames)
			return "", fmt.Errorf("场地编号 %s 对应多个场地（%s），请使用完整场地名称", name, strings.Join(matchedNames, ","))
		}
	}

	for id, venue := range venues {
		if venue.Name == name {
			return id, nil
		}
	}
	return "", fmt.Errorf("未找到名为 %s 的场地", name)
}

// findHourIDByTime 根据时间（小时）查找对应的时间段ID
func findHourIDByTime(hours []*models.SportEventsHour, hour int) (int, error) {
	targetTime := fmt.Sprintf("%02d:00", hour)

	for _, h := range hours {
		if h.BegintimeText == targetTime {
			return h.Id, nil
		}
	}
	return 0, fmt.Errorf("未找到开始时间为 %s 的时间段", targetTime)
}

func parseVenueNames(input string) ([]string, error) {
	if input == "" {
		return nil, fmt.Errorf("场地名称不能为空")
	}

	normalized := strings.ReplaceAll(input, "，", ",")
	parts := strings.Split(normalized, ",")
	names := make([]string, 0, maxBookingTargets)
	seen := make(map[string]bool)
	for _, part := range parts {
		if part == "" {
			return nil, fmt.Errorf("场地名称不能为空，请检查逗号是否多余")
		}
		if strings.TrimSpace(part) != part {
			return nil, fmt.Errorf("场地名称 %q 前后不要包含空格", part)
		}

		expandedNames, err := expandVenueNamePart(part)
		if err != nil {
			return nil, err
		}
		for _, name := range expandedNames {
			if seen[name] {
				return nil, fmt.Errorf("场地名称重复: %s", name)
			}
			seen[name] = true
			names = append(names, name)
			if len(names) > maxBookingTargets {
				return nil, fmt.Errorf("一次最多输入 %d 个场地", maxBookingTargets)
			}
		}
	}

	return names, nil
}

func expandVenueNamePart(part string) ([]string, error) {
	if !strings.Contains(part, "-") {
		return []string{part}, nil
	}

	rangeParts := strings.Split(part, "-")
	if len(rangeParts) != 2 || rangeParts[0] == "" || rangeParts[1] == "" {
		return nil, fmt.Errorf("场地范围格式错误: %s", part)
	}

	startPrefix, startNumber, startWidth, err := splitVenueNameNumber(rangeParts[0])
	if err != nil {
		return nil, fmt.Errorf("场地范围格式错误: %s", part)
	}
	endPrefix, endNumber, _, err := splitVenueNameNumber(rangeParts[1])
	if err != nil {
		return nil, fmt.Errorf("场地范围格式错误: %s", part)
	}
	if startPrefix != endPrefix {
		return nil, fmt.Errorf("场地范围前后格式不一致: %s", part)
	}
	if startNumber > endNumber {
		return nil, fmt.Errorf("场地范围不能倒序: %s", part)
	}
	if endNumber-startNumber >= maxBookingTargets {
		return nil, fmt.Errorf("一次最多输入 %d 个场地", maxBookingTargets)
	}

	count := endNumber - startNumber + 1
	names := make([]string, 0, count)
	for offset := 0; offset < count; offset++ {
		names = append(names, fmt.Sprintf("%s%0*d", startPrefix, startWidth, startNumber+offset))
	}
	return names, nil
}

func splitVenueNameNumber(name string) (string, int, int, error) {
	numberStart := len(name)
	for numberStart > 0 && name[numberStart-1] >= '0' && name[numberStart-1] <= '9' {
		numberStart--
	}
	if numberStart == len(name) {
		return "", 0, 0, fmt.Errorf("场地名称缺少数字: %s", name)
	}

	numberText := name[numberStart:]
	number, err := strconv.Atoi(numberText)
	if err != nil {
		return "", 0, 0, err
	}
	return name[:numberStart], number, len(numberText), nil
}

func buildBookingTargets(venues map[string]*models.Venue, venueNames []string, dateStr string, hourID int) ([]bookingTarget, error) {
	targets := make([]bookingTarget, 0, len(venueNames))
	seen := make(map[string]bool)
	for _, venueName := range venueNames {
		venueFieldID, err := findVenueIDByName(venues, venueName)
		if err != nil {
			return nil, err
		}
		actualName := venues[venueFieldID].Name
		if seen[venueFieldID] {
			return nil, fmt.Errorf("场地名称重复: %s", actualName)
		}
		seen[venueFieldID] = true

		targets = append(targets, bookingTarget{
			VenueName:    actualName,
			VenueFieldID: venueFieldID,
			Scene:        fmt.Sprintf("[{\"day\":\"%s\",\"fields\":{\"%s\":[%d]}}]", dateStr, venueFieldID, hourID),
		})
	}

	return targets, nil
}

// waitUntilNextAvailableTime 检查当前时间是否在允许下单的时间范围内，如果不在则等待到下一个可用时间
func waitUntilNextAvailableTime() {
	for {
		now := time.Now()
		currentTime := now.Format("15:04:05")

		if currentTime >= "07:00:00" && currentTime <= "23:30:00" {
			return
		}

		nextAvailable := now
		if currentTime > "23:30:00" {
			nextAvailable = now.Add(24 * time.Hour)
		}
		nextAvailable = time.Date(
			nextAvailable.Year(),
			nextAvailable.Month(),
			nextAvailable.Day(),
			7, 0, 0, 0,
			nextAvailable.Location(),
		)

		waitDuration := nextAvailable.Sub(now)

		fmt.Printf("\r当前时间: %s, 等待时间: %s",
			currentTime,
			waitDuration.Round(time.Second),
		)

		time.Sleep(time.Second)
	}
}

func notifyOrderDetected(venueID int, target bookingTarget, dateStr string, hour int, order models.OrderItem) {
	err := sendOrderNotice(venueID, target.VenueName, dateStr, hour, order.ID, order.Amount)
	if err != nil {
		fmt.Printf("[%s] 钉钉通知失败: %v\n", target.VenueName, err)
	}
	err = showWindowsOrderNotice(venueID, target.VenueName, dateStr, hour, order.ID, order.Amount)
	if err != nil {
		fmt.Printf("[%s] Windows 通知失败: %v\n", target.VenueName, err)
	}
	fmt.Printf("[%s] 已检测到订单 %d，正在等待支付...\n", target.VenueName, order.ID)
}

func processBookingTarget(venueID int, target bookingTarget, dateStr string, hour int, headers map[string]string) (bool, error) {
	check, err := api.CheckSportSchedule(venueID, target.Scene, headers)
	if err != nil {
		if strings.Contains(err.Error(), "预约中") {
			fmt.Printf("[%s] 检查订单失败: %v，约5秒后继续检查...\n", target.VenueName, err)
			return false, nil
		}
		if api.IsRetryableError(err) {
			fmt.Printf("[%s] 检查请求暂时失败: %v，约5秒后继续检查...\n", target.VenueName, err)
			return false, nil
		}

		return false, fmt.Errorf("[%s] 检查订单失败: %v", target.VenueName, err)
	}

	fmt.Printf("\n[%s] 检查通过，准备进入下单流程...\n", target.VenueName)
	_, err = api.SubmitOrder(venueID, check.Data.TotalAmount, target.Scene, headers)
	if err != nil {
		fmt.Printf("[%s] 提交订单失败: %v，下轮回到检查阶段...\n", target.VenueName, err)
		return false, nil
	}

	fmt.Printf("[%s] 订单提交成功，开始确认订单...\n", target.VenueName)
	for {
		orderList, err := api.GetOrderList("makeappointment", "created", "", 1, 20, headers)
		if err != nil {
			fmt.Printf("[%s] 获取订单列表失败: %v，5秒后继续确认...\n", target.VenueName, err)
			time.Sleep(checkInterval)
			continue
		}

		for _, order := range orderList.Data.List {
			if order.SportEventsID == venueID {
				notifyOrderDetected(venueID, target, dateStr, hour, order)
				return true, nil
			}
		}

		fmt.Printf("[%s] 暂未检测到未支付订单，5秒后继续确认...\n", target.VenueName)
		time.Sleep(checkInterval)
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--testdingtalk" {
		if err := sendDingTalkTestNotice(); err != nil {
			fmt.Printf("钉钉通知测试失败: %v\n", err)
			return
		}
		fmt.Println("钉钉通知测试发送成功")
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "--testwin" {
		if err := showWindowsTestNotice(); err != nil {
			fmt.Printf("Windows 通知测试失败: %v\n", err)
			return
		}
		fmt.Println("Windows 通知测试显示成功")
		return
	}

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
	}

	venueID := config.DefaultVenueID
	if venueID > 0 {
		fmt.Printf("使用配置场馆ID: %d\n", venueID)
	} else {
		fmt.Print("请输入场馆ID（25-良乡，51-中关村，54-中关村早晚场）: ")
		fmt.Scanln(&venueID)
	}

	venues, err := api.GetSportEventsField(venueID, headers)
	if err != nil {
		fmt.Printf("获取场馆信息失败: %v\n", err)
		return
	}

	hours, err := api.GetSportEventsHour(venueID, headers)
	if err != nil {
		fmt.Printf("获取时间段信息失败: %v\n", err)
		return
	}

	if config.DefaultYear > 0 {
		fmt.Printf("请输入预订日期（格式：MMDD，年份使用配置 %d）: ", config.DefaultYear)
	} else {
		fmt.Print("请输入预订日期（格式：YYYYMMDD）: ")
	}
	var dateInput string
	fmt.Scanln(&dateInput)
	dateStr, err := normalizeBookingDate(dateInput, config.DefaultYear)
	if err != nil {
		fmt.Printf("日期输入错误: %v\n", err)
		return
	}

	booked, err := api.GetSportScheduleBooked(venueID, dateStr, headers)
	if err != nil {
		fmt.Printf("获取场地预订信息失败: %v\n", err)
		return
	}

	formatAvailableVenues(venues, hours, booked)

	week, err := getWeekday(dateStr)
	if err != nil {
		fmt.Printf("计算星期失败: %v\n", err)
		return
	}

	_, err = api.GetSportEventsPrice(venueID, week, dateStr, headers)
	if err != nil {
		fmt.Printf("获取价格信息失败: %v\n", err)
		return
	}

	fmt.Printf("请输入场地名称或编号（多个用逗号分隔，最多%d个）: ", maxBookingTargets)
	var venueNamesInput string
	if _, err = fmt.Scanln(&venueNamesInput); err != nil {
		fmt.Println("场地名称输入错误: 请使用逗号分隔且不要输入空格，例如 1,2 或 场1，场2")
		return
	}

	venueNames, err := parseVenueNames(venueNamesInput)
	if err != nil {
		fmt.Printf("场地名称输入错误: %v\n", err)
		return
	}

	fmt.Print("请输入预订时间（小时，例如9表示09:00开始的时段）: ")
	var hour int
	fmt.Scanln(&hour)

	hourID, err := findHourIDByTime(hours, hour)
	if err != nil {
		fmt.Printf("查找时间段失败: %v\n", err)
		return
	}

	targets, err := buildBookingTargets(venues, venueNames, dateStr, hourID)
	if err != nil {
		fmt.Printf("查找场地失败: %v\n", err)
		return
	}

	targetInterval := checkInterval / time.Duration(len(targets))

	for {
		waitUntilNextAvailableTime()

		fmt.Printf("正在检查场馆ID %d 的订单信息...\n", venueID)

		for _, target := range targets {
			done, err := processBookingTarget(venueID, target, dateStr, hour, headers)
			if err != nil {
				fmt.Println(err)
				return
			}
			if done {
				return
			}

			time.Sleep(targetInterval)
		}
	}

}
