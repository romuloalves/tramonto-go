package db

// GetTests returns a list of tests
// func (d *OneSQLite) GetTests() ([]string, error) {
// 	d.mux.Lock()
// 	defer d.mux.Unlock()

// 	rows, err := d.db.Query("select id, name, description from tests;")
// 	if err != nil {
// 		return []string{}, err
// 	}

// 	for rows.Next() {
// 		var id int
// 		var name string
// 		var description string

// 		if err := rows.Scan(&id, &name, &description); err != nil {
// 			return []string{}, err
// 		}
// 	}
// }
