package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Define a Permissions slice, which we will use to hold the permission codes (like
// "movies:read" and "movies:write") for a single user.
type Permissions []string

// And a helper method to check whether the permissions slice contains a specific
// Permission code
func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}

// Define the PermissionModel type

type PermissionModel struct {
	DB *sql.DB
}

// The GetAllForUser() method returns all permission codes for a specific user in a
// Permissions slice. The code in this method should feel very familiar --- it uses the
// standard pattern that we've already seen before for retrieving multiple data rows in
// an SQL query
func (m PermissionModel) GetAllForUser(userId int) (Permissions, error) {
	query := `
	SELECT permissions.code
	FROM permissions
	INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
	INNER JOIN users ON users_permissions.user_id = users.id
	WHERE users.id = ?`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, userId)
	if err != nil {
		fmt.Println("shfvsdhkfbskdjf")
		return nil, err
	}
	defer rows.Close()
	var permissions Permissions
	for rows.Next() {
		permission := ""
		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return permissions, nil
}

// Add the provided permission codes for a specific user. Notice that we're using a
// variadic parameter for the codes so that we can assign multiple permissions in a
// single call.
func (m PermissionModel) AddForUser(userID int, codes string) error {
	query := `
	INSERT INTO users_permissions
	SELECT ?1, permissions.id FROM permissions WHERE permissions.code IN (?2);`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, userID, codes)
	return err
}
