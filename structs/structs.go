package structs

//this is the struct that receives the dns zone lookup
type ZoneresultMap struct {
	Zones []struct {
		ID            string   `json:"id"`
		Name          string   `json:"name"`
		TTL           int      `json:"ttl"`
		Registrar     string   `json:"registrar"`
		LegacyDNSHost string   `json:"legacy_dns_host"`
		LegacyNs      []string `json:"legacy_ns"`
		Ns            []string `json:"ns"`
		Created       string   `json:"created"`
		Verified      string   `json:"verified"`
		Modified      string   `json:"modified"`
		Project       string   `json:"project"`
		Owner         string   `json:"owner"`
		Permission    string   `json:"permission"`
		ZoneType      struct {
			ID          string      `json:"id"`
			Name        string      `json:"name"`
			Description string      `json:"description"`
			Prices      interface{} `json:"prices"`
		} `json:"zone_type"`
		Status          string `json:"status"`
		Paused          bool   `json:"paused"`
		IsSecondaryDNS  bool   `json:"is_secondary_dns"`
		TxtVerification struct {
			Name  string `json:"name"`
			Token string `json:"token"`
		} `json:"txt_verification"`
		RecordsCount int `json:"records_count"`
	} `json:"zones"`
	Meta struct {
		Pagination struct {
			Page         int `json:"page"`
			PerPage      int `json:"per_page"`
			PreviousPage int `json:"previous_page"`
			NextPage     int `json:"next_page"`
			LastPage     int `json:"last_page"`
			TotalEntries int `json:"total_entries"`
		} `json:"pagination"`
	} `json:"meta"`
}

//struct to receive a dns record lookup result
type RecordresultMap struct {
	Records []struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		Name     string `json:"name"`
		Value    string `json:"value"`
		TTL      int    `json:"ttl,omitempty"`
		ZoneID   string `json:"zone_id"`
		Created  string `json:"created"`
		Modified string `json:"modified"`
	} `json:"records"`
}


