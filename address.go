package icheck

type Address struct {
	ID       uint64
	Address  string
	City     int64
	District int64
	Email    string
}

type AddressListResp struct {
	Data []Address
}

type AddressResp struct {
	Data Address
}

type AddressBody struct {
	Address  string
	City     int64
	District int64
	Email    string
}
