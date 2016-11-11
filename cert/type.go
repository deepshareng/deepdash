package cert

// for iOS 9.0 and later
type UniversalLinkJson struct {
	Applinks AppLinks `json:"applinks"`
}

type AppLinks struct {
	Apps    []string `json:"apps"`
	Details []Detail `json:"details"`
}

type Detail struct {
	AppID string   `json:"appID"`
	Paths []string `json:"paths"`
}

// for android 6.0 and later
type AppLinkJson struct {
	AppLink []AppLink `json:"applink"`
}

type AppLink struct {
	Relation []string `json:"relation"`
	Target   Target   `json:"target"`
}

type Target struct {
	Namespace   string   `json:"namespace"`
	PackageName string   `json:"package_name"`
	SHA256      []string `json:"sha256_cert_fingerprints"`
}
