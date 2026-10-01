package world

type Location struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Region        string `json:"region"`
	Description   string `json:"description"`
	TravelSeconds int    `json:"travelSeconds"`
}

type Route struct {
	From          string `json:"from"`
	To            string `json:"to"`
	TravelSeconds int    `json:"travelSeconds"`
}

var locations = []Location{
	{ID: "central", Name: "Central District", Region: "Aster", Description: "The dense commercial heart of the world.", TravelSeconds: 0},
	{ID: "harbor", Name: "Harbor Ward", Region: "Aster", Description: "A working waterfront built around trade and transport.", TravelSeconds: 45},
	{ID: "oldtown", Name: "Old Town", Region: "Aster", Description: "A historic district where older institutions still shape daily life.", TravelSeconds: 35},
	{ID: "industrial", Name: "Industrial Belt", Region: "Aster", Description: "Factories, warehouses and logistics infrastructure.", TravelSeconds: 70},
	{ID: "highlands", Name: "North Highlands", Region: "Aster", Description: "A quieter elevated region beyond the central districts.", TravelSeconds: 120},
}

var routes = []Route{
	{From: "central", To: "harbor", TravelSeconds: 45}, {From: "harbor", To: "central", TravelSeconds: 45},
	{From: "central", To: "oldtown", TravelSeconds: 35}, {From: "oldtown", To: "central", TravelSeconds: 35},
	{From: "central", To: "industrial", TravelSeconds: 70}, {From: "industrial", To: "central", TravelSeconds: 70},
	{From: "oldtown", To: "highlands", TravelSeconds: 120}, {From: "highlands", To: "oldtown", TravelSeconds: 120},
}

func Locations() []Location { return append([]Location(nil), locations...) }
func Routes() []Route       { return append([]Route(nil), routes...) }
