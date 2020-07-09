CREATE TABLE  node (
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR(256),
  parent_id INTEGER NOT NULL
);

INSERT INTO node (name, parent_id) VALUES ("ROOT", 0);--1
INSERT INTO node (name, parent_id) VALUES ("c1", 1);
INSERT INTO node (name, parent_id) VALUES ("c1c1",2);
INSERT INTO node (name, parent_id) VALUES ("c1c2",2);
INSERT INTO node (name,parent_id) VALUES ("c2",1);--5
INSERT INTO node (name,parent_id) VALUES ("c2c1",5);
INSERT INTO node (name,parent_id) VALUES ("c2c2",5);
INSERT INTO node (name,parent_id) VALUES ("c2c2c1",7);

-- finds parents
WITH RECURSIVE tree AS (
  SELECT id, parent_id,name, 1 AS depth
  FROM node
  WHERE id=3 

  UNION ALL

  SELECT n.id, n.parent_id,n.name, t.depth + 1
  FROM node n, tree t
  WHERE n.id = t.parent_id
)

SELECT * FROM tree;

-- find children
WITH RECURSIVE tree AS (
  SELECT id, parent_id, name, 0 AS depth
  FROM node
  WHERE id = 2

  UNION ALL

  SELECT n.id, n.parent_id, n.name, t.depth -1 
  FROM node n, tree t
  WHERE t.id = n.parent_id AND depth > -2
)

SELECT * FROM tree;

