package report

import (
	"auditor/core/context"
	"auditor/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s Service) fetchSQLI(c *context.Context, id primitive.ObjectID) *entities.Page {
	reports, err := s.sqlis.FetchReport(c, id)
	if err != nil {
		return nil
	}

	return reports
}

func (s Service) fetchSecMissCon(c *context.Context, id primitive.ObjectID) *entities.Page {
	reports, err := s.mcs.FetchReport(c, id)
	if err != nil {
		return nil
	}

	return reports
}

func (s Service) fetchOutdatedCpn(c *context.Context, id primitive.ObjectID) *entities.Page {
	reports, err := s.odcs.FetchReport(c, id)
	if err != nil {
		return nil
	}

	return reports
}

func (s Service) fetchXSS(c *context.Context, id primitive.ObjectID) *entities.Page {
	reports, err := s.xsss.FetchReport(c, id)
	if err != nil {
		return nil
	}

	return reports
}

func (s Service) fetchCryptoFailure(c *context.Context, id primitive.ObjectID) *entities.Page {
	reports, err := s.cfs.FetchReport(c, id)
	if err != nil {
		return nil
	}

	return reports
}

func (s Service) fetchLFI(c *context.Context, id primitive.ObjectID) *entities.Page {
	reports, err := s.lfis.FetchReport(c, id)
	if err != nil {
		return nil
	}

	return reports
}
