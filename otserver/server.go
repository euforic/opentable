package otserver

import (
	"context"

	"github.com/euforic/opentable/opentable"
	"github.com/euforic/opentable/otpb"
)

// OTServer is the struct that implments the OTServiceServer interface
type OTServer struct{}

// New creates a new OTServer to be registerd with gRPC
func New() *OTServer {
	return &OTServer{}
}

// Search performs a search against opentable.com with the given search opts
// and returns the scraped resturants with the available reservations
func (s *OTServer) Search(ctx context.Context, req *otpb.SearchReq) (*otpb.SearchRes, error) {
	result, err := opentable.Search(opentable.SearchOpts{
		UserAgent: req.UserAgent,
		People:    req.People,
		Time:      *req.Time,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Term:      req.Term,
		Sort:      req.Sort.String(),
		Opts:      req.Opts,
	})
	if err != nil {
		return nil, err
	}

	resturants := []*otpb.Resturant{}

	for _, v := range result {

		rsv := []*otpb.Reservation{}

		for _, rv := range v.Reservations {
			rsv = append(rsv, &otpb.Reservation{
				Time: &rv.Time,
				Url:  rv.URL,
			})
		}

		resturants = append(resturants, &otpb.Resturant{
			ID:           v.ID,
			Name:         v.Name,
			URL:          v.URL,
			Recommended:  v.Recommended,
			Reservations: rsv,
		})
	}

	return &otpb.SearchRes{Resturants: resturants}, nil
}
