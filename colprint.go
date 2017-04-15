package colprint

import (
	"io"
	"os"
	"reflect"
	"strconv"
	"fmt"
	"strings"
	"math"
	"sort"
	"github.com/ryanuber/columnize"
)

const TagName = "colprint"

const (
	TagValueEmpty    = ""
	TagValueSkip     = "-"
	TagValueTraverse = "=>"
)

// Config holds configuration used when printing columns
type Config struct {
	// MaxPrintedSliceItems represents the maximum number og slice items to list.
	MaxPrintedSliceItems *int
	// FloatPrecision represents the precision used when printing floats.
	FloatPrecision *int
}

// Print prints a struct or slice of structs to stdout using default config
func Print(s interface{}) error {
	return Fprint(os.Stdout, s)
}

// Fprint prints struct or slice to provided io.Writer using provided config.
// If config is nil, default config will be used.
func Fprint(w io.Writer, s interface{}, c ... *Config) error {
	var conf *Config
	if len(c) > 0 {
		conf = c[0]
	}
	cp := cPrinter{config: mergeConfig(createDefaultConfig(), conf)}
	val := reflect.ValueOf(s)
	kind := val.Kind()

	// If its a pointer, do an indirect...
	if kind == reflect.Ptr {
		val = reflect.Indirect(reflect.ValueOf(s))
		kind = val.Kind()
	}

	// Check if s is a slice/array or not
	if kind == reflect.Slice || kind == reflect.Array {
		// add each item in slice to cPrinter
		for i := 0; i < val.Len(); i ++ {
			if err := cp.add(val.Index(i).Interface()); err != nil {
				return err
			}
		}
	} else {
		// add the item to cPrinter
		if err := cp.add(val.Interface()); err != nil {
			return err
		}
	}
	// Print to provided Writer
	cp.fprint(w)
	return nil
}

// column represents a column that will be printed by cPrinter
type column struct {
	fieldIndex *[]int
	label      string
	order      int
}

// columns is a sortable list of column structs
type columns []column

func (s columns) Len() int {
	return len(s)
}

func (s columns) Less(i, j int) bool {
	return s[i].order < s[j].order
}

func (s columns) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// cPrinter is the data structure used to print columns
type cPrinter struct {
	// columns (ordered)
	cols columns
	// Map containing values for all columns
	values map[column][]string
	// Keeps track of number of items appended to the ColPrinter
	itemCount int
	// Configuration for the printer
	config *Config
}

// add adds a struct columns and values
func (cp *cPrinter) add(s interface{}) error {
	// Init columns if it's not already done
	if cp.cols == nil {
		cp.init()
		err := cp.findColumns(s)
		if err != nil {
			return err
		}
		for _, col := range cp.cols {
			cp.initColumn(col)
		}
	}
	// Add values
	for _, col := range cp.cols {
		v := reflect.ValueOf(s)
		val := cp.valueOf(v.FieldByIndex(*col.fieldIndex).Interface())
		cp.values[col] = append(cp.values[col], val)
	}
	cp.itemCount++
	return nil
}

//fprint prints the columns to the provided io.Writer.
func (cp *cPrinter) fprint(w io.Writer) {
	// Add header line
	str := []string{}
	headers := ""
	for i, col := range cp.cols {
		headers += col.label
		if i != len(cp.cols)-1 {
			headers += "|"
		}
	}
	str = append(str, headers)

	// Add a line for each item appended
	for i := 0; i < cp.itemCount; i++ {
		vals := ""
		for j, col := range cp.cols {
			vals += cp.values[col][i]
			if j != len(cp.cols)-1 {
				vals += "|"
			}
		}
		str = append(str, vals)
	}
	// Print to given Writer
	fmt.Fprint(w, columnize.SimpleFormat(str)+"\n")
}

// init initializes the array containing columns, and the map containing the values for each column.
func (cp *cPrinter) init() {
	cp.cols = columns{}
	cp.values = make(map[column][]string)
}

// initColumn initializes the array containing column values.
func (cp *cPrinter) initColumn(col column) {
	cp.values[col] = make([]string, 0)
}

// findColumns extracts which columns should be printed and adds them to columns. Returns an error if any field
// contains a incomplete tag.
func (cp *cPrinter) findColumns(s interface{}, fieldIndex ... int) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}

	for i := 0; i < v.NumField(); i++ {
		fIndex := append(fieldIndex, i)
		field := v.Type().Field(i)
		tag := field.Tag.Get(TagName)
		switch tag {
		case TagValueEmpty, TagValueSkip:
			// Do nothing
		case TagValueTraverse:
			if err := cp.traverseStruct(v.FieldByIndex([]int{i}), fIndex...); err != nil {
				return err
			}
		default:
			if err := cp.appendColumn(tag, field, &fIndex); err != nil {
				return err
			}
		}
	}
	sort.Sort(cp.cols)
	return nil
}

// appendColumn appends a tagged field to the list of columns.
func (cp *cPrinter) appendColumn(tag string, field reflect.StructField, fieldIndex *[]int) error {
	args := strings.Split(tag, ",")
	switch len(args) {
	case 1:
		cp.cols = append(cp.cols, column{fieldIndex, args[0], math.MaxInt32})
	case 2:
		order, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("Invalid order on field %s", field.Name)
		}
		cp.cols = append(cp.cols, column{fieldIndex, args[0], order})
	default:
		return fmt.Errorf("Invalid number of tag arguments on field %s", field.Name)
	}
	return nil
}

// traverseStruct traverses field of kind reflect.Struct and finds columns
func (cp *cPrinter) traverseStruct(v reflect.Value, fieldIndex ... int) error {
	val := reflect.ValueOf(v.Interface())
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(v)
	}
	if val.Kind() == reflect.Struct {
		err := cp.findColumns(v.Interface(), fieldIndex...)
		if err != nil {
			return err
		}

	}
	return nil
}

// valueOf returns a string representation of a field.
func (cp *cPrinter) valueOf(i interface{}) string {
	v := reflect.ValueOf(i)
	kind := v.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Array, reflect.Slice:
		return cp.valueOfSlice(i)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', *cp.config.FloatPrecision, 64)
	case reflect.String:
		return v.String()
	case reflect.Ptr:
		if v.IsNil() {
			return ""
		}
		return cp.valueOf(reflect.Indirect(v).Interface())
	}
	return "<Unsupported kind:" + kind.String() + ">"
}

// valueOfSlice returns a string representation of the values in a slice field.
// Returns a maximum of Config.MaxPrintedSliceItems.
func (cp *cPrinter) valueOfSlice(s interface{}) string {
	sliceValue := reflect.ValueOf(s)
	values := ""
	for i := 0; i < sliceValue.Len(); i++ {
		values += cp.valueOf(sliceValue.Index(i).Interface())
		if i == *cp.config.MaxPrintedSliceItems-1 && sliceValue.Len() > *cp.config.MaxPrintedSliceItems {
			values += ",..."
			break
		} else if i < sliceValue.Len()-1 {
			values += ", "
		}
	}
	return values
}

// createDefaultConfig creates a default configuration.
func createDefaultConfig() *Config {
	dMPSI := 3
	dFP := 2
	return &Config{
		MaxPrintedSliceItems: &dMPSI,
		FloatPrecision:       &dFP,
	}
}

// mergeConfig merges the second argument config into the first.
func mergeConfig(a, c *Config) *Config {
	if c != nil {
		if c.MaxPrintedSliceItems != nil {
			*a.MaxPrintedSliceItems = *c.MaxPrintedSliceItems
		}

		if c.FloatPrecision != nil {
			*a.FloatPrecision = *c.FloatPrecision
		}
	}
	return a
}
