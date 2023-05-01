package report

import (
	"auditor/core/context"
	"auditor/core/utils"
	"auditor/entities"
	cf "auditor/handlers/cryptograhpical_failure"
	mc "auditor/handlers/miss_configuration"
	odc "auditor/handlers/outdated_component"
	"strings"

	"io/ioutil"
	"log"
	"net/http"

	"github.com/jinzhu/copier"
)

func (s Service) doMissConfig(c *context.Context, option *Form) *entities.Page {
	f := &mc.MCForm{}
	_ = copier.Copy(f, option)

	report, err := s.mcs.Init(c, f)
	if err != nil {
		return nil
	}

	return report.(*entities.Page)
}

func (s Service) doCryptoFailure(c *context.Context, option *Form) *entities.Page {
	f := &cf.CFForm{}
	_ = copier.Copy(f, option)

	report, err := s.cfs.Init(c, f)
	if err != nil {
		return nil
	}

	return report.(*entities.Page)
}

func (s Service) doOutdatedCpn(c *context.Context, ref []*entities.MissConfigurationReport) *entities.Page {
	if len(ref) <= 0 {
		return nil
	}

	f := &odc.OutdatedComponentForm{}
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

	return report.(*entities.Page)
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
