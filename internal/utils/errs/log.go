package errs

import klog "github.com/go-kratos/kratos/v2/log"

const msgKey = "internal_error"

var log *klog.Helper

func init() {
	log = klog.NewHelper(klog.GetLogger(), klog.WithMessageKey(msgKey))
}
