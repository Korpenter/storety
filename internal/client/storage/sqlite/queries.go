package sqlite

const (
	// createTables is a query to create the data table.
	createTableData = `CREATE TABLE IF NOT EXISTS data (
    id TEXT NOT NULL PRIMARY KEY,
	name TEXT UNIQUE,
	type TEXT CHECK (type IN ('Card', 'Cred', 'Binary', 'Text')),
	content BLOB,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	first_synced BOOLEAN NOT NULL DEFAULT 0,
	deleted BOOLEAN NOT NULL DEFAULT 0
	);`

	// createData is a query to insert a new data record.
	createData = `INSERT OR IGNORE INTO data 
	(
		id,
		name,
		type,
		content,
		updated_at
	) VALUES (
		?,
		?,
		?,
		?,
		?
	);`

	// getDataContentByName is a query to get the content and type of a data record by its name and user ID.
	getDataContentByName = `
	SELECT content, type
	FROM data
	WHERE name = ? AND deleted = 0;`

	// getAllDataInfo is a query to get all data records' name and type for a specific user ID.
	getAllDataInfo = `
	SELECT name, type
	FROM data
	WHERE deleted = 0;`

	// deleteDataByName is a quy to delete a data record by its name and user ID.
	deleteDataByName = `
	UPDATE data
	SET name = NULL, deleted = 1, content = NULL, updated_at = ?
	WHERE name = ?`

	// getNewData is a query to get all data records for user that were created after last client sync.
	getNewData = `
	SELECT id, name, type, content, updated_at, deleted
	FROM data
	WHERE first_synced = 0`

	// getUpdatedData is a query to get all data id, hash and update time for records.
	getSyncData = `
	SELECT id, content, updated_at
	FROM data`

	// setSyncedStatus is a query to set synced_at timestamp for a data record.
	setSyncedStatus = `
	UPDATE data
	SET first_synced = 1
	WHERE id = ?;`

	// insertOrReplaceData is a query to upsert a data record.
	insertOrReplaceData = `
	INSERT OR REPLACE INTO data (id, name, type, content, updated_at, deleted, first_synced)
	VALUES (?, ?, ?, ?, ?, ?, 1);
`
	getBatch = `
	SELECT id, name, type, content, updated_at, deleted
	FROM data
	WHERE id IN (?`
)
