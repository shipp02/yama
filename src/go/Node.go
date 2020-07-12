package main

const NodeSchema = `
CREATE TABLE node (
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(256) NOT NULL,
  children BOOLEAN NOT NULL,
  parent_id INTEGER NOT NULL
);
`

//type Node struct {
//	ID       sql.NullInt64
//	Name     sql.NullString
//	Children sql.NullBool
//	ParentID sql.NullInt64
//}

//func (n *Node) ToInterface(l int) (iface []interface{}) {
//	iface = make([]interface{}, l)
//	iface[0] = &n.ID
//	iface[1] = &n.Name
//	iface[2] = &n.Children
//	iface[3] = &n.ParentID
//	return iface
//}
//func (n *Node) GetNode(db *sqlx.DB) (*Node, error) {
//	var err error
//	query := "SELECT * FROM node WHERE id = %d"
//	if n.ID.Int64 == 0 {
//		err = errors.New("not enough data")
//	} else {
//		query = fmt.Sprintf(query, n.ID.Int64)
//		row, err := db.Query(query)
//		if err != nil {
//			log.Println(err.Error())
//			return n, err
//		}
//		l, _ := row.Columns()
//		for row.Next() {
//			if err := row.Scan(n.ToInterface(len(l))...); err != nil {
//				log.Fatal(err)
//			}
//		}
//		log.Println("func GetNode", n)
//
//	}
//	return n, err
//}
//
//func (n *Node) CreateNode(db *sqlx.DB) (err error) {
//	exec := "INSERT INTO node (name, children, parent_id) VALUES(\"%s\", %t, %d)"
//	if n.Name.String == "" || n.ParentID.Int64 == 0 {
//		err = errors.New("create node not enough data")
//	} else {
//		exec = fmt.Sprintf(exec, n.Name.String, n.Children.Bool, n.ParentID.Int64)
//		db.MustExec(exec)
//	}
//	return err
//}
