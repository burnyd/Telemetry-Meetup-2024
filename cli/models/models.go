package models

type Target struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Timeout      int64  `json:"timeout"`
	Insecure     bool   `json:"insecure"`
	SkipVerify   bool   `json:"skip-verify"`
	BufferSize   int    `json:"buffer-size"`
	RetryTimer   int64  `json:"retry-timer"`
	LogTLSSecret bool   `json:"log-tls-secret"`
	Gzip         bool   `json:"gzip"`
	Token        string `json:"token"`
}

type Subscriptions struct {
	Subs struct {
		Sub1 struct {
			Name           string   `json:"name"`
			Paths          []string `json:"paths"`
			Mode           string   `json:"mode"`
			StreamMode     string   `json:"stream-mode"`
			Encoding       string   `json:"encoding"`
			SampleInterval int64    `json:"sample-interval"`
		} `json:"sub1"`
	} `json:"subscriptions"`
}

type Leader struct {
	Name                  string `json:"name"`
	NumberOfLockedTargets int    `json:"number-of-locked-targets"`
	Leader                string `json:"leader"`
	Members               []struct {
		Name                string   `json:"name"`
		APIEndpoint         string   `json:"api-endpoint"`
		NumberOfLockedNodes int      `json:"number-of-locked-nodes"`
		LockedTargets       []string `json:"locked-targets"`
		IsLeader            bool     `json:"is-leader,omitempty"`
	} `json:"members"`
}

type NewTarget struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Insecure     bool   `json:"insecure"`
	Skipverify   bool   `json:"skip-verify"`
	Buffersize   int    `json:"buffer-size"`
	RetryTimer   int    `json:"retry-timer"`
	Logtlssecret bool   `json:"log-tls-secret"`
	Gzip         bool   `json:"gzip"`
	Token        string `json:"token"`
	Timeout      int    `json:"timeout"`
}
