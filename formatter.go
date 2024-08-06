package formatter

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Formatter - logrus formatter, implements logrus.Formatter
type Formatter struct {
	// FieldsOrder - default: fields sorted alphabetically
	FieldsOrder []string

	// TimestampFormat - default: time.StampMilli = "Jan _2 15:04:05.000"
	TimestampFormat string

	// HideKeys - show [fieldValue] instead of [fieldKey:fieldValue]
	HideKeys bool

	// NoColors - disable colors
	NoColors bool

	// NoFieldsColors - apply colors only to the level, default is level + fields
	NoFieldsColors bool

	// NoFieldsSpace - no space between fields
	NoFieldsSpace bool

	// ShowFullLevel - show a full level [WARNING] instead of [WARN]
	ShowFullLevel bool

	// NoUppercaseLevel - no upper case for level value
	NoUppercaseLevel bool

	// TrimMessages - trim whitespaces on messages
	TrimMessages bool

	// CallerFirst - print caller info first
	CallerFirst bool

	// CustomCallerFormatter - set custom formatter for caller info
	CustomCallerFormatter func(*runtime.Frame) string

	// Module Name
	ModuleName string
}

// Format an log entry
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	// output buffer
	b := &bytes.Buffer{}

	// write time
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.StampMilli
	}
	b.WriteString(entry.Time.Format(timestampFormat))

	// write caller (first)
	if f.CallerFirst {
		f.writeCaller(b, entry)
	}

	// write level
	levelColor := getColorByLevel(entry.Level)
	f.writeLevel(b, entry, levelColor)

	// write module
	if f.ModuleName != "" {
		f.writeModule(b, &entry.Data, levelColor)
	}

	// write fields
	if f.FieldsOrder == nil {
		f.writeFields(b, &entry.Data, levelColor)
	} else {
		f.writeOrderedFields(b, &entry.Data, levelColor)
	}

	// write message
	b.WriteByte(' ')
	if f.TrimMessages {
		b.WriteString(strings.TrimSpace(entry.Message))
	} else {
		b.WriteString(entry.Message)
	}

	// write caller (not first)
	if !f.CallerFirst {
		f.writeCaller(b, entry)
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *Formatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		if f.CustomCallerFormatter != nil {
			callStr := f.CustomCallerFormatter(entry.Caller)
			if len(callStr) > 0 {
				b.WriteByte(' ')
				b.WriteString(callStr)
			}
		} else {
			b.WriteString(" (")
			b.WriteString(entry.Caller.File)
			b.WriteByte(':')
			b.WriteString(strconv.Itoa(entry.Caller.Line))
			b.WriteByte(' ')
			b.WriteString(entry.Caller.Function)
			b.WriteByte(')')
		}
	}
}

func (f *Formatter) writeLevel(b *bytes.Buffer, entry *logrus.Entry, levelColor int) {
	var levelStr string
	if f.NoUppercaseLevel {
		levelStr = entry.Level.String()
	} else {
		levelStr = strings.ToUpper(entry.Level.String())
	}

	b.WriteByte(' ')
	f.startColor(b, levelColor, true)
	b.WriteByte('[')

	if f.ShowFullLevel {
		b.WriteString(levelStr)
	} else {
		b.WriteString(levelStr[:4])
	}

	b.WriteByte(']')
	f.stopColor(b, true)
}

func (f *Formatter) writeModule(b *bytes.Buffer, entryData *logrus.Fields, levelColor int) {
	if len(*entryData) > 0 {
		if o, ok := (*entryData)[f.ModuleName]; ok {
			b.WriteByte(' ')
			f.startColor(b, levelColor, !f.NoFieldsColors)
			fmt.Fprintf(b, "[%v]", o)
			f.stopColor(b, !f.NoFieldsColors)
		}
	}
}

func (f *Formatter) writeFields(b *bytes.Buffer, entryData *logrus.Fields, levelColor int) {
	length := len(*entryData)
	if length == 0 {
		return
	}
	if length == 1 {
		if _, ok := (*entryData)[f.ModuleName]; ok {
			return
		}
	}

	b.WriteByte(' ')
	f.startColor(b, levelColor, !f.NoFieldsColors)
	f.writeFieldsWithFilter(b, entryData, nil)
	f.stopColor(b, !f.NoFieldsColors)
}

func (f *Formatter) writeOrderedFields(b *bytes.Buffer, entryData *logrus.Fields, levelColor int) {
	length := len(*entryData)
	if _, ok := (*entryData)[f.ModuleName]; length == 0 || (length == 1 && ok) {
		// no fields to write
		return
	}

	b.WriteByte(' ')
	f.startColor(b, levelColor, !f.NoFieldsColors)

	// write ordered fields
	orderFieldsMap := map[string]struct{}{}
	for _, key := range f.FieldsOrder {
		if _, ok := (*entryData)[key]; !ok {
			continue
		}
		if len(orderFieldsMap) > 0 && !f.NoFieldsSpace {
			b.WriteByte(' ')
		}
		f.writeField(b, entryData, key)
		orderFieldsMap[key] = struct{}{}
		length--
	}

	// write remaining fields
	if _, ok := (*entryData)[f.ModuleName]; length > 1 || (length == 1 && !ok) {
		if len(orderFieldsMap) > 0 && !f.NoFieldsSpace {
			b.WriteByte(' ')
		}
		f.writeFieldsWithFilter(b, entryData, &orderFieldsMap)
	}

	f.stopColor(b, !f.NoFieldsColors)
}

func (f *Formatter) writeFieldsWithFilter(b *bytes.Buffer, entryData *logrus.Fields, filter *map[string]struct{}) {
	keys := make([]string, 0, len(*entryData))
	for key := range *entryData {
		if filter != nil {
			if _, ok := (*filter)[key]; ok {
				continue
			}
		}
		if key == f.ModuleName {
			continue
		}
		keys = append(keys, key)
	}

	sort.Strings(keys)
	i := 0
	for _, key := range keys {
		if i > 0 && !f.NoFieldsSpace {
			b.WriteByte(' ')
		}
		f.writeField(b, entryData, key)
		i++
	}
}

func (f *Formatter) writeField(b *bytes.Buffer, entryData *logrus.Fields, field string) {
	if f.HideKeys {
		fmt.Fprintf(b, "[%v]", (*entryData)[field])
	} else {
		fmt.Fprintf(b, "[%s:%v]", field, (*entryData)[field])
	}
}

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return colorGray
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}

var mapColors = map[int]string{
	colorRed:    "\x1b[31m",
	colorYellow: "\x1b[33m",
	colorBlue:   "\x1b[36m",
	colorGray:   "\x1b[37m",
}

func (f *Formatter) startColor(b *bytes.Buffer, levelColor int, flag bool) {
	if !f.NoColors && flag {
		var color = mapColors[levelColor]
		b.WriteString(color)
	}
}

func (f *Formatter) stopColor(b *bytes.Buffer, flag bool) {
	if !f.NoColors && flag {
		b.WriteString("\x1b[0m")
	}
}
