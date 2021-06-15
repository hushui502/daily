package godis

func genExpireTask(key string) string {
	return "expire: " + key
}
