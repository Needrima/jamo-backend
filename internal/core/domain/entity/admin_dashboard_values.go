package entity

type DashboardValues struct {
	TotalOrders     int     `json:"totalOrders"`
	PendingOrders   int     `json:"pendingOrders"`
	CompletedOrders int     `json:"completedOrders"`
	TotalRevenue    float64 `json:"totalRevenue"`
	NetProfit       float64 `json:"netProfit"`
}
