package database

type Device struct {
	DeviceId                string `json:"deviceId"`
	UniqueHash              string `json:"uniqueHash"`
	Name                    string `json:"name,omitempty"`
	Label                   string `json:"label,omitempty"`
	Type                    string `json:"type,omitempty"`
	Integration             string `json:"integration,omitempty"`
	EnergyEfficiencyMinutes int    `json:"energyEfficiencyMinutes,omitempty"`
}

type UniqueHashGraphs struct {
	Graphs []DeviceGraph `json:"graphs"`
}

type DeviceGraph struct {
	FromUniqueHash string `json:"fromUniqueHash"`
	FromStatus     string `json:"fromStatus"`
	ToUniqueHash   string `json:"toUniqueHash"`
	ToStatus       string `json:"toStatus"`
	Weight         int    `json:"weight"`
	Automated      bool   `json:"automated"`
	TimeAutomated  bool   `json:"timeAutomated"`
}

type TimeTable struct {
	Times []TimeBlocks `json:"timeBlocks"`
}

type TimeBlocks struct {
	UniqueHash string `json:"uniqueHash"`
	Weight     int    `json:"weight"`
	Automated  bool   `json:"automated"`
	TimeKey    string `json:"timeKey"`
}

// // queryDeviceGraphResults is a function to extract repetitive query syntax
// func (d *Database) queryDeviceGraphResults(query string, params ...interface{}) (graphResults []DeviceGraph) {
// 	conn, err := pgx.Connect(context.Background(), d.DatabaseConnectionString)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
// 		os.Exit(1)
// 	}
// 	defer conn.Close(context.Background())

// 	rows, err := conn.Query(context.Background(), query, params...)
// 	if err != nil {
// 		log.Debugf("Unable to find query for results: \n%s\n\n", query, err)
// 		return nil
// 	}

// 	for rows.Next() {
// 		var graph DeviceGraph
// 		if err := rows.Scan(
// 			&graph.FromUniqueHash,
// 			&graph.FromStatus,
// 			&graph.ToUniqueHash,
// 			&graph.ToStatus,
// 			&graph.Weight,
// 			&graph.Automated,
// 			&graph.TimeAutomated,
// 		); err != nil {
// 			log.Fatal(err)
// 		}
// 		graphResults = append(graphResults, graph)
// 	}
// 	return
// }

// // queryDeviceResults is a function to extract repetitive query syntax
// func (d *Database) queryDeviceResults(query string, params ...interface{}) (deviceResults []Device) {
// 	conn, err := pgx.Connect(context.Background(), d.DatabaseConnectionString)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
// 		os.Exit(1)
// 	}
// 	defer conn.Close(context.Background())

// 	rows, err := conn.Query(context.Background(), query, params...)
// 	if err != nil {
// 		log.Debugf("Unable to find query for results: \n%s\n\n", query, err)
// 		return nil
// 	}

// 	for rows.Next() {
// 		var device Device
// 		if err := rows.Scan(
// 			&device.DeviceId,
// 			&device.UniqueHash,
// 			&device.Name,
// 			&device.Label,
// 			&device.Type,
// 			&device.Integration,
// 			&device.EnergyEfficiencyMinutes,
// 		); err != nil {
// 			log.Fatal(err)
// 		}
// 		deviceResults = append(deviceResults, device)
// 	}
// 	return
// }

// func (d *Database) transactionQuery(query string, params ...interface{}) error {
// 	conn, err := pgx.Connect(context.Background(), d.DatabaseConnectionString)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
// 		os.Exit(1)
// 	}
// 	defer conn.Close(context.Background())

// 	tx, err := conn.Begin(context.Background())
// 	if err != nil {
// 		log.Printf("Failed to start a transaction process")
// 		return err
// 	}
// 	defer tx.Rollback(context.Background())

// 	_, err = tx.Exec(context.Background(), query, params...)

// 	if err != nil {
// 		// Ignore when we are trying to insert a device that already exists, that's why we have a unique constraint
// 		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"device_uniquehash_key\"") {
// 			return nil
// 		}
// 		log.Println("Unable to run transaction: ", query, err)
// 		log.Println(params...)
// 		return err
// 	}

// 	err = tx.Commit(context.Background())
// 	if err != nil {
// 		log.Printf("Failed to commit")
// 		return err
// 	}
// 	return nil
// }
