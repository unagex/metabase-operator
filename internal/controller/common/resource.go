package common

func GetLabels(name string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "metabase",
		"app.kubernetes.io/instance":   name,
		"app.kubernetes.io/managed-by": "metabase-operator",
	}
}
