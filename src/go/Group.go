package main


type Group struct{
	ID int64      `db:"id"`
	Name string   `db:"name"`
}

// GroupSchema
const GroupSchema = `
CREATE TABLE grp (
	id int NOT NULL AUTO_INCREMENT,
	name VARCHAR(200), 
	PRIMARY KEY(id) 
	);
`

// GroupUserSchema provides many-many scheme for groups
const GroupUserSchema = `
CREATE TABLE usergroups (
	group_id int,
	user_id int,
	FOREIGN KEY (group_id) REFERENCES grp(id),
	FOREIGN KEY (user_id) REFERENCES users(id)
)
`