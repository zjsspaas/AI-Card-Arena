package model

var rolePriority = map[UserRole]int{
	Admin:     3,
	Developer: 2,
	Viewer:    1,
}

// 从多角色中选主角色（用于 JWT / 鉴权）
func PickMainRole(roles []Role) UserRole {
	result := Viewer
	max := 0

	for _, r := range roles {
		role := UserRole(r.Name)
		if p, ok := rolePriority[role]; ok && p > max {
			max = p
			result = role
		}
	}
	return result
}
