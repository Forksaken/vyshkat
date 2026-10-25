package main

import (
	"fmt"
	"math"
	"strings"
)

// Inspect возвращает true, если транспорт можно принять в парк
type Inspector interface {
	Inspect(Vehicle) bool
}

type ServiceCenter struct{}

func (ServiceCenter) Inspect(v Vehicle) bool { return v.IsWorking() }

type Park struct {
	inspector Inspector
	items     []Inventory
}

func NewPark(inspector Inspector) *Park {
	return &Park{inspector: inspector}
}

// Add ставит объект на баланс
//  транспорт сначала проходит техосмотр.
func (p *Park) Add(item Inventory) error {
	if item == nil || item.InventoryNumber() <= 0 || strings.TrimSpace(item.Name()) == "" {
		return fmt.Errorf("нужны положительный инвентарный номер и название")
	}
	for _, existing := range p.items {
		if existing.InventoryNumber() == item.InventoryNumber() {
			return fmt.Errorf("инвентарный номер %d уже занят", item.InventoryNumber())
		}
	}
	if consumer, ok := item.(EnergyConsumer); ok {
		energy := consumer.DailyEnergy()
		if energy < 0 || math.IsNaN(energy) || math.IsInf(energy, 0) {
			return fmt.Errorf("расход энергии должен быть конечным неотрицательным числом")
		}
	}
	if vehicle, ok := item.(Vehicle); ok {
		if vehicle.Ease() < 1 || vehicle.Ease() > 10 {
			return fmt.Errorf("простота для новичка должна быть от 1 до 10")
		}
		if p.inspector == nil {
			return fmt.Errorf("техосмотр не настроен")
		}
		if !p.inspector.Inspect(vehicle) {
			return fmt.Errorf("техосмотр: отказ в приёмке")
		}
	}
	p.items = append(p.items, item)
	return nil
}

func (p *Park) Counts() (vehicles, things int) {
	for _, item := range p.items {
		if _, ok := item.(Vehicle); ok {
			vehicles++
		} else {
			things++
		}
	}
	return
}

func (p *Park) TotalEnergy() float64 {
	var total float64
	for _, item := range p.items {
		if consumer, ok := item.(EnergyConsumer); ok {
			total += consumer.DailyEnergy()
		}
	}
	return total
}

func (p *Park) ForBeginners() []Vehicle {
	var result []Vehicle
	for _, item := range p.items {
		if vehicle, ok := item.(Vehicle); ok && vehicle.Ease() >= 6 {
			result = append(result, vehicle)
		}
	}
	return result
}

func (p *Park) Inventory() []Inventory {
	return append([]Inventory(nil), p.items...)
}
