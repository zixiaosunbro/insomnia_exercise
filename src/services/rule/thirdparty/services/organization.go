package services

import "context"

type IOrganizationService interface {
	GetUserProjects(ctx context.Context, userId uint64, organId int32) ([]*ProjectInfo, error)
}

type OrganizationService struct {
}

type ProjectInfo struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

func (o *OrganizationService) GetUserProjects(ctx context.Context, userId uint64, organId int32) ([]*ProjectInfo, error) {
	if organId == 101 && userId == 1010 {
		return []*ProjectInfo{
			{
				Id:   14,
				Name: "insomnia",
			},
		}, nil
	}
	return nil, nil
}
