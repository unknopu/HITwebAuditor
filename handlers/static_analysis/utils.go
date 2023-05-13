package static_analysis

import (
	"auditor/core/utils"
	"auditor/entities"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/labstack/echo/v4"
)

func buildPageInfomation(e []*entities.SQLiReport) *entities.Page {
	if len(e) <= 0 {
		return nil
	}

	var low, medium, high, critical int

	for _, r := range e {
		switch r.Level {
		case entities.LOW:
			low += 1
		case entities.MEDIUM:
			medium += 1
		case entities.HIGH:
			high += 1
		case entities.CRITICAL:
			critical += 1
		}
	}

	pif := &entities.PageInformation{
		Vulnerabilities: len(e),
		Low:             low,
		Medium:          medium,
		High:            high,
		Critical:        critical,
	}

	return entities.NewPage(*pif, e)
}

func fetchSourceCode(c echo.Context, f *StaticAnalysisForm) (string, error) {
	if f.URL != nil {
		return utils.GetPageHTML(*f.URL, ""), nil
	}

	content, err := fileContent(c)
	if err != nil {
		return "", err
	}

	return content, nil
}

func fileContent(c echo.Context) (string, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", err
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// remove the file when it’s no longer needed.
	temp, err := ioutil.TempFile("temp", "prefix")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(temp.Name())

	// Copy
	if _, err = io.Copy(temp, src); err != nil {
		return "", err
	}

	content, err := os.ReadFile(temp.Name())
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func validatePHPSource(content string) string {
	re := regexp.MustCompile(`(?s)<\?php(.*?)\?>`)

	// var sourceCode string
	if re.MatchString(content) {
		// 	matches := re.FindAllStringSubmatch(string(content), -1)
		// 	for _, match := range matches {
		// 		sourceCode += match[1]
		// 	}
		return content
	}

	return ""
	// return strings.Trim(sourceCode, "\n")
}
