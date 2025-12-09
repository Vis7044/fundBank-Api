package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/funcBank_Api/models"
	"github.com/xuri/excelize/v2"
)

func FilterNavByInterval(navData []models.FundNav, interval string, swpDate int, startDate string) ([]models.FundNav, error) {

	if len(navData) == 0 {
		return navData, nil
	}

	layout := "02-01-2006"
	layoutISO := "2006-01-02"

	// 1️⃣ Include Investment Date
	investmentDate, err := time.Parse(layout, navData[0].Date)
	if err != nil {
		return nil, err
	}

	filteredNav := []models.FundNav{navData[0]}
	lastIncludedDate := investmentDate

	// 2️⃣ Parse Input Start Date
	start, err := time.Parse(layoutISO, startDate)
	if err != nil {
		return nil, err
	}

	// 3️⃣ Construct next SWP date based on calendar day
	var swpDay int = swpDate
	if swpDay < 1 || swpDay > 28 {
		return nil, fmt.Errorf("SWP date must be between 1 and 28")
	}

	// Create SWP date in same month/year as startDate
	firstSWPDate := time.Date(start.Year(), start.Month(), swpDay, 0, 0, 0, 0, start.Location())

	// If that SWP date is BEFORE start date → move to next month
	if firstSWPDate.Before(start) {
		firstSWPDate = firstSWPDate.AddDate(0, 1, 0)
	}

	// 4️⃣ Find first NAV on or after calculated SWP date
	var firstSWPNav *models.FundNav
	for _, nav := range navData {
		navDate, _ := time.Parse(layout, nav.Date)
		if !navDate.Before(firstSWPDate) {
			firstSWPNav = &nav
			lastIncludedDate = navDate
			break
		}
	}

	// If found, add it
	if firstSWPNav != nil {
		filteredNav = append(filteredNav, *firstSWPNav)
	} else {
		return filteredNav, nil
	}

	// 5️⃣ Continue filtering based on interval
	for _, navPoint := range navData {
		currentDate, err := time.Parse(layout, navPoint.Date)
		if err != nil {
			return nil, err
		}

		if currentDate.Before(lastIncludedDate) {
			continue
		}

		diff := currentDate.Sub(lastIncludedDate).Hours() / 24

		switch interval {
		case "weekly":
			if diff >= 7 {
				filteredNav = append(filteredNav, navPoint)
				lastIncludedDate = currentDate
			}
		case "fortnightly":
			if diff >= 14 {
				filteredNav = append(filteredNav, navPoint)
				lastIncludedDate = currentDate
			}
		case "monthly":
			if diff >= 30 {
				filteredNav = append(filteredNav, navPoint)
				lastIncludedDate = currentDate
			}
		case "quarterly":
			if diff >= 90 {
				filteredNav = append(filteredNav, navPoint)
				lastIncludedDate = currentDate
			}
		default:
			return navData, nil
		}
	}

	return filteredNav, nil
}

/*
	type SWPResponse struct {
		TotalInvestedAmount float64 `json:"total_invested_amount"`
		Units               float64 `json:"units"`
		CumulativeUnits     float64 `json:"cumulative_units"`
		RemainingUnits      float64 `json:"remaining_units"`
		CashFlow            float64 `json:"cash_flow"`
		NetAmount           float64 `json:"net_amount"`
		CapitalGainsLoss    float64 `json:"capital_gains_loss"`
		CurrentNAV          float64 `json:"current_nav"`
		CurrentValue        float64 `json:"current_value"`
	}
*/
func CalculateSWPResponse(navData []models.FundNav, total_invested_amount float64, withdrawal_amount float64) ([]models.SWPResponse, error) {
	var swpResponses []models.SWPResponse

	var swpResponse models.SWPResponse
	swpResponse.NetAmount = total_invested_amount
	swpResponse.CashFlow = total_invested_amount
	swpResponse.CapitalGainsLoss = 0
	navValue, err := strconv.ParseFloat(navData[0].Nav, 32)
	if err != nil {
		return nil, err
	}
	swpResponse.Units = total_invested_amount / navValue
	swpResponse.CumulativeUnits = swpResponse.Units
	swpResponse.CurrentValue = swpResponse.CumulativeUnits * navValue
	swpResponse.CurrentNAV = navValue
	swpResponse.CurrentDate = navData[0].Date
	swpResponses = append(swpResponses, swpResponse)

	for _, navPoint := range navData[1:] {
		var swpResponse models.SWPResponse

		navValue, err := strconv.ParseFloat(navPoint.Nav, 32)
		if err != nil {
			return nil, err
		}
		swpResponse.CurrentNAV = navValue
		swpResponse.CurrentDate = navPoint.Date

		// Calculate units redeemed
		unitsRedeemed := withdrawal_amount / navValue
		flag := false
		if unitsRedeemed > swpResponses[len(swpResponses)-1].CumulativeUnits {
			unitsRedeemed = swpResponses[len(swpResponses)-1].CumulativeUnits // can't redeem more than available
			flag = true
		}
		swpResponse.Units = -unitsRedeemed
		swpResponse.CumulativeUnits = swpResponses[len(swpResponses)-1].CumulativeUnits - unitsRedeemed
		swpResponse.CashFlow = -withdrawal_amount
		swpResponse.NetAmount = swpResponses[len(swpResponses)-1].NetAmount - withdrawal_amount

		currentValue := swpResponse.CumulativeUnits * navValue
		capitalGainsLoss := currentValue - swpResponse.NetAmount
		swpResponse.CurrentValue = currentValue
		swpResponse.CapitalGainsLoss = capitalGainsLoss

		swpResponses = append(swpResponses, swpResponse)
		if flag {
			break
		}
	}

	return swpResponses, nil
}

func ExportSWPResponseToExcel(swpResponses []models.SWPResponse, schemeCode string) (string, error) {
	file := excelize.NewFile()
	sheetName := "SWP Plan"
	index, err := file.NewSheet(sheetName)
	if err != nil {
		return "", err
	}

	// Set headers
	headers := []string{"Date", "Units", "Cumulative Units", "Cash Flow", "Net Amount", "Capital Gains/Loss", "Current NAV", "Current Value"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		file.SetCellValue(sheetName, cell, header)
	}

	// Fill data
	for rowIndex, swp := range swpResponses {
		values := []interface{}{swp.CurrentDate, swp.Units, swp.CumulativeUnits, swp.CashFlow, swp.NetAmount, swp.CapitalGainsLoss, swp.CurrentNAV, swp.CurrentValue}
		for colIndex, value := range values {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			file.SetCellValue(sheetName, cell, value)
		}
	}

	file.SetActiveSheet(index)
	filePath := fmt.Sprintf("%s_swp_plan.xlsx", schemeCode)
	if err := file.SaveAs(filePath); err != nil {
		return "", err
	}
	return filePath, nil
}

func ReverseFundNavSlice(navData []models.FundNav) []models.FundNav {
	for i, j := 0, len(navData)-1; i < j; i, j = i+1, j-1 {
		navData[i], navData[j] = navData[j], navData[i]
	}
	return navData
}
