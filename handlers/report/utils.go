package report

import (
	"auditor/core/context"
	"auditor/core/utils"
	"auditor/entities"
	cf "auditor/handlers/cryptograhpical_failure"
	lfi "auditor/handlers/lfi"
	mc "auditor/handlers/miss_configuration"
	odc "auditor/handlers/outdated_component"
	sqli "auditor/handlers/sqli"
	xss "auditor/handlers/xss"

	"strings"

	"io/ioutil"
	"log"
	"net/http"

	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s Service) doMissConfig(c *context.Context, option *Form) []*entities.MissConfigurationReport {
	f := &mc.MCForm{}
	_ = copier.Copy(f, option)

	report, err := s.mcs.Init(c, f)
	if err != nil {
		return nil
	}

	return report
}

func (s Service) doCryptoFailure(c *context.Context, option *Form) []*entities.CryptoFailureReport {
	f := &cf.CFForm{}
	_ = copier.Copy(f, option)

	report, err := s.cfs.Init(c, f)
	if err != nil {
		return nil
	}

	return report
}

func (s Service) doOutdatedCpn(c *context.Context, ref []*entities.MissConfigurationReport, id primitive.ObjectID) []*entities.OutdatedComponentsReport {
	if len(ref) <= 0 {
		return nil
	}

	f := &odc.OutdatedComponentForm{}
	f.ReportNumber = id
	for _, r := range ref {
		if len(r.Payload) == 0 {
			continue
		}
		if strings.ContainsAny(r.Payload[0], "phpPHP") {
			f.Refer = append(f.Refer, r.Payload[0])
		}
	}

	report, err := s.odcs.Init(c, f)
	if err != nil {
		return nil
	}

	return report
}

func (s Service) doXSS(c *context.Context, option *Form) []*entities.XSSReport {
	f := &xss.XSSForm{}
	_ = copier.Copy(f, option)

	report, err := s.xsss.Init(c, f)
	if err != nil {
		return nil
	}

	return report
}

func (s Service) doSQLI(c *context.Context, option *Form) interface{} {
	f := &sqli.SqliForm{}
	_ = copier.Copy(f, option)

	reports, err := s.sqlis.Init(c, f)
	if err != nil {
		return nil
	}

	return reports
}

func (s Service) doLFI(c *context.Context, option *Form) []*entities.LFIReport {
	f := &lfi.LFIForm{}
	_ = copier.Copy(f, option)

	report, err := s.lfis.Init(c, f)
	if err != nil {
		return nil
	}

	return report
}

func injectPayload(option entities.LFI, payload string) string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = utils.C

	q := option.URL.Query()
	q.Set(option.Parameter, payload)
	option.URL.RawQuery = q.Encode()

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodGet, option.URL.String(), nil)

	res, err := client.Do(r)
	if err != nil {
		log.Println("[*] GET HTML: ", err)
	}
	if res != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}
