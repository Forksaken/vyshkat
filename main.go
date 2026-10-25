package main

import (
	"fmt"
	"io"
	"os"
)

// Container собирает зависимости один раз и передаёт техосмотр парку.
type Container struct {
	Park *Park
}

func NewContainer(inspector Inspector) *Container {
	if inspector == nil {
		inspector = ServiceCenter{}
	}
	return &Container{Park: NewPark(inspector)}
}

func main() {
	runDemo(os.Stdout)
}

func runDemo(out io.Writer) {
	park := NewContainer(nil).Park
	fmt.Fprintln(out, "ВышКат — Доедем до дедлайна")
	printReport(out, park)

	items := []Inventory{
		Scooter{VehicleData{Item{1, "Самокат"}, 6, true}, EnergyUsage{1.5}},
		Bicycle{VehicleData{Item{2, "Велосипед"}, 5, true}},
		ElectricBicycle{VehicleData{Item{3, "Электровелосипед"}, 8, true}, EnergyUsage{2.5}},
		Scooter{VehicleData{Item{4, "Неисправный самокат"}, 9, false}, EnergyUsage{10}},
		Helmet{Item{4, "Шлем"}},
		ChargingStation{Item{5, "Зарядная станция"}, EnergyUsage{3}},
		Helmet{Item{1, "Шлем с занятым номером"}},
		Bicycle{VehicleData{Item{6, "Велосипед с неверной оценкой"}, 11, true}},
	}
	fmt.Fprintln(out, "\nПриёмка:")
	for _, item := range items {
		if err := park.Add(item); err != nil {
			fmt.Fprintf(out, "%s №%d: %v\n", item.Name(), item.InventoryNumber(), err)
		} else {
			fmt.Fprintf(out, "%s №%d: принят на баланс\n", item.Name(), item.InventoryNumber())
		}
	}
	fmt.Fprintln(out, "\nИтоговый отчёт:")
	printReport(out, park)
}

func printReport(out io.Writer, park *Park) {
	vehicles, things := park.Counts()
	fmt.Fprintf(out, "Транспорт: %d, вещи: %d, всего: %d\n", vehicles, things, vehicles+things)
	fmt.Fprintf(out, "Суточный расход энергии: %.1f кВтч\n", park.TotalEnergy())
	fmt.Fprintln(out, "Для новичков (оценка ≥ 6):")
	beginnerVehicles := park.ForBeginners()
	if len(beginnerVehicles) == 0 {
		fmt.Fprintln(out, "  нет")
	}
	for _, vehicle := range beginnerVehicles {
		fmt.Fprintf(out, "  %s №%d, оценка %d\n", vehicle.Name(), vehicle.InventoryNumber(), vehicle.Ease())
	}
	fmt.Fprintln(out, "Всё на балансе:")
	items := park.Inventory()
	if len(items) == 0 {
		fmt.Fprintln(out, "  нет")
	}
	for _, item := range items {
		fmt.Fprintf(out, "  %s №%d\n", item.Name(), item.InventoryNumber())
	}
}
