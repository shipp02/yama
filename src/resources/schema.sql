CREATE TABLE users
(
    id            int          NOT NULL AUTO_INCREMENT,
    name          VARCHAR(50)  NOT NULL,
    username      VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(373) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE posts
(
    id       int NOT NULL AUTO_INCREMENT,
    owner_id int NOT NULL,
    text     TEXT,
    PRIMARY KEY (id),
    FOREIGN KEY (owner_id) REFERENCES users (id)
);

CREATE TABLE document
(
    id      INTEGER      NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name    VARCHAR(256) NOT NULL,
    content BLOB,
    type    ENUM ('pdf','docx','jpg','png','pages', 'odt', 'txt', 'html', 'gif', 'mpeg', 'mp3', 'aac', 'ai')
);
CREATE TABLE node
(
    id          INTEGER      NOT NULL PRIMARY KEY AUTO_INCREMENT,
    document_id INTEGER UNIQUE,
    name        VARCHAR(256) NOT NULL,
    children    BOOLEAN      NOT NULL DEFAULT false,
    parent_id   INTEGER      NOT NULL,
    FOREIGN KEY (document_id) REFERENCES document (id)
);
# SELECT id, username, name FROM users
# -- name: GetUserByID : one
# SELECT * FROM users
# WHERE id = $1 LIMIT 1;
#
# -- name: GetUserByUsername : one
# SELECT * FROM users
# WHERE username = $1 LIMIT 1;
#
# -- name: CreateUser :exec
# INSERT INTO users (name, username, password_hash)
# VALUES ($1, $2, $3)



