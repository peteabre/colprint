package colprint

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"os"
	"errors"
)

type UnitTests struct {
	suite.Suite
}

func TestRunUnitTests(t *testing.T) {
	suite.Run(t, new(UnitTests))
}

func (s *UnitTests) TestColumns_Len() {
	cols := make(columns, 0)
	col := column{}

	s.Equal(0, cols.Len())

	cols = append(cols, col)
	cols = append(cols, col)
	s.Equal(2, cols.Len())
}

func (s *UnitTests) TestColumns_Swap() {
	cols := make(columns, 2)
	col1 := column{FieldName: "Name1", Label: "Label1", Order: 1}
	col2 := column{FieldName: "Name2", Label: "Label2", Order: 2}
	cols[0] = col1
	cols[1] = col2

	cols.Swap(0, 1)
	s.True(cols[0] == col2)
	s.True(cols[1] == col1)
}

func (s *UnitTests) TestColumns_Less() {
	cols := make(columns, 2)
	col1 := column{FieldName: "Name1", Label: "Label1", Order: 10}
	col2 := column{FieldName: "Name2", Label: "Label2", Order: 20}
	cols[0] = col1
	cols[1] = col2
	s.True(cols.Less(0, 1))
	s.False(cols.Less(1, 0))

	cols[0].Order = 20
	s.False(cols.Less(0, 1))
	s.False(cols.Less(1, 0))
}

func (s *UnitTests) TestCPrinter_init() {
	cp := cPrinter{}
	s.Nil(cp.values)
	s.Nil(cp.cols)
	cp.init()
	s.NotNil(cp.values)
	s.NotNil(cp.cols)
}

func (s *UnitTests) TestCPrinter_initColumn() {
	cp := cPrinter{}
	cp.init()
	col := column{"name", "label", 2}
	val := cp.values[col]
	s.Nil(val)
	cp.initColumn(col)
	val = cp.values[col]
	s.NotNil(val)
}

func (s *UnitTests) TestCPrinter_add() {
	cp := cPrinter{config: createDefaultConfig()}
	d := DummyData{Name: "name", Description: "description", Version: float32(35)}
	s.NotPanics(func() {
		s.NoError(cp.add(d))
	})

	cp = cPrinter{config: createDefaultConfig()}
	e := Errornous{Error: errors.New("Error")}
	s.NotPanics(func() {
		s.Error(cp.add(e))
	})
}

func (s *UnitTests) TestMergeConfig() {
	defaultConf := createDefaultConfig()

	mergedConf := mergeConfig(defaultConf, &Config{})
	s.Equal(*defaultConf.FloatPrecision, *mergedConf.FloatPrecision)
	s.Equal(*defaultConf.MaxPrintedSliceItems, *mergedConf.MaxPrintedSliceItems)

	mpsli := 10
	mergedConf = mergeConfig(defaultConf, &Config{MaxPrintedSliceItems: &mpsli})
	s.Equal(*defaultConf.FloatPrecision, *mergedConf.FloatPrecision)
	s.Equal(mpsli, *mergedConf.MaxPrintedSliceItems)

	fp := 10
	mergedConf = mergeConfig(defaultConf, &Config{FloatPrecision: &fp})
	s.Equal(fp, *mergedConf.FloatPrecision)
	s.Equal(*defaultConf.MaxPrintedSliceItems, *mergedConf.MaxPrintedSliceItems)
}

func (s *UnitTests) TestCreateDefaultConfig() {
	c := createDefaultConfig()
	s.NotNil(c)
	s.NotNil(c.MaxPrintedSliceItems)
	s.NotNil(c.FloatPrecision)
	s.Equal(2, *c.FloatPrecision)
	s.Equal(3, *c.MaxPrintedSliceItems)
}

func (s *UnitTests) TestDefaultPrint() {
	persons := []Person{
		{
			FirstName: "Ola",
			LastName:  "Nordmann",
			Age:       35,
			Groups:    []string{"group1", "group2", "group3", "group4"},
		},
		{
			FirstName: "Kari",
			LastName:  "Nordmann",
			Age:       37,
			Groups:    []string{"group1", "group2", "group3"},
		},
	}
	s.NotPanics(func() {
		s.NoError(DefaultPrint(persons))
	})

	errs := []Errornous{{Error: errors.New("Error")}}
	s.NotPanics(func() {
		s.Error(DefaultPrint(errs))
	})
}

func (s *UnitTests) TestDefaultFPrint() {
	age := 40
	d := DummyData{Age: &age, Name: "name", Description: "description", Version: float32(35)}
	s.NotPanics(func() {
		s.NoError(DefaultFprint(os.Stdout, d))
	})

	err := Errornous{Error: errors.New("Error")}
	s.NotPanics(func() {
		s.Error(DefaultFprint(os.Stdout, err))
	})
}

type Errornous struct {
	Error error `colprint:"Error,a"`
}

type DummyData struct {
	Name        string  `colprint:"Name,3"`
	Description string  `colprint:"Description"`
	Valid       bool    `colprint:"Valid"`
	Age         *int    `colprint:"Age,1"`
	Version     float32 `colprint:"Version,2"`
}

type Person struct {
	FirstName string    `colprint:"First name,1"`
	LastName  string    `colprint:"Last name,2"`
	Age       int       `colprint:"Age,3"`
	Groups    []string  `colprint:"Groups,4"`
	Address   string    `colprint:""`
	Address2  string    `colprint:"-"`
	Spouse    *Person   `colprint:"Spouse"`
	Data      DummyData `colprint:"Data"`
}
