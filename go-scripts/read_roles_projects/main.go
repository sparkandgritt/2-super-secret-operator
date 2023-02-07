package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/iam/v1"
)

func main() {
	ctx := context.Background()
	iamService, err := iam.NewService(ctx)
	if err != nil {
		fmt.Printf("error occured while getting service %v", err)
		return
	}

	roles, err := iamService.Projects.Roles.List("projects/" + "kubebuilder-try-1").Do()
	if err != nil {
		log.Fatalf("failed to list roles %v", err)
	}

	nrole := &iam.CreateRoleRequest{
		Role: &iam.Role{
			Deleted:     false,
			Description: "new role from go lang",
			Stage:       "GA",
			Title:       "Go Lang 2",
		},
		RoleId: "go_lang_2",
	}

	for _, role := range roles.Roles {
		fmt.Println(" - " + role.Name)
	}

	_, err = iamService.Projects.Roles.Create("projects/kubebuilder-try-1", nrole).Context(ctx).Do()
	if err != nil {
		fmt.Printf("error occured while creating role %v", err)
	}
	fmt.Println("=======printing again ==============")
	fmt.Println("=======printing again ==============")
	rolesw, err := iamService.Projects.Roles.List("projects/" + "kubebuilder-try-1").Do()
	if err != nil {
		log.Fatalf("failed to list roles %v", err)
	}
	for _, role := range rolesw.Roles {
		fmt.Println(" - " + role.Name)
	}

}
