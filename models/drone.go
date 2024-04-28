package models

import (
	"fmt"
)

type DroneMovement interface {
	MoveForward() error
	Accent() error
	Decend() error
}

type Drone struct {
	CurrentHeight   uint8
	Estate          *Estate
	MappedTrees     [][]uint8
	Travelled       uint16
	MaximumBattery  *uint16
	BatteryDrains   bool
	LastCoordinateX uint16
	LastCoordinateY uint16
}

func NewDrone(estate *Estate, estateTrees *[]Tree, maxDistance *uint16) *Drone {
	mappedTrees := make([][]uint8, estate.Length)
	for i := range mappedTrees {
		mappedTrees[i] = make([]uint8, estate.Width)
	}

	for _, tree := range *estateTrees {
		mappedTrees[tree.X-1][tree.Y-1] = tree.Height
	}

	drone := Drone{
		Estate:         estate,
		MappedTrees:    mappedTrees,
		MaximumBattery: maxDistance,
	}

	return &drone
}

func (d *Drone) TestFlight() {
	fmt.Print("Drone Start Flight At Plot 1,1 \n")

	xValue := uint16(1)
	for y := uint16(1); y <= d.Estate.Width; y++ {
		if xValue == uint16(1) {
			for x := xValue; x <= uint16(d.Estate.Length); x++ {
				fmt.Println(x, y)
			}
			xValue = d.Estate.Length
		} else if xValue == d.Estate.Length {
			for x := xValue; x > 0; x-- {
				fmt.Println(x, y)
			}
			xValue = 1
		}
	}
}

func (d *Drone) StartFlight() {
	fmt.Print("Drone Start Flight At Plot 1,1 \n\n")
	if d.MaximumBattery != nil {
		fmt.Printf("Drone Maximum Flight is %v m \n\n", *d.MaximumBattery)
	}

	xValue := uint16(1)
	for y := uint16(1); y <= d.Estate.Width; y++ {
		// y += 1
		if d.BatteryDrains {
			break
		}

		if xValue == uint16(1) {
			for x := xValue; x <= uint16(d.Estate.Length); x++ {
				if d.BatteryDrains {
					break
				}
				d.ReadData(x-1, y-1)
			}
			xValue = d.Estate.Length
		} else if xValue == d.Estate.Length {
			for x := xValue; x > 0; x-- {
				if d.BatteryDrains {
					break
				}
				d.ReadData(x-1, y-1)
			}
			xValue = 1
		}
	}
}

func (d *Drone) ReadData(x uint16, y uint16) {
	d.LastCoordinateX = x + 1
	d.LastCoordinateY = y + 1

	fmt.Printf("Drone At %v, %v\n", d.LastCoordinateX, d.LastCoordinateY)

	fmt.Println(x, y)
	// For the first plot, the drone must fly from ground level to a position exactly 1 meter above the plot or tree to get the data
	treeHigh := d.MappedTrees[x][y]
	var nextTreeHigh uint8
	if x < d.Estate.Length-1 {
		nextTreeHigh = d.MappedTrees[x+1][y]
	}

	// First Plot
	if x == 0 && y == 0 {
		if treeHigh > 0 {
			d.Accend(uint16(treeHigh+1), x, y)
		} else {
			d.Accend(uint16(1), x, y)
			if d.CurrentHeight < nextTreeHigh {
				d.Accend(uint16(nextTreeHigh), x, y)
			}
		}
	} else if x == d.Estate.Length-1 && y == d.Estate.Width-1 { // Last Plot
		d.Forward(x, y)
		d.Decend(uint16(d.CurrentHeight), x, y)

		fmt.Printf("Drone Landed at %v,%v . . .\n\n", x+1, y+1)
		return
	} else {
		if treeHigh == 0 {
			// Its a ground
			d.Forward(x, y)
		} else if d.CurrentHeight > treeHigh+1 {
			// If Current Position is HIGHER than next plot tree / ground
			// then move forward, and decend to 1 m above the next plot tree or ground

			// Calculate distance require to decend
			distance := d.CurrentHeight - treeHigh

			d.Forward(x, y)
			d.Decend(uint16(distance-1), x, y)
		} else if d.CurrentHeight < treeHigh+1 {
			// If Current Position is LOWER than next plot tree / ground
			// then move forward, and decend to 1 m above the next plot tree or ground

			// Calculate distance require to decend
			distance := treeHigh - d.CurrentHeight

			d.Accend(uint16(distance+1), x, y)
			d.Forward(x, y)
		} else {
			// If Current Position is SAME LEVEL with next plot tree / ground
			// just need to move forward

			d.Forward(x, y)
		}

		if d.BatteryDrains {
			fmt.Printf("Drone Rest at %v,%v . . .\n\n", x+1, y+1)
			return
		}

		fmt.Print("Drone Reading The Data . . .\n\n")
	}
}

func (d *Drone) Forward(x uint16, y uint16) {
	nextDistance := d.CheckBatteryBeforeDrains(10)

	fmt.Printf("Drone Move for %v m\n", nextDistance)

	if nextDistance <= 5 {
		d.LastCoordinateX = d.LastCoordinateX - 1
	}

	if d.BatteryDrains {
		return
	}
}

func (d *Drone) Accend(distance uint16, x uint16, y uint16) {
	nextDistance := d.CheckBatteryBeforeDrains(distance)

	d.CurrentHeight = d.CurrentHeight + uint8(nextDistance)
	fmt.Printf("Drone Accend for %v m\n", nextDistance)

	if d.BatteryDrains {
		return
	}
}

func (d *Drone) Decend(distance uint16, x uint16, y uint16) {
	nextDistance := d.CheckBatteryBeforeDrains(distance)

	d.CurrentHeight = d.CurrentHeight - uint8(nextDistance)
	fmt.Printf("Drone Decend for %v m\n", nextDistance)

	if d.BatteryDrains {
		return
	}
}

func (d *Drone) CheckBatteryBeforeDrains(nextDistance uint16) uint16 {
	// No Maximum Distance
	if d.MaximumBattery == nil {
		d.Travelled += uint16(nextDistance)
		return uint16(nextDistance)
	}

	// If battery is not enough for next distance will be travelled
	if (d.Travelled + uint16(nextDistance)) > *d.MaximumBattery {
		nextDistance = uint16(*d.MaximumBattery - d.Travelled)

		d.BatteryDrains = true
	} else if (d.Travelled + uint16(nextDistance)) == *d.MaximumBattery {
		d.BatteryDrains = true
	}

	d.Travelled += uint16(nextDistance)

	return uint16(nextDistance)
}
