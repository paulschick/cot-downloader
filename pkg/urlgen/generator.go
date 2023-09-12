package urlgen

import "fmt"

func GetLegacyArchiveURLForYear(year int) string {
	if year <= 2003 {
		return fmt.Sprintf("https://www.cftc.gov/files/dea/history/deafut_xls_%d.zip", year)
	} else {
		return fmt.Sprintf("https://www.cftc.gov/files/dea/history/dea_fut_xls_%d.zip", year)
	}
}

func GenerateURLs(reportType string, start, end int) []string {
	var urls []string

	for i := start; i <= end; i++ {
		// TODO - add report types
		switch reportType {
		case "legacy":
			url := GetLegacyArchiveURLForYear(i)
			urls = append(urls, url)
		default:
			url := GetLegacyArchiveURLForYear(i)
			urls = append(urls, url)
		}
	}

	return urls
}
