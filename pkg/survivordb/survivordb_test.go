package survivordb

import (
	"os"
	"testing"
	"time"
)

// TestSurvivorDB_Open checks if opening the database works
func TestSurvivorDB_Open(t *testing.T) {
	got := Open("./test.db")
	if got == nil {
		t.Errorf("Open(\"./test.db\"): want: %v, got: %v", true, got != nil)
	}
}

// TestSurvivorDB_Setup checks if setting up the survivors table works
func TestSurvivorDB_Setup(t *testing.T) {
	os.Remove("./test.db")
	apocalypse := Open("./test.db")
	err := apocalypse.Setup()
	if err != nil {
		t.Errorf("SurvivorDB.Setup(): want: %v, got: %v", nil, err)
	}
}

// TestSurvivorDB_Save checks if saving a survivor works
func TestSurvivorDB_Save(t *testing.T) {
	os.Remove("./test.db")
	survivordb := Open("./test.db")
	err := survivordb.Setup()
	if err != nil {
		t.Errorf("SurvivorDB.Setup(): Failed to setup database")
		return
	}

	testCases := []struct {
		survivor *Survivor
		want     bool // is nil
	}{
		{
			survivor: &Survivor{
				Name:     "Jane Doe",
				Age:      1,
				Gender:   "Female",
				IdNumber: "HD138VOP34219",
				LastLocation: LastLocation{
					Longitude: 0,
					Latitude:  0,
				},
				Resources: Resources{
					Water:      2000,
					Food:       "Beef, Fish, Milk, Pasta, Ceriel",
					Medication: "Antibiotics, venteze cfc free, cough syrup",
				},
				Infected:       false,
				LastUpdateTime: time.Now(),
			},
			want: true,
		},
		{
			survivor: &Survivor{
				Name:     "John Doe",
				Age:      1,
				Gender:   "Female",
				IdNumber: "HD138VOP34220",
				LastLocation: LastLocation{
					Longitude: 0,
					Latitude:  0,
				},
				Resources: Resources{
					Water:      2000,
					Food:       "Beef, Fish, Milk, Pasta, Ceriel",
					Medication: "Antibiotics, venteze cfc free, cough syrup",
				},
				Infected:       false,
				LastUpdateTime: time.Now(),
			},
			want: true,
		},
	}

	for _, tc := range testCases {
		got := survivordb.Save(tc.survivor)
		if (got == nil) != tc.want {
			t.Errorf("SurvivorDB.Save() - %q: want: %v, got: %v", tc.survivor.Name, tc.want, (got == nil))
		}
	}
}

// TestSurvivorDB_UpdateLocation checks if updating a survivor location works
func TestSurvivorDB_UpdateLocation(t *testing.T) {
	os.Remove("./test.db")
	survivordb := Open("./test.db")
	err := survivordb.Setup()
	if err != nil {
		t.Errorf("SurvivorDB.Setup(): Failed to setup database")
		return
	}

	survivor := &Survivor{
		Name:     "Jane Doe",
		Age:      1,
		Gender:   "Female",
		IdNumber: "HD138VOP34219",
		LastLocation: LastLocation{
			Longitude: 0,
			Latitude:  0,
		},
		Resources: Resources{
			Water:      2000,
			Food:       "Beef, Fish, Milk, Pasta, Ceriel",
			Medication: "Antibiotics, venteze cfc free, cough syrup",
		},
		Infected:       false,
		LastUpdateTime: time.Now(),
	}

	err = survivordb.Save(survivor)
	if err != nil {
		t.Errorf("SurvivorDB.Save() - %q: want: %v, got: %v", survivor.Name, nil, err)
	}

	err = survivordb.UpdateLocation(survivor.IdNumber, 1, 2)
	if err != nil {
		t.Errorf("SurvivorDB.UpdateLocation() - %q: want: %v, got: %v", survivor.Name, nil, err)
	}
	newSurvivor := survivordb.GetSurvivor(survivor.IdNumber)
	if newSurvivor == nil || (newSurvivor.LastLocation.Longitude != 1 || newSurvivor.LastLocation.Latitude != 2) {
		t.Errorf("SurvivorDB.GetSurvivor() - %q: want: not nil, got: %v", survivor.Name, newSurvivor)
	}
}

// TestSurvivorDB_UpdateResource checks if updating a survivor resource works
func TestSurvivorDB_UpdateResource(t *testing.T) {
	os.Remove("./test.db")
	survivordb := Open("./test.db")
	err := survivordb.Setup()
	if err != nil {
		t.Errorf("SurvivorDB.Setup(): Failed to setup database")
		return
	}

	survivor := &Survivor{
		Name:     "Jane Doe",
		Age:      1,
		Gender:   "Female",
		IdNumber: "HD138VOP34219",
		LastLocation: LastLocation{
			Longitude: 0,
			Latitude:  0,
		},
		Resources: Resources{
			Water:      2000,
			Food:       "Beef, Fish, Milk, Pasta, Ceriel",
			Medication: "Antibiotics, venteze cfc free, cough syrup",
		},
		Infected:       false,
		LastUpdateTime: time.Now(),
	}

	err = survivordb.Save(survivor)
	if err != nil {
		t.Errorf("SurvivorDB.Save() - %q: want: %v, got: %v", survivor.Name, nil, err)
	}

	err = survivordb.UpdateResource(survivor.IdNumber,
		survivor.Resources.Water,
		survivor.Resources.Food,
		survivor.Resources.Medication,
		4000)
	if err != nil {
		t.Errorf("SurvivorDB.UpdateResource() - %q: want: %v, got: %v", survivor.Name, nil, err)
	}

	newSurvivor := survivordb.GetSurvivor(survivor.IdNumber)
	if newSurvivor == nil || (newSurvivor.Resources.Ammunition != 4000) {
		t.Errorf("SurvivorDB.GetSurvivor() - %q: want: not nil, got: %v", survivor.Name, newSurvivor)
	}
}

// TestSurvivorDB_UpdateResource checks if updating a survivor resource works
func TestSurvivorDB_UpdateInfected(t *testing.T) {
	os.Remove("./test.db")
	survivordb := Open("./test.db")
	err := survivordb.Setup()
	if err != nil {
		t.Errorf("SurvivorDB.Setup(): Failed to setup database")
		return
	}

	survivor := &Survivor{
		Name:     "Jane Doe",
		Age:      1,
		Gender:   "Female",
		IdNumber: "HD138VOP34219",
		LastLocation: LastLocation{
			Longitude: 0,
			Latitude:  0,
		},
		Resources: Resources{
			Water:      2000,
			Food:       "Beef, Fish, Milk, Pasta, Ceriel",
			Medication: "Antibiotics, venteze cfc free, cough syrup",
		},
		Infected:       false,
		LastUpdateTime: time.Now(),
	}

	err = survivordb.Save(survivor)
	if err != nil {
		t.Errorf("SurvivorDB.Save() - %q: want: %v, got: %v", survivor.Name, nil, err)
	}

	err = survivordb.UpdateInfected(survivor.IdNumber)
	if err != nil {
		t.Errorf("SurvivorDB.UpdateInfected() - %q: want: %v, got: %v", survivor.Name, nil, err)
	}

	newSurvivor := survivordb.GetSurvivor(survivor.IdNumber)
	if newSurvivor == nil || (newSurvivor.Infected != true) {
		t.Errorf("SurvivorDB.GetSurvivor() - %q: want: not nil, got: %v", survivor.Name, newSurvivor)
	}
}

// TestSurvivorDB_GetAllSurvivors checks if retrieving survivors works
func TestSurvivorDB_GetAllSurvivors(t *testing.T) {
	os.Remove("./test.db")
	survivordb := Open("./test.db")
	err := survivordb.Setup()
	if err != nil {
		t.Errorf("SurvivorDB.Setup(): Failed to setup database")
		return
	}

	testCases := []struct {
		survivor *Survivor
		want     bool
	}{
		{
			survivor: &Survivor{
				Name:     "Jane Doe",
				Age:      1,
				Gender:   "Female",
				IdNumber: "HD138VOP34219",
				LastLocation: LastLocation{
					Longitude: 0,
					Latitude:  0,
				},
				Resources: Resources{
					Water:      2000,
					Food:       "Beef, Fish, Milk, Pasta, Ceriel",
					Medication: "Antibiotics, venteze cfc free, cough syrup",
				},
				Infected:       false,
				LastUpdateTime: time.Now(),
			},
			want: true,
		},
		{
			survivor: &Survivor{
				Name:     "John Doe",
				Age:      1,
				Gender:   "Male",
				IdNumber: "HD138VOP34220",
				LastLocation: LastLocation{
					Longitude: 0,
					Latitude:  0,
				},
				Resources: Resources{
					Water:      2000,
					Food:       "Beef, Fish, Milk, Pasta, Ceriel",
					Medication: "Antibiotics, venteze cfc free, cough syrup",
				},
				Infected:       false,
				LastUpdateTime: time.Now(),
			},
			want: true,
		},
	}

	for _, tc := range testCases {
		got := survivordb.Save(tc.survivor)
		if (got == nil) != tc.want {
			t.Errorf("SurvivorDB.Save() - %q: want: %v, got: %v", tc.survivor.Name, tc.want, (got == nil))
		}
	}

	all := survivordb.GetAllSurvivors()
	if all == nil {
		t.Errorf("SurvivorDB.GetAllActions(): want: %v, got: %v", nil, err)
		return
	}
	if len(all) != 2 {
		t.Errorf("SurvivorDB.GetAllActions(): want: %v, got: %v", 2, len(all))
		return
	}
}

// TestSurvivorDB_GetSurvivors_Infected_NotInfected checks if retrieving survivors works
func TestSurvivorDB_GetSurvivors_Infected_NotInfected(t *testing.T) {
	os.Remove("./test.db")
	survivordb := Open("./test.db")
	err := survivordb.Setup()
	if err != nil {
		t.Errorf("SurvivorDB.Setup(): Failed to setup database")
		return
	}

	testCases := []struct {
		survivor *Survivor
		want     bool
	}{
		{
			survivor: &Survivor{
				Name:     "Jane Doe",
				Age:      1,
				Gender:   "Female",
				IdNumber: "HD138VOP34219",
				LastLocation: LastLocation{
					Longitude: 0,
					Latitude:  0,
				},
				Resources: Resources{
					Water:      2000,
					Food:       "Beef, Fish, Milk, Pasta, Ceriel",
					Medication: "Antibiotics, venteze cfc free, cough syrup",
				},
				Infected:       false,
				LastUpdateTime: time.Now(),
			},
			want: true,
		},
		{
			survivor: &Survivor{
				Name:     "John Doe",
				Age:      1,
				Gender:   "Male",
				IdNumber: "HD138VOP34220",
				LastLocation: LastLocation{
					Longitude: 0,
					Latitude:  0,
				},
				Resources: Resources{
					Water:      2000,
					Food:       "Beef, Fish, Milk, Pasta, Ceriel",
					Medication: "Antibiotics, venteze cfc free, cough syrup",
				},
				Infected:       true,
				LastUpdateTime: time.Now(),
			},
			want: true,
		},
	}

	for _, tc := range testCases {
		got := survivordb.Save(tc.survivor)
		if (got == nil) != tc.want {
			t.Errorf("SurvivorDB.Save() - %q: want: %v, got: %v", tc.survivor.Name, tc.want, (got == nil))
		}
	}

	infected := survivordb.GetSurvivors(true)
	if infected == nil || len(infected) != 1 || infected[0].IdNumber != "HD138VOP34220" {
		t.Errorf("SurvivorDB.GetSurvivors(true): want: %v, got: %v", nil, err)
		return
	}

	healthy := survivordb.GetSurvivors(false)
	if healthy == nil || len(healthy) != 1 || healthy[0].IdNumber != "HD138VOP34219" {
		t.Errorf("SurvivorDB.GetSurvivors(false): want: %v, got: %v", nil, err)
		return
	}
}

// TestSurvivorDB_CountSurvivors_Infected_NotInfected checks if counting survivors works
func TestSurvivorDB_CountSurvivors_Infected_NotInfected(t *testing.T) {
	os.Remove("./test.db")
	survivordb := Open("./test.db")
	err := survivordb.Setup()
	if err != nil {
		t.Errorf("SurvivorDB.Setup(): Failed to setup database")
		return
	}

	testCases := []struct {
		survivor *Survivor
		want     bool
	}{
		{
			survivor: &Survivor{
				Name:     "Jane Doe",
				Age:      1,
				Gender:   "Female",
				IdNumber: "HD138VOP34219",
				LastLocation: LastLocation{
					Longitude: 0,
					Latitude:  0,
				},
				Resources: Resources{
					Water:      2000,
					Food:       "Beef, Fish, Milk, Pasta, Ceriel",
					Medication: "Antibiotics, venteze cfc free, cough syrup",
				},
				Infected:       false,
				LastUpdateTime: time.Now(),
			},
			want: true,
		},
		{
			survivor: &Survivor{
				Name:     "John Doe",
				Age:      1,
				Gender:   "Male",
				IdNumber: "HD138VOP34220",
				LastLocation: LastLocation{
					Longitude: 0,
					Latitude:  0,
				},
				Resources: Resources{
					Water:      2000,
					Food:       "Beef, Fish, Milk, Pasta, Ceriel",
					Medication: "Antibiotics, venteze cfc free, cough syrup",
				},
				Infected:       true,
				LastUpdateTime: time.Now(),
			},
			want: true,
		},
		{
			survivor: &Survivor{
				Name:     "Jill Doe",
				Age:      1,
				Gender:   "Female",
				IdNumber: "HD138VOP34220",
				LastLocation: LastLocation{
					Longitude: 0,
					Latitude:  0,
				},
				Resources: Resources{
					Water:      2000,
					Food:       "Beef, Fish, Milk, Pasta, Ceriel",
					Medication: "Antibiotics, venteze cfc free, cough syrup",
				},
				Infected:       true,
				LastUpdateTime: time.Now(),
			},
			want: true,
		},
	}

	for _, tc := range testCases {
		got := survivordb.Save(tc.survivor)
		if (got == nil) != tc.want {
			t.Errorf("SurvivorDB.Save() - %q: want: %v, got: %v", tc.survivor.Name, tc.want, (got == nil))
		}
	}

	infected := survivordb.CountSurvivors(true)
	if infected != 2 {
		t.Errorf("SurvivorDB.CountSurvivors(true): want: %v, got: %v", 1, infected)
		return
	}

	healthy := survivordb.CountSurvivors(false)
	if healthy != 1 {
		t.Errorf("SurvivorDB.CountSurvivors(false): want: %v, got: %v", 1, healthy)
		return
	}
}
