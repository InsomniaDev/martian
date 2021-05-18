package cassandra

type Config struct {
	Name   string `cql:"name"`
	Record string `cql:"record"`
}

// GetConfig will get the configuration value stored for the provided name
func (s *Session) GetConfig(name string) string {
	var tagsToRecords []Config
	m := map[string]interface{}{}
	query := "SELECT * FROM config WHERE name = ?"
	iterable := s.Connection.Query(query, name).Iter()
	for iterable.MapScan(m) {
		tagsToRecords = append(tagsToRecords, Config{
			Name:   m["name"].(string),
			Record: m["record"].(string),
		})
		m = map[string]interface{}{}
	}
	// There should only ever be one config per name
	return tagsToRecords[0].Record
}
