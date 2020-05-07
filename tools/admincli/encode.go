package main

import (
	"encoding/json"
	"github.com/teploff/antibruteforce/endpoints/admin"
)

func encodeResetBucketByLoginRequest(login string) ([]byte, error) {
	r := admin.ResetBucketByLoginRequest{Login: login}
	marshalReq, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return marshalReq, nil
}

func encodeResetBucketByPasswordRequest(password string) ([]byte, error) {
	r := admin.ResetBucketByPasswordRequest{Password: password}
	marshalReq, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return marshalReq, nil
}

func encodeResetBucketByIPRequest(ip string) ([]byte, error) {
	r := admin.ResetBucketByIPRequest{IP: ip}
	marshalReq, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return marshalReq, nil
}

func encodeSubnetRequest(subnet string) ([]byte, error) {
	r := admin.SubnetRequest{IPWithMask: subnet}
	marshalReq, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return marshalReq, nil
}
