package authority

import (
	"errors"

	"gorm.io/gorm"
)

// Authority helps deal with permissions
type Authority struct {
	DB *gorm.DB
}

type Options struct {
	TablesPrefix string
	DB           *gorm.DB
}

var tablePrefix string
var auth *Authority

func New(opts Options) *Authority {
	tablePrefix = opts.TablesPrefix
	auth = &Authority{
		DB: opts.DB,
	}

	migrateTables(opts.DB)
	return auth
}

// Resolve returns the initiated instance
func Resolve() *Authority {
	return auth
}

// CreateRole stores a role in the database
// it accepts the role name. it returns an error
// in case of any
func (a *Authority) CreateRole(roleName string) error {
	var dbRole Role
	res := a.DB.Where("name = ?", roleName).First(&dbRole)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			// create
			a.DB.Create(&Role{Name: roleName})
			return nil
		}
	}

	return res.Error
}

func (a *Authority) CreatePermission(permName string) error {
	var dbPerm Permission
	res := a.DB.Where("name = ?", permName).First(dbPerm)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			a.DB.Create(&Permission{Name: permName})
		}
	}

	return res.Error
}

func (a *Authority) AssignPermissions(roleName string, permNames []string) error {
	var role Role
	rRes := a.DB.Where("name = ?", roleName).First(&role)
	if rRes.Error != nil {
		if errors.Is(rRes.Error, gorm.ErrRecordNotFound) {
			return errors.New("role record not found")
		}
	}

	var perms []Permission
	for _, permName := range permNames {
		var perm Permission
		pRes := a.DB.Where("name = ?", permName).First(&perm)
		if pRes.Error != nil {
			if errors.Is(pRes.Error, gorm.ErrRecordNotFound) {
				return errors.New("a permission record not found")
			}

		}
		perms = append(perms, perm)
	}

	for _, perm := range perms {
		// ignore if permission had assigned
		var rolePerm RolePermission
		res := a.DB.Where("role_id = ?", role.ID).Where("permission_id = ?", perm.ID).First(&rolePerm)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				cRes := a.DB.Create(&RolePermission{RoleID: role.ID, PermissionID: perm.ID})
				if cRes.Error != nil {
					return cRes.Error
				}
			}
		}
	}

	return nil
}

func (a *Authority) AssignRole(userID uint, roleName string) error {
	var role Role
	res := a.DB.Where("name = ?", roleName).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.New("missing role record")
		}
	}

	// check if the role is already assigned
	var userRole UserRole
	res = a.DB.Where("user_id = ?", userID).Where("role_id = ?", role.ID).First(&userRole)
	if res.Error == nil {
		//found a record, this role is already assigned to the same user
		return errors.New("this role is already assinged to the user")
	}

	a.DB.Create(&UserRole{UserID: userID, RoleID: role.ID})

	return nil
}

func (a *Authority) CheckRole(userID uint, roleName string) (bool, error) {
	// find the role
	var role Role
	res := a.DB.Where("name = ?", roleName).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("Role not found")
		}

	}

	// check if the role is a ssigned
	var userRole UserRole
	res = a.DB.Where("user_id = ?", userID).Where("role_id = ?", role.ID).First(&userRole)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}

	}

	return true, nil
}

func (a *Authority) CheckPermission(userID uint, permName string) (bool, error) {
	var userRoles []UserRole
	res := a.DB.Where("user_id = ?", userID).Find(&userRoles)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, res.Error
		}
	}

	var roleIDs []UserRole
	for _, role := range userRoles {
		roleIDs = append(roleIDs, role)
	}

	var perm Permission
	res = a.DB.Where("name = ?", permName).First(&perm)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("permission not found")
		}
	}

	var rolePermission RolePermission
	res = a.DB.Where("role_id IN (?)", roleIDs).Where("permission_id = ?", perm.ID).First(&rolePermission)
	if res.Error != nil {
		return false, nil
	}

	return true, nil
}

func (a *Authority) CheckRolePermission(roleName string, permName string) (bool, error) {
	// find the role
	var role Role
	res := a.DB.Where("name = ?", roleName).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("role not found")
		}

	}

	// find the permission
	var perm Permission
	res = a.DB.Where("name = ?", permName).First(&perm)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("permission not found")
		}

	}

	// find the rolePermission
	var rolePermission RolePermission
	res = a.DB.Where("role_id = ?", role.ID).Where("permission_id = ?", perm.ID).First(&rolePermission)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}

	}

	return true, nil
}

func (a *Authority) RevokeRole(userID uint, roleName string) error {
	// find the role
	var role Role
	res := a.DB.Where("name = ?", roleName).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}

	}

	a.DB.Where("user_id = ?", userID).Where("role_id = ?", role.ID).Delete(&UserRole{})

	return nil
}

func (a *Authority) RevokePermission(userID uint, permName string) error {
	// revoke the permission from all roles of the user
	// find the user roles
	var userRoles []UserRole
	res := a.DB.Where("user_id = ?", userID).Find(&userRoles)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil
		}

	}

	// find the permission
	var perm Permission
	res = a.DB.Where("name = ?", permName).First(&perm)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.New("permission not found")
		}

	}

	for _, r := range userRoles {
		// revoke the permission
		a.DB.Where("role_id = ?", r.RoleID).Where("permission_id = ?", perm.ID).Delete(RolePermission{})
	}

	return nil
}

func (a *Authority) RevokeRolePermission(roleName string, permName string) error {
	// find the role
	var role Role
	res := a.DB.Where("name = ?", roleName).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}

	}

	// find the permission
	var perm Permission
	res = a.DB.Where("name = ?", permName).First(&perm)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.New("permission not found")
		}

	}

	// revoke the permission
	a.DB.Where("role_id = ?", role.ID).Where("permission_id = ?", perm.ID).Delete(RolePermission{})

	return nil
}

func (a *Authority) GetRoles() ([]string, error) {
	var result []string
	var roles []Role
	a.DB.Find(&roles)

	for _, role := range roles {
		result = append(result, role.Name)
	}

	return result, nil
}

// GetUserRoles returns all user assigned roles
func (a *Authority) GetUserRoles(userID uint) ([]string, error) {
	var result []string
	var userRoles []UserRole
	a.DB.Where("user_id = ?", userID).Find(&userRoles)

	for _, r := range userRoles {
		var role Role
		// for every user role get the role name
		res := a.DB.Where("id = ?", r.RoleID).Find(&role)
		if res.Error == nil {
			result = append(result, role.Name)
		}
	}

	return result, nil
}

// GetPermissions returns all stored permissions
func (a *Authority) GetPermissions() ([]string, error) {
	var result []string
	var perms []Permission
	a.DB.Find(&perms)

	for _, perm := range perms {
		result = append(result, perm.Name)
	}

	return result, nil
}

// DeleteRole deletes a given role
// if the role is assigned to a user it returns an error
func (a *Authority) DeleteRole(roleName string) error {
	// find the role
	var role Role
	res := a.DB.Where("name = ?", roleName).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}

	}

	// check if the role is assigned to a user
	var userRole UserRole
	res = a.DB.Where("role_id = ?", role.ID).First(&userRole)
	if res.Error == nil {
		// role is assigned
		return errors.New("cannot delete assigned role")
	}

	// revoke the assignment of permissions before deleting the role
	a.DB.Where("role_id = ?", role.ID).Delete(RolePermission{})

	// delete the role
	a.DB.Where("name = ?", roleName).Delete(Role{})

	return nil
}

// DeletePermission deletes a given permission
// if the permission is assigned to a role it returns an error
func (a *Authority) DeletePermission(permName string) error {
	// find the permission
	var perm Permission
	res := a.DB.Where("name = ?", permName).First(&perm)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.New("permission not found")
		}

	}

	// check if the permission is assigned to a role
	var rolePermission RolePermission
	res = a.DB.Where("permission_id = ?", perm.ID).First(&rolePermission)
	if res.Error == nil {
		// role is assigned
		return errors.New("cannot delete assigned permission")
	}

	// delete the permission
	a.DB.Where("name = ?", permName).Delete(Permission{})

	return nil
}

func migrateTables(db *gorm.DB) {
	db.AutoMigrate(&Role{})
	db.AutoMigrate(&Permission{})
	db.AutoMigrate(&RolePermission{})
	db.AutoMigrate(&UserRole{})
}
