CREATE TABLE  node (
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR(256),
  parent_id INTEGER NOT NULL
);

INSERT INTO node (name, parent_id) VALUES ("ROOT", 0);--1
INSERT INTO node (name, parent_id) VALUES ("c1",1);
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

CREATE PROCEDURE CreateChild (main INTEGER, newName VARCHAR(256))
BEGIN
    UPDATE node
    SET children = true
    WHERE id =  main;
    INSERT INTO node (name, parent_id)
    VALUES (NodeName(parent_id, newName), main);
    SELECT * FROM node WHERE parent_id = main;
end;

DROP PROCEDURE IF EXISTS CreateChild;
CALL CreateChild(1, "c2");

CREATE FUNCTION NodeName(main INTEGER, newName VARCHAR(256))RETURNS VARCHAR(256)
    DETERMINISTIC
BEGIN
    WITH nodes AS (
        SELECT COUNT(*) AS countNodes
        FROM node
        WHERE parent_id = 1 AND name LIKE CONCAT(newName,"-%") OR name = newName
    )
    SELECT IF(countNodes = 0, newName, CONCAT(newName,"-",countNodes)) INTO newName
    FROM nodes;
    RETURN newName;
end;
DROP FUNCTION NodeName;
SELECT NodeName(1,"c2") AS name;


CREATE FUNCTION FindChild(current INTEGER, childName VARCHAR(256)) RETURNS INTEGER
    DETERMINISTIC
BEGIN
    DECLARE child_id INTEGER;
    WITH children AS (
        SELECT *
        FROM node WHERE parent_id = current
    )
    SELECT id INTO child_id FROM children WHERE name = childName;
    RETURN child_id;
end;

SELECT FindChild(5, "child-5");