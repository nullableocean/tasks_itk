package oop

import (
	"errors"
	"fmt"
	"testing"
)

func TestEngineWork(t *testing.T) {
	car := NewCar("Oka")

	gotStatus := car.GetEngineStatus()
	expectedStatus := false

	if gotStatus != expectedStatus {
		t.Fatalf("got engine status: %v, expected: %v", gotStatus, expectedStatus)
	}

	err := car.StartEngine()
	if err != nil {
		t.Fatalf("start engine. got error: %v, expected: nil", err)
	}

	err = car.StartEngine()
	expectedErr := ErrEngineAlreadyRunning

	if !errors.Is(err, expectedErr) {
		t.Fatalf("repeat start engine. got error: %v, expected: %v", err, expectedErr)
	}

	err = car.StopEngine()
	if err != nil {
		t.Fatalf("stop engine. got error: %v, expected: nil", err)
	}

	err = car.StopEngine()
	expectedErr = ErrEngineOff

	if !errors.Is(err, expectedErr) {
		t.Fatalf("repeat stop engine. got error: %v, expected: %v", err, expectedErr)
	}
}

func TestVehiclePolymorph(t *testing.T) {
	var veh interface{}

	car := NewCar("Kia")
	veh = car
	_, ok := veh.(Vehicle)
	if !ok {
		t.Fatalf("car not implement vehicle interface")
	}

	elCar := NewElectricCar("Electro-kia")
	veh = elCar
	_, ok = veh.(Vehicle)
	if !ok {
		t.Fatalf("electro car not implement vehicle interface")
	}

	truck := NewTruck("Truck", 34.1)
	veh = truck
	_, ok = veh.(Vehicle)
	if !ok {
		t.Fatalf("truck not implement vehicle interface")
	}
}

func TestVehicleMethodsOverload(t *testing.T) {
	car := NewCar("Honda")
	truck := NewTruck("Men", 32.2)

	honkCar := car.Honk()
	honkTruck := truck.Honk()

	fmt.Println(truck.GetInfo())

	if honkCar == honkTruck {
		t.Fatalf("method honk not overloaded. car honk == truck honk")
	}
}
