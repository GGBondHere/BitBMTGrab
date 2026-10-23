package app

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	maxBookingTargets  = 5
	maxNetworkFailures = 5
	checkInterval      = 5 * time.Second
)

type bookingTarget struct {
	VenueName    string
	VenueFieldID string
	HourID       int
	Scene        string
}

// formatAvailableVenues 格式化并输出场地可用信息
func formatAvailableVenues(venues map[string]*Venue, hours []*SportEventsHour, booked *SportScheduleBooked) {
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
func findVenueIDByName(venues map[string]*Venue, name string) (string, error) {
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
func findHourIDByTime(hours []*SportEventsHour, hour int) (int, error) {
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

func buildBookingTargets(venues map[string]*Venue, venueNames []string, dateStr string, hourID int) ([]bookingTarget, error) {
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
			HourID:       hourID,
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

func matchesBookingTarget(order OrderItem, venueID int, target bookingTarget, dateStr string) bool {
	if order.SportEventsID != venueID {
		return false
	}

	scenes := order.Config.Scene
	if len(scenes) == 0 && order.ConfigJSON != "" {
		var details OrderConfig
		err := json.Unmarshal([]byte(order.ConfigJSON), &details)
		if err != nil {
			return false
		}
		scenes = details.Scene
	}
	for _, scene := range scenes {
		if scene.Day != dateStr {
			continue
		}
		for _, hourID := range scene.Fields[target.VenueFieldID] {
			if hourID == target.HourID {
				return true
			}
		}
	}
	return false
}

func notifyBookingSuccess(settings Settings, venueID int, target bookingTarget, dateStr string, hour int, order OrderItem) {
	if strings.TrimSpace(settings.DingTalkWebhook) != "" {
		err := sendOrderNotice(settings, venueID, target.VenueName, dateStr, hour, order.ID, order.Amount)
		if err != nil {
			fmt.Printf("[%s] 钉钉通知失败: %v\n", target.VenueName, err)
		} else {
			fmt.Printf("[%s] 钉钉通知发送成功\n", target.VenueName)
		}
	}
	if settings.EnableWindowsNotification {
		err := showWindowsOrderNotice(settings, venueID, target.VenueName, dateStr, hour, order.ID, order.Amount)
		if err != nil {
			fmt.Printf("[%s] Windows 通知失败: %v\n", target.VenueName, err)
		} else {
			fmt.Printf("[%s] Windows 弹窗已启动\n", target.VenueName)
		}
	}
}

func processBookingTarget(settings Settings, venueID int, target bookingTarget, dateStr string, hour int, headers map[string]string, networkFailures *int) (bool, error) {
	check, err := CheckSportSchedule(venueID, target.Scene, headers)
	if err != nil {
		if IsRetryableError(err) {
			*networkFailures += 1
			if *networkFailures >= maxNetworkFailures {
				return false, fmt.Errorf("[%s] 检查请求连续失败 %d 次，停止抢票: %w", target.VenueName, *networkFailures, err)
			}
			fmt.Printf("[%s] 检查请求暂时失败（连续 %d/%d 次）: %v，下轮继续检查...\n", target.VenueName, *networkFailures, maxNetworkFailures, err)
			return false, nil
		}
		*networkFailures = 0
		if strings.Contains(err.Error(), "预约中") {
			fmt.Printf("[%s] 检查订单失败: %v，下轮继续检查...\n", target.VenueName, err)
			return false, nil
		}

		return false, fmt.Errorf("[%s] 检查订单失败: %v", target.VenueName, err)
	}
	*networkFailures = 0

	fmt.Printf("\n[%s] 检查通过，准备进入下单流程...\n", target.VenueName)
	_, err = SubmitOrder(venueID, check.Data.TotalAmount, target.Scene, settings.OpenId, headers)
	if err != nil {
		fmt.Printf("[%s] 提交订单失败: %v，下轮回到检查阶段...\n", target.VenueName, err)
		return false, nil
	}

	fmt.Printf("[%s] 订单提交成功，开始确认订单...\n", target.VenueName)
	for {
		orderList, err := GetOrderList("makeappointment", "created", "", 1, 20, headers)
		if err != nil {
			fmt.Printf("[%s] 获取订单列表失败: %v，5秒后继续确认...\n", target.VenueName, err)
			time.Sleep(checkInterval)
			continue
		}

		for _, order := range orderList.Data.List {
			if matchesBookingTarget(order, venueID, target, dateStr) {
				fmt.Printf("[%s] 已确认目标订单 %d，抢票成功，请尽快支付。\n", target.VenueName, order.ID)
				notifyBookingSuccess(settings, venueID, target, dateStr, hour, order)
				return true, nil
			}
		}

		fmt.Printf("[%s] 暂未确认目标未支付订单，5秒后继续确认...\n", target.VenueName)
		time.Sleep(checkInterval)
	}
}

func runBooking(settings Settings) error {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
		"token":        settings.ApiToken,
	}

	venueID := settings.DefaultVenueID
	if venueID > 0 {
		fmt.Printf("使用配置场馆ID: %d\n", venueID)
	} else {
		var venueInput string
		err := inputForm("场馆 ID（25-良乡，51-中关村，54-中关村早晚场）", &venueInput, positiveNumber)
		if err != nil {
			return err
		}
		venueID, _ = strconv.Atoi(venueInput)
	}

	venues, err := GetSportEventsField(venueID, headers)
	if err != nil {
		return fmt.Errorf("获取场馆信息失败: %w", err)
	}

	hours, err := GetSportEventsHour(venueID, headers)
	if err != nil {
		return fmt.Errorf("获取时间段信息失败: %w", err)
	}

	dateTitle := "预约日期（YYYYMMDD）"
	if settings.DefaultYear > 0 {
		dateTitle = fmt.Sprintf("预约日期（MMDD，年份 %d）", settings.DefaultYear)
	}
	var dateInput string
	err = inputForm(dateTitle, &dateInput, func(value string) error {
		_, err := normalizeBookingDate(value, settings.DefaultYear)
		return err
	})
	if err != nil {
		return err
	}
	dateStr, _ := normalizeBookingDate(dateInput, settings.DefaultYear)

	booked, err := GetSportScheduleBooked(venueID, dateStr, headers)
	if err != nil {
		return fmt.Errorf("获取场地预订信息失败: %w", err)
	}

	formatAvailableVenues(venues, hours, booked)

	week, err := getWeekday(dateStr)
	if err != nil {
		return err
	}

	_, err = GetSportEventsPrice(venueID, week, dateStr, headers)
	if err != nil {
		return fmt.Errorf("获取价格信息失败: %w", err)
	}

	var venueNamesInput string
	err = inputForm(fmt.Sprintf("场地名称或编号（逗号 / 范围，最多%d个）", maxBookingTargets), &venueNamesInput, func(value string) error {
		names, err := parseVenueNames(value)
		if err != nil {
			return err
		}
		_, err = buildBookingTargets(venues, names, dateStr, 0)
		return err
	})
	if err != nil {
		return err
	}
	venueNames, _ := parseVenueNames(venueNamesInput)

	var hourInput string
	err = inputForm("预约开始小时", &hourInput, func(value string) error {
		hour, err := strconv.Atoi(value)
		if err != nil || hour < 0 || hour > 23 {
			return fmt.Errorf("请输入 0 到 23 的整数小时")
		}
		_, err = findHourIDByTime(hours, hour)
		return err
	})
	if err != nil {
		return err
	}
	hour, _ := strconv.Atoi(hourInput)
	hourID, _ := findHourIDByTime(hours, hour)

	targets, err := buildBookingTargets(venues, venueNames, dateStr, hourID)
	if err != nil {
		return fmt.Errorf("查找场地失败: %w", err)
	}

	targetInterval := checkInterval / time.Duration(len(targets))
	networkFailures := make([]int, len(targets))

	for {
		waitUntilNextAvailableTime()

		fmt.Printf("正在检查场馆ID %d 的订单信息...\n", venueID)

		for i, target := range targets {
			done, err := processBookingTarget(settings, venueID, target, dateStr, hour, headers, &networkFailures[i])
			if err != nil {
				return err
			}
			if done {
				return nil
			}

			time.Sleep(targetInterval)
		}
	}

}
