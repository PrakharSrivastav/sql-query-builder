package querybuilder

import (
	"testing"
)

func TestPgsqlQB_Create(t *testing.T) {
	type fields struct {
		table     string
		idPrimary bool
	}

	type args struct {
		param map[string]interface{}
	}

	// setup test cases
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{

		{ // Test without table name
			name:    "test-no-tablename",
			fields:  fields{idPrimary: true},
			wantErr: true,
		},
		{ // Test with primary key
			name:   "table with PKey",
			fields: fields{table: "xyz", idPrimary: true},
			want:   "CREATE TABLE xyz (id UUID primary key);",
		},
		{ // Test without primary key
			name:   "create table without PKey",
			fields: fields{table: "xyz", idPrimary: false},
			want:   "CREATE TABLE xyz (id UUID);",
		},
		{ // Test with one extra column
			name:   "table Pkey with columns",
			fields: fields{table: "xyz", idPrimary: true},
			want:   "CREATE TABLE xyz (col1 text ,id UUID primary key);",
			args:   args{param: map[string]interface{}{"col1": "blah"}},
		},
		{ // Test with some columns
			name:   "table Pkey with columns",
			fields: fields{table: "xyz", idPrimary: true},
			want:   "CREATE TABLE xyz (amount float ,col1 text ,desc text ,tid varchar(200)  ,id UUID primary key);",
			args: args{param: map[string]interface{}{
				"col1":   "blah",
				"amount": "float",
				"desc":   "",
				"tid":    "keyword",
			}},
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := &PgsqlQB{
				Table:     tt.fields.table,
				IDPrimary: tt.fields.idPrimary,
			}
			got, err := qb.Create(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("\rtest(%s):\rerror\t%v\rwantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("\rtest(%s):\ractual\t\t%v\rexpected\t%v", tt.name, got, tt.want)
			}
		})
	}
}

func TestPgsqlQB_Get(t *testing.T) {
	type fields struct {
		table     string
		idPrimary bool
	}
	type args struct {
		columns []string
		where   map[string]interface{}
		limit   int
		offset  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{ // no table name should throw err
			name: "complete",
			args: args{
				columns: []string{"column1", "column2"},

				where: map[string]interface{}{
					"column1": "value1",
					"column2": 10,
				},
				limit:  1,
				offset: 10,
			},
			wantErr: true,
		},
		{ // no column should throw error
			name:   "complete",
			fields: fields{table: "xyz"},
			args: args{
				where: map[string]interface{}{
					"column1": "value1",
					"column2": 10,
				},
				limit:  1,
				offset: 10,
			},
			wantErr: true,
		},
		{ // Complete select clause with limit and offset
			name: "complete",
			fields: fields{
				table: "xyz",
			},
			args: args{
				columns: []string{"column1", "column2"},
				where: map[string]interface{}{
					"column1": "value1",
					"column2": 10,
				},
				limit:  1,
				offset: 10,
			},
			want:    "SELECT column1,column2 FROM xyz WHERE column1 = 'value1' AND column2 = 10 LIMIT 1  OFFSET 10 ;",
			wantErr: false,
		},
		{ // select clause without offset and limit
			name: "without limit and offset",
			fields: fields{
				table: "xyz",
			},
			args: args{
				columns: []string{"column1", "column2"},
				where: map[string]interface{}{
					"column1": 1.543,
					"column2": "value2",
				},
			},
			want:    "SELECT column1,column2 FROM xyz WHERE column1 = 1.543 AND column2 = 'value2';",
			wantErr: false,
		},
		{
			name: "without where",
			fields: fields{
				table: "xyz",
			},
			args: args{
				columns: []string{"column1", "column2"},
			},
			want:    "SELECT column1,column2 FROM xyz;",
			wantErr: false,
		},
		{
			name: "select *",
			fields: fields{
				table: "xyz",
			},
			args: args{
				columns: []string{"*"},
			},
			want:    "SELECT * FROM xyz;",
			wantErr: false,
		},
		{
			name: "select * with limit",
			fields: fields{
				table: "xyz",
			},
			args: args{
				columns: []string{"*"},
				limit:   1,
			},
			want:    "SELECT * FROM xyz LIMIT 1 ;",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &PgsqlQB{
				Table:     tt.fields.table,
				IDPrimary: tt.fields.idPrimary,
			}
			got, err := b.Get(tt.args.columns, tt.args.where, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("\rtest(%s):\rerror\t%v\rwantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("\rtest(%s):\ractual\t\t%v\rexpected\t%v", tt.name, got, tt.want)
			}
		})
	}
}

func TestPgsqlQB_Update(t *testing.T) {
	type fields struct {
		table     string
		idPrimary bool
	}
	type args struct {
		columns map[string]interface{}
		where   map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test-no-table-name",
			args:    args{columns: map[string]interface{}{"col1": 5.42, "col2": "val2"}},
			wantErr: true,
		},
		{
			name:    "test-no-columns",
			fields:  fields{table: "xyz"},
			args:    args{},
			wantErr: true,
		},
		{
			name:   "simple update",
			fields: fields{table: "xyz"},
			args:   args{columns: map[string]interface{}{"col1": 5.42, "col2": "val2"}},
			want:   "UPDATE xyz SET col1=5.42,col2='val2' ;",
		},
		{
			name:   "update with where clause",
			fields: fields{table: "xyz"},
			args: args{
				columns: map[string]interface{}{"col1": 5.42, "col2": "val2"},
				where:   map[string]interface{}{"col3": 123, "col4": "aaa"},
			},
			want: "UPDATE xyz SET col1=5.42,col2='val2' WHERE col3=123 AND col4='aaa' ;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &PgsqlQB{
				Table:     tt.fields.table,
				IDPrimary: tt.fields.idPrimary,
			}
			got, err := b.Update(tt.args.columns, tt.args.where)
			if (err != nil) != tt.wantErr {
				t.Errorf("\rtest(%s):\rerror\t%v\rwantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("\rtest(%s):\ractual\t\t%v\rexpected\t%v", tt.name, got, tt.want)
			}
		})
	}
}

func TestPgsqlQB_Insert(t *testing.T) {
	type fields struct {
		table     string
		idPrimary bool
	}
	type args struct {
		columns []string
		data    []map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test - no table name",
			args: args{
				columns: []string{"b", "a"},
				data: []map[string]interface{}{
					map[string]interface{}{
						"a": 123,
						"b": "sss",
					},
					map[string]interface{}{
						"a": 234,
						"b": "aaa",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Test - no columns",
			fields: fields{
				table: "XYZ",
			},
			args: args{
				data: []map[string]interface{}{
					map[string]interface{}{
						"a": 123,
						"b": "sss",
					},
					map[string]interface{}{
						"a": 234,
						"b": "aaa",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Test - no data",
			fields: fields{
				table: "XYZ",
			},
			args: args{
				columns: []string{"b", "a"},
			},
			wantErr: true,
		},
		{ // OK
			name: "123",
			fields: fields{
				table: "XYZ",
			},
			args: args{
				columns: []string{"b", "a"},
				data: []map[string]interface{}{
					map[string]interface{}{
						"a": 123,
						"b": "sss",
					},
					map[string]interface{}{
						"a": 234,
						"b": "aaa",
					},
				},
			},
			want: "INSERT INTO XYZ (a,b) VALUES (123,'sss'),(234,'aaa');",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &PgsqlQB{
				Table:     tt.fields.table,
				IDPrimary: tt.fields.idPrimary,
			}
			got, err := b.Insert(tt.args.columns, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("\rtest(%s):\rerror\t%v\rwantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("\rtest(%s):\ractual\t\t%v\rexpected\t%v", tt.name, got, tt.want)
			}
		})
	}
}
