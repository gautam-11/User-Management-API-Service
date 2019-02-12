package modules

// RoleMap function maps user/admin to its allowed operations map
func RoleMap(role string, crud uint8) bool {
	RoleMapper := make(map[string][4]bool)
	RoleMapper["admin"] = [4]bool{true, true, true, true} // C , R , U , D
	RoleMapper["user"] = [4]bool{false, true, false, false}

	return RoleMapper[role][crud]
}
