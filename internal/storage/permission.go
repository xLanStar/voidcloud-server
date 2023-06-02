package storage

type Permission struct {
	Index uint8
	Name  string
	Can   map[string]bool
}

func (permission *Permission) ToString() string {
	return permission.Name
}

var (
	NO_PERMISSION *Permission = &Permission{
		0,
		"無權限",
		map[string]bool{
			"COPY":      false,
			"DELETE":    false,
			"GET":       false,
			"HEAD":      false,
			"LOCK":      false,
			"MKCOL":     false,
			"MOVE":      false,
			"OPTIONS":   false,
			"POST":      false,
			"PROPFIND":  false,
			"PROPPATCH": false,
			"PUT":       false,
			"TRACE":     false,
			"UNLOCK":    false,
		}}
	READ_ONLY *Permission = &Permission{
		1,
		"只能讀取",
		map[string]bool{
			"COPY":      true,
			"DELETE":    false,
			"GET":       true,
			"HEAD":      true,
			"LOCK":      true,
			"MKCOL":     false,
			"MOVE":      false,
			"OPTIONS":   true,
			"POST":      false,
			"PROPFIND":  true,
			"PROPPATCH": false,
			"PUT":       false,
			"TRACE":     true,
			"UNLOCK":    false,
		}}
	ALL *Permission = &Permission{
		7,
		"全部權限",
		map[string]bool{
			"COPY":      true,
			"DELETE":    true,
			"GET":       true,
			"HEAD":      true,
			"LOCK":      true,
			"MKCOL":     true,
			"MOVE":      true,
			"OPTIONS":   true,
			"POST":      true,
			"PROPFIND":  true,
			"PROPPATCH": true,
			"PUT":       true,
			"TRACE":     true,
			"UNLOCK":    true,
		}}
)

var PermissionMap = []*Permission{
	NO_PERMISSION,
	READ_ONLY,
	READ_ONLY,
	READ_ONLY,
	READ_ONLY,
	READ_ONLY,
	READ_ONLY,
	ALL,
}

func ParsePermission(index uint8) *Permission {
	if index > 7 {
		return ALL
	}
	return PermissionMap[index]
}
