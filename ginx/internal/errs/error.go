package errs

import "errors"

var ErrUnauthorized = errors.New("未授权")
var ErrSessionKeyNotFound = errors.New("session 中没找到对应的 key")

// ErrNoResponse 是一个 sentinel 错误。
// 也就是说，你可以通过返回这个 ErrNoResponse 来告诉 ginx 不需要继续写响应。
// 大多数情况下，这意味着你已经写回了响应。
var ErrNoResponse = errors.New("不需要返回 response")
