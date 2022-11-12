package usecase

import "github.com/go-lang-api/internal/order/entity"

type GetTotalOutputDTO struct {
	Total int
}

type GetTotalUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewGetTotalUseCase(orderRepository entity.OrderRepositoryInterface) *GetTotalUseCase {
	return &GetTotalUseCase{
		OrderRepository: orderRepository,
	}
}

func (g *GetTotalUseCase) Execute() (*GetTotalOutputDTO, error) {
	total, err := g.OrderRepository.GetTotal()
	if err != nil {
		return nil, err
	}
	return &GetTotalOutputDTO{
		Total: total,
	}, nil
}
