CREATE TABLE  node (
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR(256),
  parent_id INTEGER NOT NULL
);

INSERT INTO node (name, parent_id) VALUES ('ROOT', 0);#--1
INSERT INTO node (name, parent_id) VALUES ('c1',1);
INSERT INTO node (name, parent_id) VALUES ('c1c1',2);
INSERT INTO node (name, parent_id) VALUES ('c1c2',2);
INSERT INTO node (name,parent_id) VALUES ('c2',1);#--5
INSERT INTO node (name,parent_id) VALUES ('c2c1',5);
INSERT INTO node (name,parent_id) VALUES ('c2c2',5);
INSERT INTO node (name,parent_id) VALUES ('c2c2c1',7);

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

CREATE PROCEDURE CreateChild (new_parent_id INTEGER, newName VARCHAR(256))
BEGIN
    UPDATE node
    SET children = true
    WHERE id =  new_parent_id;
    INSERT INTO node (name, parent_id)
    VALUES (NodeName(new_parent_id, newName), new_parent_id);
    SELECT * FROM node WHERE parent_id = new_parent_id;
end;

DROP PROCEDURE IF EXISTS CreateChild;
CALL CreateChild(2 ,'child');

CREATE FUNCTION NodeName(main INTEGER, newName VARCHAR(256))RETURNS VARCHAR(256)
    DETERMINISTIC
BEGIN
    WITH nodes AS (
        SELECT COUNT(*) AS countNodes
        FROM node
        WHERE parent_id = main AND (name REGEXP CONCAT('^', newName, '-{1}', '[:digit:]*$') OR name = newName)
    )
    SELECT IF(countNodes = 0, newName, CONCAT(newName,'-',countNodes)) INTO newName
    FROM nodes;
    RETURN newName;
end;
DROP FUNCTION NodeName;
SELECT NodeName(1,'c2') AS name;


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

CREATE PROCEDURE FindChildDetails(current INTEGER, childName VARCHAR(256))
BEGIN
    SELECT id AS 'node.id',
           name AS 'node.name' ,
           children AS 'node.children',
           parent_id AS 'node.parent_id',
           document_id AS 'node.document_id'
           FROM node WHERE FindChild(current, childName);

end;

SELECT FindChild(2, 'child');