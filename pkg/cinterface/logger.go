package cinterface

import c "vm/pkg/constants"

type Logger interface {
	Init()

	Debug(cat c.Category, sub c.SubCategory, msg string, extra map[c.ExtraKey]interface{})
	Debugf(templateName string, args ...interface{})

	Info(cat c.Category, sub c.SubCategory, msg string, extra map[c.ExtraKey]interface{})
	Infof(templateName string, args ...interface{})

	Warn(cat c.Category, sub c.SubCategory, msg string, extra map[c.ExtraKey]interface{})
	Warnf(templateName string, args ...interface{})

	Error(cat c.Category, sub c.SubCategory, msg string, extra map[c.ExtraKey]interface{})
	Errorf(templateName string, args ...interface{})

	Fatal(cat c.Category, sub c.SubCategory, msg string, extra map[c.ExtraKey]interface{})
	Fatalf(templateName string, args ...interface{})
}
