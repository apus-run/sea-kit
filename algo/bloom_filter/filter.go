package bloom_filter

type Filter interface {
	Add(key []byte)                // 添加 key 到过滤器
	Exist(bitmap, key []byte) bool // 是否存在 key
	Hash() []byte                  // 生成过滤器对应的 bitmap
	Reset()                        // 重置过滤器
	KeyLen() int                   // 存在多少个 key
}
