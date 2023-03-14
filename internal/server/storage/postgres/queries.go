package postgres

const (
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

	getUserByName = `SELECT id, password FROM users WHERE username = $1`

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

	getUserBySession = `
	SELECT user_id
	FROM sessions
	WHERE id=$1 AND refresh_token=$2`

	deleteOldSession = `
	DELETE FROM sessions
	WHERE id=$1 AND refresh_token=$2`

	createData = `INSERT INTO data (
		id,
		user_id,
		name,
		type,
		content
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
	)`

	getDataContentByName = `
    SELECT  content
    FROM data
    WHERE name = $1 AND user_id = $2`

	getAllDataInfo = `
    SELECT  name, type
    FROM data
    WHERE user_id = $1`

	updateDataVersion = `
	UPDATE users
	SET data_version = data_version + 1
	WHERE id = $1`

	deleteDataByName = `
	DELETE FROM data
	WHERE name = $1 AND user_id = $2`
)
