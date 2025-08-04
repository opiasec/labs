package labide

type GetLabResponse struct {
	Status        string `json:"status"`
	LastHeartbeat int64  `json:"lastHeartbeat"`
}
