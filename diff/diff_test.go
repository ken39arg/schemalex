package diff_test

import (
	"bytes"
	"testing"

	"github.com/schemalex/schemalex/diff"
	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	type Spec struct {
		Before string
		After  string
		Expect string
	}

	specs := []Spec{
		// drop table
		{
			Before: "CREATE TABLE `hoge` ( `id` integer not null ); CREATE TABLE `fuga` ( `id` INTEGER NOT NULL );",
			After:  "CREATE TABLE `fuga` ( `id` INTEGER NOT NULL );",
			Expect: "DROP TABLE `hoge`;",
		},
		// create table
		{
			Before: "CREATE TABLE `fuga` ( `id` INTEGER NOT NULL );",
			After:  "CREATE TABLE `hoge` ( `id` INTEGER NOT NULL ) ENGINE=InnoDB; CREATE TABLE `fuga` ( `id` INTEGER NOT NULL );",
			Expect: "CREATE TABLE `hoge` (\n`id` INTEGER NOT NULL\n) ENGINE = InnoDB;",
		},
		// drop column
		{
			Before: "CREATE TABLE `fuga` ( `id` INTEGER NOT NULL, `c` VARCHAR (20) NOT NULL DEFAULT 'xxx' );",
			After:  "CREATE TABLE `fuga` ( `id` INTEGER NOT NULL );",
			Expect: "ALTER TABLE `fuga` DROP COLUMN `c`;",
		},
		// add column
		{
			Before: "CREATE TABLE `fuga` ( `id` INTEGER NOT NULL );",
			After:  "CREATE TABLE `fuga` ( `id` INTEGER NOT NULL, `c` VARCHAR (20) NOT NULL DEFAULT 'xxx');",
			Expect: "ALTER TABLE `fuga` ADD COLUMN `c` VARCHAR (20) NOT NULL DEFAULT \"xxx\";",
		},
		// change column
		{
			Before: "CREATE TABLE `fuga` ( `id` INTEGER NOT NULL );",
			After:  "CREATE TABLE `fuga` ( `id` BIGINT NOT NULL );",
			Expect: "ALTER TABLE `fuga` CHANGE COLUMN `id` `id` BIGINT NOT NULL;",
		},
		// drop primary key
		{
			Before: "CREATE TABLE `fuga` ( `id` INTEGER NOT NULL AUTO_INCREMENT, PRIMARY KEY (`id`) );",
			After:  "CREATE TABLE `fuga` ( `id` INTEGER NOT NULL AUTO_INCREMENT );",
			Expect: "ALTER TABLE `fuga` DROP INDEX PRIMARY KEY;",
		},
		// add primary key
		{
			Before: "CREATE TABLE `fuga` ( `id` INTEGER NOT NULL AUTO_INCREMENT );",
			After:  "CREATE TABLE `fuga` ( `id` INTEGER NOT NULL AUTO_INCREMENT, PRIMARY KEY (`id`) );",
			Expect: "ALTER TABLE `fuga` ADD PRIMARY KEY (`id`);",
		},
		// drop unique key
		{
			Before: "CREATE TABLE `fuga` ( `id` INTEGER NOT NULL AUTO_INCREMENT, CONSTRAINT `symbol` UNIQUE KEY `uniq_id` USING BTREE (`id`) );",
			After:  "CREATE TABLE `fuga` ( `id` INTEGER NOT NULL AUTO_INCREMENT );",
			Expect: "ALTER TABLE `fuga` DROP INDEX `uniq_id`;",
		},
		// add unique key
		{
			Before: "CREATE TABLE `fuga` ( `id` INTEGER NOT NULL AUTO_INCREMENT );",
			After:  "CREATE TABLE `fuga` ( `id` INTEGER NOT NULL AUTO_INCREMENT, CONSTRAINT `symbol` UNIQUE KEY `uniq_id` USING BTREE (`id`) );",
			Expect: "ALTER TABLE `fuga` ADD CONSTRAINT `symbol` UNIQUE INDEX `uniq_id` USING BTREE (`id`);",
		},
		// not change index
		{
			Before: "CREATE TABLE `fuga` ( `id` INTEGER NOT NULL AUTO_INCREMENT, CONSTRAINT `symbol` UNIQUE KEY `uniq_id` USING BTREE (`id`) );",
			After:  "CREATE TABLE `fuga` ( `id` INTEGER NOT NULL AUTO_INCREMENT, CONSTRAINT `symbol` UNIQUE KEY `uniq_id` USING BTREE (`id`) );",
			Expect: "",
		},
	}

	var buf bytes.Buffer
	for _, spec := range specs {
		buf.Reset()

		if !assert.NoError(t, diff.Strings(&buf, spec.Before, spec.After), "diff.String should succeed") {
			return
		}
		if !assert.Equal(t, spec.Expect, buf.String(), "result SQL should match") {
			return
		}
	}
}
