package config

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/funcBank_Api/models"
)

var (
	planRegex   = regexp.MustCompile(`(?i)direct|regular|retail|institutional`)
	optionRegex = regexp.MustCompile(`(?i)growth|idcw|dividend|bonus|reinvestment|payout|option`)
	freqRegex   = regexp.MustCompile(`(?i)monthly|quarterly|daily|weekly|annual`)
	dashRegex   = regexp.MustCompile(`\s*-\s*`)
)

// ---------- Parent Extraction ----------

func ExtractParentName(name string) string {
	n := dashRegex.ReplaceAllString(name, "-")
	parts := strings.Split(n, "-")

	base := []string{}
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if planRegex.MatchString(t) || optionRegex.MatchString(t) || freqRegex.MatchString(t) {
			break
		}
		base = append(base, t)
	}

	if len(base) == 0 {
		return name
	}

	return strings.TrimSpace(strings.Join(base, " - "))
}

func CreateParentKey(name string) string {
	key := strings.ToLower(name)
	key = strings.ReplaceAll(key, "&", "and")
	key = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(key, "-")
	key = regexp.MustCompile(`-+`).ReplaceAllString(key, "-")
	return strings.Trim(key, "-")
}

// ---------- Extract Plan/Option/Frequency ----------

func ExtractPlanType(name string) string {
	l := strings.ToLower(name)
	switch {
	case strings.Contains(l, "direct"):
		return "Direct"
	case strings.Contains(l, "regular"):
		return "Regular"
	case strings.Contains(l, "retail"):
		return "Retail"
	case strings.Contains(l, "institutional"):
		return "Institutional"
	default:
		return "Unknown"
	}
}

func ExtractOptionType(name string) string {
	l := strings.ToLower(name)
	switch {
	case strings.Contains(l, "growth"):
		return "Growth"
	case strings.Contains(l, "idcw"):
		return "IDCW"
	case strings.Contains(l, "bonus"):
		return "Bonus"
	case strings.Contains(l, "dividend"):
		return "Dividend"
	case strings.Contains(l, "reinvestment"):
		return "Reinvestment"
	default:
		return "Unknown"
	}
}

func ExtractFrequency(name string) string {
	l := strings.ToLower(name)
	switch {
	case strings.Contains(l, "monthly"):
		return "Monthly"
	case strings.Contains(l, "quarterly"):
		return "Quarterly"
	case strings.Contains(l, "daily"):
		return "Daily"
	case strings.Contains(l, "weekly"):
		return "Weekly"
	case strings.Contains(l, "annual"):
		return "Annual"
	default:
		return "None"
	}
}

// ---------- AMC Logo Key ----------

func NormalizeFundHouseKey(s string) string {
	key := strings.ToLower(s)
	key = strings.ReplaceAll(key, "&", "and")
	key = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(key, "-")
	return strings.Trim(key, "-")
}

// ---------- Active Check ----------

func IsActive(date time.Time) bool {
	return time.Since(date).Hours() < 24*60 // 60 days
}
func ParseCategory(category string) (string, string) {
    // Extract text inside parentheses
    start := strings.Index(category, "(")
    end := strings.LastIndex(category, ")")
    if start == -1 || end == -1 || end <= start {
        return "", ""
    }

    inside := strings.TrimSpace(category[start+1 : end])

    // Split into category and subcategory
    parts := strings.SplitN(inside, "-", 2)
    if len(parts) != 2 {
        return strings.TrimSpace(inside), ""
    }

    cat := strings.TrimSpace(parts[0])
    sub := strings.TrimSpace(parts[1])
    return cat, sub
}


// ---------- MAIN PARSER ----------

func ParseNAVAll(path string) ([]models.SchemeDetail, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var records []models.SchemeDetail
	sc := bufio.NewScanner(f)

	var category string
	var fundHouse string

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}

		// Category header
		if strings.HasPrefix(line, "Open Ended Schemes(") ||
			strings.HasPrefix(line, "Close Ended Schemes(") ||
			strings.HasPrefix(line, "Solution Oriented Schemes(") ||
			strings.HasPrefix(line, "Interval Fund(") {
			category = line
			continue
		}

		// Fund House
		if !strings.Contains(line, ";") && !strings.ContainsAny(line, "0123456789") {
			fundHouse = line
			continue
		}

		// Scheme row
		parts := strings.Split(line, ";")
		if len(parts) < 6 {
			continue
		}

		rawName := strings.TrimSpace(parts[3])
		parentName := ExtractParentName(rawName)
		parentKey := CreateParentKey(parentName)
		fundHouseKey := NormalizeFundHouseKey(fundHouse)

		planType := ExtractPlanType(rawName)
		optionType := ExtractOptionType(rawName)
		frequency := ExtractFrequency(rawName)

		navVal := strings.TrimSpace(parts[4])
		dateVal := strings.TrimSpace(parts[5])

		navDate, _ := time.Parse("02-Jan-2006", dateVal) // AMFI format
		if IsActive(navDate) == false {
			continue // skip inactive
		}
		mainCat, subCat := ParseCategory(category)
		record := models.SchemeDetail{
			SchemeCode:   strings.TrimSpace(parts[0]),
			ISINGrowth:   strings.TrimSpace(parts[1]),
			ISINReinvest: strings.TrimSpace(parts[2]),

			SchemeName: rawName,
			ParentName: parentName,
			ParentKey:  parentKey,

			FundHouse:    fundHouse,
			FundHouseKey: fundHouseKey,

			CategoryHeader: category,
			Category: mainCat,
			SubCategory: subCat,
			PlanType: planType,
			OptionType: optionType,
			Frequency: frequency,

			NAV:      parseFloat(navVal),

			DisplayName: parentName + " (" + planType + " - " + optionType + ")",
		}

		records = append(records, record)
	}

	return records, sc.Err()
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
