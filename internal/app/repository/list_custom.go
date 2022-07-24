package repository

import (
	"encoding/json"
	"fmt"
	"reflect"
	"restapi/internal/app/model"
	"restapi/internal/db/postgres"
	"sort"
	"strings"

	"gorm.io/gorm"
)

type CustomRepo interface {
	List(req model.RequestDataTable, dataStruct interface{}, rawQuwey string) (*model.ResultDataTable, error)
}

func NewCustom(postgres postgres.Client) CustomRepo {
	return &listCustomRepo{postgres}
}

type listCustomRepo struct {
	postgres postgres.Client
}

func (r *listCustomRepo) List(req model.RequestDataTable, dataStruct interface{}, rawQuery string) (*model.ResultDataTable, error) {
	query := r.postgres.Conn().Table(rawQuery)

	for f, v := range req.Additional {
		inputProcessNew(query, f, v)
	}

	temp, err := dataTable(query, req, dataStruct)
	if err != nil {
		return nil, err
	}

	return temp, nil
}

func inputProcessNew(base *gorm.DB, column string, value interface{}) {
	m := reflect.ValueOf(value)
	switch m.Kind() {
	case reflect.Map:
		val := convertToMap(m.Interface())
		k := m.MapKeys()
		if k[0].String() != "from" && k[1].String() != "to" {
			slc := []string{}
			// var stmt string
			for f, req := range val {
				r := reflect.ValueOf(req)
				if f == "range" && r.Len() >= 1 {
					for i := 0; i < r.Len(); i++ {
						d := reflect.ValueOf(r.Index(i).Interface())
						data := convertToMap(r.Index(i).Interface())
						keys := d.MapKeys()
						sort.Slice(keys, func(i, j int) bool {
							return keys[i].String() < keys[j].String()
						})

						if data[keys[0].String()] != "" && data[keys[1].String()] != "" {
							slc = append(slc, fmt.Sprintf("%s between '%v' and '%v'", column, data[keys[0].String()], data[keys[1].String()]))
						} else if data[keys[0].String()] != "" && data[keys[1].String()] == "" {
							slc = append(slc, fmt.Sprintf("%v = '%v'", column, data[keys[0].String()]))
						} else if data[keys[0].String()] == "" && data[keys[1].String()] != "" {
							slc = append(slc, fmt.Sprintf("%v <= '%v'", column, data[keys[1].String()]))
						}
					}
				} else if f == "multiple" && r.Len() >= 1 {
					slc = append(slc, fmt.Sprintf("%v in ?", column))
				}
			}

			if len(slc) == 1 {
				if strings.Contains(slc[0], "in") {
					base.Where(slc[0], val["multiple"])
				} else {
					base.Where(slc[0])
				}
			} else if len(slc) > 1 {
				str := strings.Join(slc, " or ")

				if strings.Contains(str, "in") {
					base.Where(str, val["multiple"])
				} else {
					base.Where(str)
				}
			}
		} else {
			if val[k[0].String()] != "" && val[k[1].String()] != "" {
				base.Where(column+" between ? and ?", val[k[0].String()], val[k[1].String()])
			} else if val[k[0].String()] != "" && val[k[1].String()] == "" {
				base.Where(column+" = ?", val[k[0].String()])
			} else if val[k[0].String()] == "" && val[k[1].String()] != "" {
				base.Where(column+" = ?", val[k[1].String()])
			}
		}
	case reflect.Slice:
		if m.Len() > 1 {
			var data []interface{}

			for i := 0; i < m.Len(); i++ {
				if m.Index(i).Interface() != "" {
					data = append(data, m.Index(i).Interface())
				}
			}
			base.Where(column+" in ?", data)
		} else if m.Len() == 1 {
			if m.Index(0).Interface() != "" {
				base.Where(column+" = ?", m.Index(0).Interface())
			}
		}
	case reflect.String:
		if m.String() != "" {
			base.Where(column+" = ?", value)
		}
	default:
		break
	}
}

func dataTable(base *gorm.DB, request model.RequestDataTable, dataStruct interface{}) (*model.ResultDataTable, error) {
	var (
		query   []map[string]interface{}
		search  = map[string]interface{}{}
		results model.ResultDataTable
	)
	count := int64(0)
	if len(request.Filter) > 0 {
		for _, v := range request.Filter {
			if v.Operator == "in" {
				in := strings.Split(v.Input, "^~")
				base.Where(
					v.Column+" in (?)", in,
				)
			} else {
				base.Where(
					v.Column+" "+v.Operator+" ?", v.Input,
				)
			}
		}
	}

	columns := []string{}
	x := reflect.TypeOf(dataStruct)
	for i := 0; i < x.NumField(); i++ {
		f := x.Field(i)
		tag := f.Tag.Get("datatable")
		if tag == "-" {
			continue
		}
		json := f.Tag.Get("json")
		columns = append(columns, json)
	}

	requestSearch := strings.Trim(request.Search, " ")
	if requestSearch != "" {
		groupcond := "("
		for i, col := range columns {
			if col == "no" {
				continue
			}
			groupcond += "LOWER(CAST(" + col + " AS TEXT)) like @" + col
			if i < len(columns)-1 {
				groupcond += " OR "
			}
			search[col] = "%" + strings.ToLower(requestSearch) + "%"
		}
		groupcond += ")"
		base.Where(base.Where(groupcond, search))
	}
	base.Count(&count)
	order := request.OrderBy
	if request.OrderBy == "" {
		order = columns[0]
	}
	base.Order(order + " " + request.OrderDesc).
		Limit(
			request.Length,
		).Offset(
		request.Start,
	).Where("deleted_at is null").Find(&query)
	if base.Error != nil {
		return &model.ResultDataTable{}, base.Error
	}
	var no = request.Start + 1
	for _, val := range query {
		val["no"] = no
		no++
		results.Data = append(results.Data, val)
	}
	results.Count = count
	return &results, nil
}

func convertToMap(data interface{}) (result map[string]interface{}) {
	b, _ := json.Marshal(data)
	json.Unmarshal(b, &result)
	return result
}
