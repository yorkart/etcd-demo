package etcd3

import (
	"log"

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
	"github.com/coreos/etcd/auth/authpb"
)

type Permission struct {
	Key string
	RangeEnd string
	Type authpb.Permission_Type
}

func InitAuth() {
}

func enableAuth(authAPI clientv3.Auth) {
	perms := []*Permission{
		&Permission{Key:"/", RangeEnd:"", Type: clientv3.PermReadWrite},
	}
	if err := createRoleWithPermission("root", perms , authAPI); err != nil {
		log.Fatal(err)
	}

	if err := listRolePerm("root", authAPI); err != nil {
		log.Fatal(err)
	}

	if err := addUser("root", "root", "P@ssw0rd", authAPI); err != nil {
		log.Fatal(err)
	}

	if _, err := authAPI.AuthEnable(context.TODO()); err != nil {
		log.Fatal(err)
	}
}

func createWebAuth(authAPI clientv3.Auth) {
	perms := []*Permission{
		&Permission{Key:"/warden/", RangeEnd:"/warden/z", Type: clientv3.PermReadWrite},
		&Permission{Key:"/web/", RangeEnd:"/web/z", Type: clientv3.PermReadWrite},
	}
	if err := createRoleWithPermission("web", perms , authAPI); err != nil {
		log.Fatal(err)
	}

	if err := addUser("web", "web", "warden@web", authAPI); err != nil {
		log.Fatal(err)
	}
}

func createRoleWithPermission(role string, perms []*Permission, authAPI clientv3.Auth) error {
	if _, err := authAPI.RoleAdd(context.TODO(), role); err != nil {
		return err
	}

	for _, perm := range perms {
		if _, err := authAPI.RoleGrantPermission(
			context.TODO(),
			role, // role name
			perm.Key, // key
			perm.RangeEnd, // range end
			clientv3.PermissionType(perm.Type),
		); err != nil {
			return err
		}
	}

	return nil
}

func listRolePerm(role string, authAPI clientv3.Auth ) error {
	rres, err := authAPI.RoleGet(context.TODO(), "root")
	if err != nil {
		log.Fatal("role list", err)
	}
	log.Println(rres.Perm)

	return nil
}

func addUser(role, user, pass string, authAPI clientv3.Auth) error{
	if _, err := authAPI.UserAdd(context.TODO(), user, pass); err != nil {
		return err
	}

	if _, err := authAPI.UserGrantRole(context.TODO(), user, role); err != nil {
		return err
	}

	return nil
}
