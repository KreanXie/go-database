package table

import (
	bufferpool "go-database/src/buffer_pool"
)

type Table struct {
	name   string
	schema *Schema
	pages  []int
	buffer *bufferpool.Manager
}

func NewTable(name string, schema *Schema, buffer *bufferpool.Manager) *Table {
	return &Table{
		name:   name,
		schema: schema,
		pages:  []int{},
		buffer: buffer,
	}
}

func (t *Table) Name() string {
	return t.name
}

func (t *Table) Schema() *Schema {
	return t.schema
}

func (t *Table) Pages() []int {
	return t.pages
}

func (t *Table) Insert(row Row) {

}
