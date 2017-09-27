package model

type DashboardInfo struct {
	User    *User     `json:"user"`
	Devices []*Device `json:"devices"`
}
