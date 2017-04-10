package colprint

import (
	"github.com/stretchr/testify/suite"
	"testing"
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

	cols.Swap(0,1)
	s.True(cols[0] == col2)
	s.True(cols[1] == col1)
}

func (s *UnitTests) TestColumns_Less() {
	cols := make(columns, 2)
	col1 := column{FieldName: "Name1", Label: "Label1", Order: 10}
	col2 := column{FieldName: "Name2", Label: "Label2", Order: 20}
	cols[0] = col1
	cols[1] = col2
	s.True(cols.Less(0,1))
	s.False(cols.Less(1,0))

	cols[0].Order=20
	s.False(cols.Less(0,1))
	s.False(cols.Less(1,0))
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

func (s *UnitTests) TestFprint() {
	data := []simpleStruct{{Name: "name", Description: "description", Version:float32(35)}, {Name: "Navn", Description: "beskrivelse"}}
	s.NotPanics(func() {
		Print(data)
	})


	d := simpleStruct{Name: "name", Description: "description", Version:float32(35)}
	s.NotPanics(func() {
		Print(d)
	})
}

type simpleStruct struct {
	Name        string `column:"Name,3"`
	Description string `column:"Tittentei"`
	Valid       bool `column:"Valid"`
	Age         int `column:"Age,1"`
	Version float32 `column:"Version,2"`
}
