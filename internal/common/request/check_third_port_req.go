package request

type CheckThirdPortRequest struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}
