// Copyright (c) 2017 Femtowiki authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"github.com/s-gv/femtowiki/models/db"
	"log"
	"time"
)

const ModelVersion = 1

func Migration1() {
	db.Exec(`CREATE TABLE configs(name VARCHAR(250), val TEXT);`)
	db.Exec(`CREATE UNIQUE INDEX configs_key_index on configs(name);`)

	db.Exec(`CREATE TABLE users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(32) NOT NULL,
		passwdhash VARCHAR(250) NOT NULL,
		email VARCHAR(250) DEFAULT '',
		reset_token VARCHAR(250) DEFAULT '',
		reset_token_date INTEGER DEFAULT 0,
		is_banned INTEGER DEFAULT 0,
		is_superuser INTEGER DEFAULT 0,
		created_date INTEGER DEFAULT 0,
		updated_date INTEGER DEFAULT 0
	);`)
	db.Exec(`CREATE UNIQUE INDEX users_username_index ON users(username);`)
	db.Exec(`CREATE INDEX users_email_index ON users(email);`)
	db.Exec(`CREATE INDEX users_reset_token_index ON users(reset_token);`)
	db.Exec(`CREATE INDEX users_created_index ON users(created_date);`)

	db.Exec(`CREATE TABLE groups(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(250) DEFAULT '',
		created_date INTEGER DEFAULT 0,
		updated_date INTEGER DEFAULT 0
	);`)
	db.Exec(`CREATE UNIQUE INDEX groups_name_index ON groups(name);`)

	db.Exec(`CREATE TABLE groupmembers(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
		groupid INTEGER REFERENCES groups(id) ON DELETE CASCADE,
		created_date INTEGER DEFAULT 0
	);`)
	db.Exec(`CREATE INDEX groupmembers_userid_index ON groupmembers(userid);`)
	db.Exec(`CREATE INDEX groupmembers_groupid_index ON groupmembers(groupid);`)

	db.Exec(`CREATE TABLE sessions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sessionid VARCHAR(250) DEFAULT '',
		csrftoken VARCHAR(250) DEFAULT '',
		userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
		msg VARCHAR(250) DEFAULT '',
		created_date INTEGER DEFAULT 0,
		updated_date INTEGER DEFAULT 0
	);`)
	db.Exec(`CREATE INDEX sessions_sessionid_index ON sessions(sessionid);`)
	db.Exec(`CREATE INDEX sessions_userid_index ON sessions(userid);`)

	db.Exec(`CREATE TABLE pages(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title VARCHAR(250) DEFAULT '',
		content TEXT DEFAULT '',
		discussion TEXT DEFAULT '',
		is_file INTEGER DEFAULT 0,
		editgroupid INTEGER REFERENCES groups(id) ON DELETE SET NULL,
		readgroupid INTEGER REFERENCES groups(id) ON DELETE SET NULL,
		created_date INTEGER DEFAULT 0,
		updated_date INTEGER DEFAULT 0
	);`)
	db.Exec(`CREATE INDEX pages_title_index on pages(title);`)
	db.Exec(`CREATE INDEX pages_editgroupid_index ON pages(editgroupid);`)
	db.Exec(`CREATE INDEX pages_readgroupid_index ON pages(readgroupid);`)

	db.Exec(`CREATE TABLE uploads(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(250) DEFAULT '',
		location VARCHAR(250) DEFAULT '',
		editgroupid INTEGER REFERENCES groups(id) ON DELETE SET NULL,
		readgroupid INTEGER REFERENCES groups(id) ON DELETE SET NULL,
		created_date INTEGER DEFAULT 0,
		updated_date INTEGER DEFAULT 0
	);`)
	db.Exec(`CREATE INDEX uploads_name_index ON uploads(name);`)
	db.Exec(`CREATE INDEX uploads_editgroupid_index ON uploads(editgroupid);`)
	db.Exec(`CREATE INDEX uploads_readgroupid_index ON uploads(readgroupid);`)

	if db.DriverName == "sqlite3" {
		db.Exec(`CREATE VIRTUAL TABLE pages_search USING fts5(title, content);`)
		db.Exec(`CREATE TRIGGER after_page_insert AFTER INSERT ON pages FOR EACH ROW BEGIN
			INSERT INTO pages_search(rowid, title, content) VALUES(new.id, new.title, new.content);
		END;`)
		db.Exec(`CREATE TRIGGER after_page_update UPDATE OF content ON pages FOR EACH ROW BEGIN
			UPDATE pages_search SET content=new.content WHERE rowid=old.id;
		END;`)
		db.Exec(`CREATE TRIGGER after_page_delete AFTER DELETE ON pages FOR EACH ROW BEGIN
			DELETE FROM pages_search WHERE rowid=old.id;
		END;`)
	} else if db.DriverName == "postgres" {
		db.Exec(`ALTER TABLE pages ADD COLUMN vectors tsvector;`)
		db.Exec(`CREATE INDEX pages_search_index ON pages USING GIN(vectors);`)
		db.Exec(`CREATE TRIGGER trigger_pages_vectors BEFORE INSERT OR UPDATE ON pages
			FOR EACH ROW EXECUTE PROCEDURE tsvector_update_trigger(vectors, 'pg_catalog.english', content);`)
	} else {
		log.Fatalf("[ERROR] DB Driver %s not supported\n", db.DriverName)
	}

}

func IsMigrationNeeded() bool {
	return db.Version() != ModelVersion
}

func Migrate() {
	dbver := db.Version()
	if dbver == ModelVersion {
		log.Panicf("[ERROR] DB migration not needed. DB up-to-date.\n")
	} else if dbver > ModelVersion {
		log.Panicf("[ERROR] DB version (%d) is greater than binary version (%d). Use newer binary.\n", dbver, ModelVersion)
	}
	for dbver < ModelVersion {
		if dbver == 0 {
			log.Printf("[INFO] Migrating to version 1...")
			Migration1()

			WriteConfig(Version, "1")
			WriteConfig(ConfigJSON, DefaultConfigJSON)
			WriteConfig(HeaderLinks, DefaultHeaderLinks)
			WriteConfig(FooterLinks, DefaultFooterLinks)
			WriteConfig(NavSections, DefaultNavSections)
			WriteConfig(IllegalNames, DefaultIllegalNames)
			WriteConfig(PageMasterGroup, DefaultPageMasterGroup)
			db.Exec(`INSERT INTO pages(title, content, created_date, updated_date) VALUES(?, ?, ?, ?);`, IndexPage, "# Home Page", time.Now().Unix(), time.Now().Unix())
		}
		dbver = db.Version()
	}
}