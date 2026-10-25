package main

type Inventory interface {
	InventoryNumber() int
	Name() string
}

// DailyEnergy возвращает суточный расход в кВтч
type EnergyConsumer interface {
	DailyEnergy() float64
}

type Vehicle interface {
	Inventory
	Ease() int
	IsWorking() bool
}

type Item struct {
	Number int
	Title  string
}

func (i Item) InventoryNumber() int { return i.Number }
func (i Item) Name() string         { return i.Title }

type VehicleData struct {
	Item
	Simplicity int
	Working    bool
}

func (v VehicleData) Ease() int       { return v.Simplicity }
func (v VehicleData) IsWorking() bool { return v.Working }

type EnergyUsage struct {
	KWh float64
}

func (e EnergyUsage) DailyEnergy() float64 { return e.KWh }

type Scooter struct {
	VehicleData
	EnergyUsage
}

type Bicycle struct {
	VehicleData
}

type ElectricBicycle struct {
	VehicleData
	EnergyUsage
}

type Helmet struct {
	Item
}

type ChargingStation struct {
	Item
	EnergyUsage
}
