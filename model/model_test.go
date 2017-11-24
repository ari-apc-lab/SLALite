/*
   Copyright 2017 Atos

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package model

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

var pr = Provider{Id: "id", Name: "name"}
var cl = Provider{Id: "id", Name: "name"}

func TestProviders(t *testing.T) {
	p := Provider{Id: "id", Name: "name"}
	checkNumber(t, &p, 0)

	if p.GetId() != p.Id {
		t.Errorf("Provider.Id and Provider.GetId() do not match")
	}

	p = Provider{Id: "", Name: "name"}
	checkNumber(t, &p, 1)

	p = Provider{Id: "id", Name: ""}
	checkNumber(t, &p, 1)

	p = Provider{Id: "", Name: ""}
	checkNumber(t, &p, 2)
}

func TestAssessment(t *testing.T) {
	a := Assessment{FirstExecution: time.Now(), LastExecution: time.Now()}
	checkNumber(t, &a, 0)
}

func TestGuarantee(t *testing.T) {
	g := Guarantee{Name: "name", Constraint: "a LT 10"}
	checkNumber(t, &g, 0)

	g = Guarantee{Name: "", Constraint: "a LT 10"}
	checkNumber(t, &g, 1)

	g = Guarantee{Name: "name", Constraint: ""}
	checkNumber(t, &g, 1)

}

func TestAgreementText(t *testing.T) {
	at := AgreementText{Id: "id", Name: "name", Provider: pr, Client: cl}
	checkNumber(t, &at, 0)

	at = AgreementText{Id: "", Name: "name", Provider: pr, Client: cl}
	checkNumber(t, &at, 1)

	at = AgreementText{Id: "id", Name: "", Provider: pr, Client: cl}
	checkNumber(t, &at, 1)

	at = AgreementText{Id: "id", Name: "name", Client: cl}
	checkNumber(t, &at, 2)

	at = AgreementText{Id: "id", Name: "name", Provider: pr}
	checkNumber(t, &at, 2)

	at = AgreementText{
		Id:       "id",
		Name:     "name",
		Provider: pr,
		Client:   cl,
		Guarantees: []Guarantee{
			Guarantee{Name: ""},
		},
	}
	checkNumber(t, &at, 2)
}

func TestAgreement(t *testing.T) {

	a := Agreement{
		Id:         "id",
		Name:       "name",
		State:      STOPPED,
		Assessment: Assessment{},
		Text: AgreementText{
			Id:       "id",
			Name:     "name",
			Provider: pr,
			Client:   cl,
		},
	}
	checkNumber(t, &a, 0)

	if a.GetId() != a.Id {
		t.Errorf("Agreement.Id and Agreement.GetId() do not match")
	}

	a.Id = ""
	a.Text.Id = ""
	a.Name = "name"
	a.Text.Name = "name"
	checkNumber(t, &a, 2) // one error per empty id

	a.Id = "id"
	a.Text.Id = "id"
	a.Name = ""
	a.Text.Name = ""
	checkNumber(t, &a, 2) // one error per empty name

	a.Id = "id1"
	a.Text.Id = "id2"
	a.Name = "name"
	a.Text.Name = "name"
	checkNumber(t, &a, 1)

	a.Id = "id"
	a.Text.Id = "id"
	a.Name = "name1"
	a.Text.Name = "name2"
	checkNumber(t, &a, 1)
}

func TestStates(t *testing.T) {
	a := Agreement{State: STOPPED}
	if !a.IsStopped() {
		t.Error("Agreement should be stopped")
	}
	a = Agreement{State: STARTED}
	if !a.IsStarted() {
		t.Error("Agreement should be started")
	}
	a = Agreement{State: TERMINATED}
	if !a.IsTerminated() {
		t.Error("Agreement should be terminated")
	}
}

func TestProviderSerialization(t *testing.T) {
	var p Provider

	s := `{"id": "id", "name": "name" }`
	json.NewDecoder(strings.NewReader(s)).Decode(&p)
	checkNumber(t, &p, 0)
}

func TestAgreementSerialization(t *testing.T) {
	var a1, a2 Agreement
	var err error

	s := `{
		"id": "id", 
		"name": "name", 
		"text": {
			"id": "id",
			"name": "name",
			"provider": { "id": "pr-id", "name": "pr-name" },
			"client": { "id": "cl-id", "name": "cl-name" }
		},
		"state": "stopped"
	}`
	err = json.NewDecoder(strings.NewReader(s)).Decode(&a1)
	if err != nil {
		t.Fatalf("Error decoding %v", err)
	}
	checkNumber(t, &a1, 0)

	// state is empty. Validate should normalize to STOPPED
	s = `{
		"id": "id", 
		"name": "name", 
		"text": {
			"id": "id",
			"name": "name",
			"provider": { "id": "pr-id", "name": "pr-name" },
			"client": { "id": "cl-id", "name": "cl-name" }
		}
	}`
	err = json.NewDecoder(strings.NewReader(s)).Decode(&a2)
	if err != nil {
		t.Fatalf("Error decoding %v", err)
	}
	checkNumber(t, &a2, 0)
	if a2.State != STOPPED {
		t.Errorf("State=%s is not STOPPED", a2.State)
	}

}

func checkNumber(t *testing.T, v Validable, expected int) {

	if errs := v.Validate(); len(errs) != expected {
		t.Errorf("Error validating %s%v. Errors = %v", reflect.TypeOf(v), v, errs)
	}
}