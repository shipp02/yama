
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

CREATE TABLE grp
(
    id   int NOT NULL AUTO_INCREMENT,
    name VARCHAR(200),
    anonymous BOOLEAN DEFAULT false,
    PRIMARY KEY (id)
);

CREATE TABLE usergroups
(
    group_id int,
    user_id  int,
    FOREIGN KEY (group_id) REFERENCES grp (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
)

