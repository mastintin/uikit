package uikit

import (
	"fmt"
	"time"
)

// ValToStr convierte cualquier valor a su representación en string mediante fmt.Sprint.
func ValToStr(xval any) string {
	return fmt.Sprint(xval)
}

// dateToStrBase maneja los formatos
func dateToStrBase(dfecha time.Time, lEuropa bool) string {
	if lEuropa {
		return dfecha.Format("02/01/2006")
	}
	return dfecha.Format("01/02/2006")
}

// DateToStrEu equivale a tu función por defecto (formato europeo DD/MM/AAAA)
func DateToStrEu(dfecha time.Time) string {
	return dateToStrBase(dfecha, true)
}

// DateToStrUs formatea la fecha en formato estadounidense (MM/DD/AAAA)
func DateToStrUs(dfecha time.Time) string {
	return dateToStrBase(dfecha, false)
}
