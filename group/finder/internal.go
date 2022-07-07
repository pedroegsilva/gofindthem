package finder

import (
	"fmt"
	"reflect"
	"strings"
)

func (rf *GroupFinder) getRulesInfo(
	data interface{},
	fieldName string,
	includePaths []string,
	excludePaths []string,
	matchedExpByFieldByByTag map[string]map[string]map[string]struct{},
) (err error) {
	t := reflect.TypeOf(data)

	val := reflect.ValueOf(data)

	switch val.Kind() {
	case reflect.String:
		if !isValidateFieldPath(fieldName, includePaths, excludePaths) {
			return
		}
		expRes, err := rf.findthem.ProcessText(val.String())
		if err != nil {
			return err
		}

		for _, er := range expRes {
			if _, ok := matchedExpByFieldByByTag[er.Tag]; !ok {
				matchedExpByFieldByByTag[er.Tag] = make(map[string]map[string]struct{})
			}
			if _, ok := matchedExpByFieldByByTag[er.Tag][fieldName]; !ok {
				matchedExpByFieldByByTag[er.Tag][fieldName] = make(map[string]struct{})
			}
			matchedExpByFieldByByTag[er.Tag][fieldName][er.ExpresionStr] = struct{}{}
		}

	case reflect.Struct:
		numField := t.NumField()

		for i := 0; i < numField; i++ {
			structField := t.Field(i)
			fn := structField.Name
			if fieldName != "" {
				fn = fieldName + "." + fn
			}
			if !val.Field(i).CanInterface() {
				continue
			}
			err := rf.getRulesInfo(val.Field(i).Interface(), fn, includePaths, excludePaths, matchedExpByFieldByByTag)
			if err != nil {
				return err
			}
		}

	case reflect.Map:
		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key()
			if k.Type().Kind() != reflect.String {
				break
			}

			v := iter.Value()
			fn := k.String()
			if fieldName != "" {
				fn = fieldName + "." + fn
			}
			if !v.CanInterface() {
				continue
			}
			err := rf.getRulesInfo(v.Interface(), fn, includePaths, excludePaths, matchedExpByFieldByByTag)
			if err != nil {
				return err
			}
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			fn := fmt.Sprintf("index(%d)", i)
			if fieldName != "" {
				fn = fieldName + "." + fn
			}
			if !val.Index(i).CanInterface() {
				continue
			}
			err := rf.getRulesInfo(val.Index(i).Interface(), fn, includePaths, excludePaths, matchedExpByFieldByByTag)
			if err != nil {
				return err
			}
		}
	}

	return
}

// isValidateFieldPath returns true if the field path is valid for tagging
func isValidateFieldPath(fieldPath string, includePaths []string, excludePaths []string) bool {
	if len(excludePaths) > 0 {
		for _, excP := range excludePaths {
			if strings.HasPrefix(fieldPath, excP) {
				return false
			}
		}
	}

	if len(includePaths) > 0 {
		for _, incP := range includePaths {
			if strings.HasPrefix(fieldPath, incP) {
				return true
			}
		}
		return false
	}

	return true
}
