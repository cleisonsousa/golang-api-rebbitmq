package entity

import "errors"

type Order struct {
	ID         string
	Price      float64
	Tax        float64
	FinalPrice float64
}

func NewOrder(id string, price float64, tax float64) (*Order, error) {
	order := &Order{
		ID:    id,
		Price: price,
		Tax:   tax,
	}
	err := order.isValid()
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *Order) isValid() error {
	if o.ID == "" {
		return errors.New("invalida id")
	}
	if o.Price <= 0 {
		return errors.New("invalida id")
	}
	if o.Tax <= 0 {
		return errors.New("invalida tax")
	}
	return nil
}

func (o *Order) CalculateFinalPrice() error {
	o.FinalPrice = o.Price + o.Tax
	err := o.isValid()
	if err != nil {
		return err
	}
	return nil
}
