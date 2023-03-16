package postgres

const (
	// createUser is a query to insert a new user record.
	createUser = `
	INSERT INTO users (
		id,
		username,
		password
	) VALUES (
		$1,
		$2,
		$3
	) 	
	ON CONFLICT DO NOTHING
	RETURNING id`

	// getUserByName is a query to get a user record by its username.
	getUserByName = `SELECT id, password FROM users WHERE username = $1`

	// createNewSession is a query to insert a new session record.
	createNewSession = `
	INSERT INTO sessions (
		id,
		user_id,
		auth_token,
		refresh_token
	) VALUES (
		$1,
		$2,
		$3,
		$4
	)
	ON CONFLICT DO NOTHING`

	// getUserBySession is a query to get a user ID by its session ID and refresh token.
	getUserBySession = `
	SELECT user_id
	FROM sessions
	WHERE id=$1 AND refresh_token=$2`

	// deleteOldSession is a query to delete a session record by its ID and refresh token.
	deleteOldSession = `
	DELETE FROM sessions
	WHERE id=$1 AND refresh_token=$2`

	// createData is a query to insert a new data record.
	createData = `INSERT INTO data (
    id,
    user_id,
    name,
    type,
    content,
	version
	)
	SELECT
		$1,
		$2,
		CASE 
			WHEN EXISTS (
				SELECT 1 
				FROM data 
				WHERE user_id = $2 AND name = $3
			)
				THEN (
					SELECT 
						CONCAT($3, '_', COALESCE(MAX(SUBSTRING(name FROM '.*_(\d+)')::INTEGER), 0) + 1)
					FROM data 
					WHERE user_id = $2 AND name LIKE $3 || '_%'
				)
			ELSE $3
		END,
		$4,
		$5,
		(SELECT data_version FROM users WHERE id = $2)
`

	// getDataContentByName is a query to get the content and type of a data record by its name and user ID.
	getDataContentByName = `
    SELECT  content, type
    FROM data
    WHERE name = $1 AND user_id = $2`

	// getAllDataInfo is a query to get all data records' name and type for a specific user ID.
	getAllDataInfo = `
    SELECT  name, type
    FROM data
    WHERE user_id = $1`

	// updateDataVersion is a query to increment the data_version field of a user record by 1.
	updateDataVersion = `
	UPDATE users
	SET data_version = data_version + 1
	WHERE id = $1`

	// deleteDataByName is a query to delete a data record by its name and user ID.
	deleteDataByName = `
	UPDATE data SET deleted = true, content = null, version = version + 1,
	WHERE name = $1 AND user_id = $2`
)
