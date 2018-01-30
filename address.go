package icheck

type Address struct {
	ID       uint64
	Address  string `json:"address"`
	City     int64  `json:"city"`
	District int64  `json:"district"`
	Email    string `json:"email"`
}

type AddressListResp struct {
	Data []Address
}

type AddressResp struct {
	Data Address
}

type AddressBody struct {
	Address  string `json:"address"`
	City     int64  `json:"city"`
	District int64  `json:"district"`
	Email    string `json:"email"`
}
