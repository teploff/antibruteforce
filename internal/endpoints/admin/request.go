package admin

type ResetBucketByLoginRequest struct {
	Login string `json:"login"`
}

type ResetBucketByPasswordRequest struct {
	Password string `json:"password"`
}

type ResetBucketByIPRequest struct {
	IP string `json:"ip"`
}

type SubnetRequest struct {
	IPWithMask string `json:"ip_with_mask"`
}
