package local_cache

// 搭配local_cache_mgr使用, 暫時用不到了
// cache key統一設置在這以確保唯一性, 避免不小心覆蓋了別的業務的緩存
// 命名規則: Key{表名}{緩存名}
// 例如: Key{Config}{TemplateVersionAll}
// 命名: KeyConfigTemplateVersionAll

const KeyConfigTemplateVersionAll = "KeyConfigTemplateVersionAll"
