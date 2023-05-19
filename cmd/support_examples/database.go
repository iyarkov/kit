package main

import (
	"github.com/iyarkov/foundation/schema"
)

var changeset = []schema.Change{
	{
		Version: "0.0.1",
		Commands: []string{
			`CREATE TABLE table_a (
				id SERIAL,
				created_at TIMESTAMP(3) WITHOUT TIME ZONE,
				name varchar(255) CONSTRAINT name_idx UNIQUE,
				description VARCHAR(255),
				b_flag bool,
				partition int,
				subpartition float,
				PRIMARY KEY (id)
			)`,
			"CREATE INDEX partition_idx ON table_a(partition, subpartition)",
		},
	},
	{
		Version: "0.0.2",
		Commands: []string{
			`CREATE TABLE table_b (
				id SERIAL,
				parent_id integer not null,
				created_at TIMESTAMP(3) WITHOUT TIME ZONE,
				name varchar(255),
				PRIMARY KEY (id),
				CONSTRAINT table_b_parent_fk
					FOREIGN KEY(parent_id) 
					REFERENCES table_a(id)
					ON DELETE CASCADE
					ON UPDATE RESTRICT
			)`,
			"CREATE UNIQUE INDEX parent_name_idx ON table_b(parent_id, name)",
		},
	},
}

var expectedSchema = schema.Schema{
	Name: "public",
	Tables: map[string]schema.Table{
		"table_a": {
			Columns: map[string]schema.Column{
				"id": {
					Type: "int4",
				},
			},
		},
	},
}
