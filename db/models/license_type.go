package models

const (
	LicenseTypeNonInclusive     = "non-inclusive"
	LicenseTypeInclusive        = "inclusive"
	LicenseTypeTransferOfRights = "tor"
)

func LicenseTypesList() []string {
	return []string{LicenseTypeNonInclusive, LicenseTypeInclusive, LicenseTypeTransferOfRights}
}

type LicenseType struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (l *LicenseType) IsNonInclusive() bool {
	return l.Type == LicenseTypeNonInclusive
}

func (l *LicenseType) IsInclusive() bool {
	return l.Type == LicenseTypeInclusive
}

func (l *LicenseType) IsTransferOfRights() bool {
	return l.Type == LicenseTypeTransferOfRights
}
