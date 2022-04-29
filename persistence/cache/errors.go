package cache

import "errors"

var ErrorKeyDoesNotExist = errors.New("key does not exist in cache")
var ErrorCacheServerPluginDoesNotExist = errors.New("required cache server plugin does not exist")
