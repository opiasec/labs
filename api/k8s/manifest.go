package k8s

import (
	networkingv1 "k8s.io/api/networking/v1"
)

func int32Ptr(i int32) *int32 { return &i }
func boolPtr(b bool) *bool    { return &b }
func pathTypePtr(p string) *networkingv1.PathType {
	pt := networkingv1.PathType(p)
	return &pt
}
