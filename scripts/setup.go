package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	instancepb "cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	databasepb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

const (
	projectID  = "test-project"
	instanceID = "test-instance"
	databaseID = "test-database"
)

func main() {
	// Set up context
	ctx := context.Background()

	// Setup clients to use emulator
	opts := []option.ClientOption{
		option.WithEndpoint("localhost:9010"),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	}

	// Create instance admin client
	instanceAdminClient, err := instance.NewInstanceAdminClient(ctx, opts...)
	if err != nil {
		log.Fatalf("Failed to create instance admin client: %v", err)
	}
	defer instanceAdminClient.Close()

	// Create instance
	fmt.Println("Creating instance...")
	instancePath := fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID)
	op, err := instanceAdminClient.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", projectID),
		InstanceId: instanceID,
		Instance: &instancepb.Instance{
			DisplayName: "Test Instance",
			Config:      "emulator-config",
			NodeCount:   1,
		},
	})
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			fmt.Println("Instance already exists ✓")
		} else {
			log.Printf("Failed to create instance (may already exist): %v", err)
		}
	} else {
		// Wait for instance creation
		instance, err := op.Wait(ctx)
		if err != nil {
			log.Printf("Failed to wait for instance creation: %v", err)
		} else {
			fmt.Printf("Instance created: %s ✓\n", instance.Name)
		}
	}

	// Create database admin client
	databaseAdminClient, err := database.NewDatabaseAdminClient(ctx, opts...)
	if err != nil {
		log.Fatalf("Failed to create database admin client: %v", err)
	}
	defer databaseAdminClient.Close()

	// Read schema file
	schemaBytes, err := ioutil.ReadFile("test/schema.sql")
	if err != nil {
		log.Fatalf("Failed to read schema file: %v", err)
	}

	// Parse DDL statements
	schemaContent := string(schemaBytes)
	ddlStatements := parseDDL(schemaContent)

	// Create database
	fmt.Println("Creating database...")
	op2, err := databaseAdminClient.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
		Parent:          instancePath,
		CreateStatement: fmt.Sprintf("CREATE DATABASE `%s`", databaseID),
		ExtraStatements: ddlStatements,
	})
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			fmt.Println("Database already exists ✓")
		} else {
			log.Printf("Failed to create database (may already exist): %v", err)
		}
	} else {
		// Wait for database creation
		db, err := op2.Wait(ctx)
		if err != nil {
			log.Printf("Failed to wait for database creation: %v", err)
		} else {
			fmt.Printf("Database created: %s ✓\n", db.Name)
		}
	}

	fmt.Println("Database setup complete!")
}

func parseDDL(content string) []string {
	var statements []string
	
	// Split by lines and process
	lines := strings.Split(content, "\n")
	var currentStatement strings.Builder
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		
		currentStatement.WriteString(line)
		
		// If line ends with semicolon, we have a complete statement
		if strings.HasSuffix(line, ";") {
			stmt := strings.TrimSpace(strings.TrimSuffix(currentStatement.String(), ";"))
			if stmt != "" {
				statements = append(statements, stmt)
			}
			currentStatement.Reset()
		} else {
			currentStatement.WriteString(" ")
		}
	}
	
	// Handle any remaining statement
	if currentStatement.Len() > 0 {
		stmt := strings.TrimSpace(currentStatement.String())
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}
	
	return statements
}