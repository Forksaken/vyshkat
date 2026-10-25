package main

import (
	"bytes"
	"math"
	"strings"
	"testing"
)

type stubInspector struct {
	accept bool
	seen   Vehicle
	calls  int
}

func (s *stubInspector) Inspect(v Vehicle) bool {
	s.seen = v
	s.calls++
	return s.accept
}

func TestInspectionDecision(t *testing.T) {
	for _, accept := range []bool{false, true} {
		name := "reject"
		if accept {
			name = "accept"
		}
		t.Run(name, func(t *testing.T) {
			stub := &stubInspector{accept: accept}
			park := NewContainer(stub).Park
			// специально противоположно решению заглушки
			vehicle := Scooter{VehicleData{Item{1, "Самокат"}, 6, !accept}, EnergyUsage{2}}
			err := park.Add(vehicle)
			if (err == nil) != accept || stub.calls != 1 || stub.seen != vehicle {
				t.Fatalf("приёмка: err=%v, вызовов=%d, объект=%v", err, stub.calls, stub.seen)
			}
			vehicles, things := park.Counts()
			if accept {
				if vehicles != 1 || things != 0 || park.TotalEnergy() != 2 || len(park.ForBeginners()) != 1 {
					t.Fatal("принятый транспорт не попал в отчёты")
				}
			} else {
				if vehicles != 0 || things != 0 || park.TotalEnergy() != 0 || len(park.ForBeginners()) != 0 {
					t.Fatal("отклонённый транспорт попал в отчёты")
				}
				if err := park.Add(Helmet{Item{1, "Шлем"}}); err != nil || stub.calls != 1 {
					t.Fatalf("отказ занял номер или вещь отправлена на техосмотр: %v", err)
				}
			}
		})
	}
}

func TestServiceCenter(t *testing.T) {
	for _, working := range []bool{false, true} {
		park := NewContainer(nil).Park
		err := park.Add(Bicycle{VehicleData{Item{1, "Велосипед"}, 5, working}})
		if (err == nil) != working {
			t.Errorf("исправность=%v, ошибка=%v", working, err)
		}
	}
}

func TestEmptyPark(t *testing.T) {
	park := NewContainer(nil).Park
	vehicles, things := park.Counts()
	if vehicles != 0 || things != 0 || park.TotalEnergy() != 0 || len(park.Inventory()) != 0 || len(park.ForBeginners()) != 0 {
		t.Fatal("у пустого парка должны быть пустые отчёты")
	}
}

func TestTotalsAndInventory(t *testing.T) {
	park := NewContainer(nil).Park
	items := []Inventory{
		Scooter{VehicleData{Item{1, "Самокат"}, 6, true}, EnergyUsage{1.5}},
		Bicycle{VehicleData{Item{2, "Велосипед"}, 5, true}},
		ElectricBicycle{VehicleData{Item{3, "Электровелосипед"}, 8, true}, EnergyUsage{2.5}},
		Helmet{Item{4, "Шлем"}},
		ChargingStation{Item{5, "Зарядная станция"}, EnergyUsage{3}},
	}
	for _, item := range items {
		if err := park.Add(item); err != nil {
			t.Fatal(err)
		}
	}
	vehicles, things := park.Counts()
	if vehicles != 3 || things != 2 || park.TotalEnergy() != 7 {
		t.Fatalf("транспорт=%d, вещи=%d, энергия=%v", vehicles, things, park.TotalEnergy())
	}
	inventory := park.Inventory()
	if len(inventory) != len(items) {
		t.Fatalf("на балансе %d объектов вместо %d", len(inventory), len(items))
	}
	for i, item := range inventory {
		if item.Name() != items[i].Name() || item.InventoryNumber() != items[i].InventoryNumber() {
			t.Errorf("потеряны название или номер объекта %d", i)
		}
	}
	inventory[0] = Helmet{Item{99, "Подмена"}}
	if park.Inventory()[0].InventoryNumber() != 1 {
		t.Fatal("изменение выданного списка испортило баланс")
	}
	if err := park.Add(Helmet{Item{1, "Другой шлем"}}); err == nil {
		t.Fatal("номер транспорта повторно занят вещью")
	}
	if len(park.Inventory()) != 5 || park.TotalEnergy() != 7 {
		t.Fatal("ошибка добавления изменила баланс")
	}
}

func TestForBeginners(t *testing.T) {
	park := NewContainer(nil).Park
	for _, ease := range []int{1, 5, 6, 10} {
		if err := park.Add(Bicycle{VehicleData{Item{ease, "Велосипед"}, ease, true}}); err != nil {
			t.Fatal(err)
		}
	}
	if err := park.Add(Helmet{Item{20, "Шлем"}}); err != nil {
		t.Fatal(err)
	}
	got := park.ForBeginners()
	if len(got) != 2 || got[0].Ease() != 6 || got[1].Ease() != 10 {
		t.Fatalf("ожидались только оценки 6 и 10, получено %v", got)
	}
}

func TestInvalidData(t *testing.T) {
	tests := []struct {
		name string
		item Inventory
	}{
		{"nil", nil},
		{"zero_number", Helmet{Item{0, "Шлем"}}},
		{"negative_number", Helmet{Item{-1, "Шлем"}}},
		{"empty_name", Helmet{Item{1, "  "}}},
		{"ease_zero", Bicycle{VehicleData{Item{1, "Велосипед"}, 0, true}}},
		{"ease_eleven", Bicycle{VehicleData{Item{1, "Велосипед"}, 11, true}}},
		{"negative_energy", ChargingStation{Item{1, "Станция"}, EnergyUsage{-1}}},
		{"nan_energy", ChargingStation{Item{1, "Станция"}, EnergyUsage{math.NaN()}}},
		{"infinite_energy", ChargingStation{Item{1, "Станция"}, EnergyUsage{math.Inf(1)}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &stubInspector{accept: true}
			park := NewContainer(stub).Park
			if err := park.Add(tt.item); err == nil {
				t.Fatal("ожидалась ошибка")
			}
			if len(park.Inventory()) != 0 || stub.calls != 0 {
				t.Fatal("неверные данные изменили баланс или попали на техосмотр")
			}
		})
	}
}

func TestMissingInspector(t *testing.T) {
	park := NewPark(nil)
	if err := park.Add(Bicycle{VehicleData{Item{1, "Велосипед"}, 6, true}}); err == nil {
		t.Fatal("без техосмотра транспорт нельзя принять")
	}
}

func TestDemo(t *testing.T) {
	var out bytes.Buffer
	runDemo(&out)
	for _, want := range []string{
		"Транспорт: 0, вещи: 0, всего: 0",
		"Самокат №1: принят на баланс",
		"Неисправный самокат №4: техосмотр: отказ",
		"Шлем №4: принят на баланс",
		"инвентарный номер 1 уже занят",
		"простота для новичка должна быть от 1 до 10",
		"Транспорт: 3, вещи: 2, всего: 5",
		"Суточный расход энергии: 7.0 кВтч",
		"Самокат №1, оценка 6",
		"Электровелосипед №3, оценка 8",
		"Зарядная станция №5",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("в демонстрации отсутствует %q", want)
		}
	}
}
