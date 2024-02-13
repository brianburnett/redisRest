package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"os"
	"strconv"
)

const HOST = "https://172.16.22.21:9443"

func main() {
	fmt.Println("Start...")

	// not secure
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	//getAllDB() // used primarily for init access testing

	fmt.Println("Create DB...")

	// Hard coding 3 for example code - don't rely on the uid actually being 3
	if !getDB(3) {
		createDB()
	} else {
		fmt.Println("DB exists - skip creating...")
	}

	// Do this once...
	/*
	addRole("DB Viewer", "db_viewer") // uid: 2
	addRole("DB Member", "db_member") // uid: 3
	*/

	fmt.Println("Get roles...")
	getRoles()

	fmt.Println("Add users...")
	// Note: role uids are:
	// admin: 1
	// db_viewer: 2
	// db_member: 3
	addUser("john.doe@example.com", "John Doe", 2)
	addUser("mike.smith@example.com", "Mike Smith", 3)
	addUser("cary.johnson@example.com", "Cary Johnson", 1)

}

func basicAuth() string {
	pw := os.Getenv("password")
	credentials := fmt.Sprintf("admin@rl.org:%s", pw)

	b64Cred := base64.StdEncoding.EncodeToString([]byte(credentials))

	basicAuthStr := fmt.Sprintf("Basic %s", b64Cred)
	return basicAuthStr
}

func getAllDB() {
	urlStr := fmt.Sprintf("%s%s", HOST, "/v1/bdbs")
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", basicAuth())
	req.Header.Add("Accept", "application/json")

	resp, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		fmt.Printf("Error: %s\n", err2)
		os.Exit(1)
	}

	defer resp.Body.Close()

	bytes, err3 := io.ReadAll(resp.Body)
	if err3 != nil {
		fmt.Printf("Error: %s\n", err3)
		os.Exit(1)
	}

	fmt.Printf("Resp: %s %s\n", resp.Status, string(bytes))
}

func getDB(uid int) bool {
	urlStr := fmt.Sprintf("%s%s%d", HOST, "/v1/bdbs/", uid)
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}


	req.Header.Add("Authorization", basicAuth())
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		fmt.Printf("Error: %s\n", err2)
		os.Exit(1)
	}


	defer resp.Body.Close()

	bytes, err3 := io.ReadAll(resp.Body)
	if err3 != nil {
		fmt.Printf("Error: %s\n", err3)
		os.Exit(1)
	}

	fmt.Printf("Resp: %s %s\n", resp.Status, string(bytes))

	return is2XX(resp.Status)
}

func createDB() bool {
	type DB struct {
		Name string      `json:"name"`
		MemSizeBytes int `json:"memory_size"`
	}

	dbParams := DB{
		Name: "brian-test",
		MemSizeBytes: 2000000000,
	}

	postBodyJson, _ := json.Marshal(dbParams)

	urlStr := fmt.Sprintf("%s%s", HOST, "/v1/bdbs")
	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewReader(postBodyJson))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", basicAuth())
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Add dry-run query param
	/*
	values := req.URL.Query()
	values.Add("dry_run", "true")
	req.URL.RawQuery = values.Encode()
	*/

	resp, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		fmt.Printf("Error: %s\n", err2)
		os.Exit(1)
	}

	defer resp.Body.Close()

	bytes, err3 := io.ReadAll(resp.Body)
	if err3 != nil {
		fmt.Printf("Error: %s\n", err3)
		os.Exit(1)
	}

	fmt.Printf("Resp: %s %s\n", resp.Status, string(bytes))

	return is2XX(resp.Status)
}

func getRoles() {
	urlStr := fmt.Sprintf("%s%s", HOST, "/v1/roles")
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", basicAuth())
	req.Header.Add("Accept", "application/json")

	resp, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		fmt.Printf("Error: %s\n", err2)
		os.Exit(1)
	}

	defer resp.Body.Close()

	bytes, err3 := io.ReadAll(resp.Body)
	if err3 != nil {
		fmt.Printf("Error: %s\n", err3)
		os.Exit(1)
	}

	fmt.Printf("Resp: %s %s\n", resp.Status, string(bytes))
}

func addRole(name string, management string) bool {
	type Role struct {
		Name string       `json:"name"`
		Management string `json:"management"`
	}

	roleParams := Role{
		Name: name,
		Management: management,
	}

	postBodyJson, _ := json.Marshal(roleParams)

	urlStr := fmt.Sprintf("%s%s", HOST, "/v1/roles")
	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewReader(postBodyJson))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", basicAuth())
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		fmt.Printf("Error: %s\n", err2)
		os.Exit(1)
	}

	defer resp.Body.Close()

	bytes, err3 := io.ReadAll(resp.Body)
	if err3 != nil {
		fmt.Printf("Error: %s\n", err3)
		os.Exit(1)
	}

	fmt.Printf("Resp: %s %s\n", resp.Status, string(bytes))

	return is2XX(resp.Status)
}

func addUser(email string, name string, role int) bool {
	type User struct {
		Email string      `json:"email"`
		Name string       `json:"name"`
		Password string   `json:"password"`
		AuthMethod string `json:"auth_method"`
		RoleUids []int    `json:"role_uids"`
	}

	ruid := make([]int, 1)
	ruid[0] = role

	userParams := User{
		Email: email,
		Name: name,
		Password: "password",
		AuthMethod: "regular",
		RoleUids: ruid,
	}

	postBodyJson, _ := json.Marshal(userParams)

	urlStr := fmt.Sprintf("%s%s", HOST, "/v1/users")
	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewReader(postBodyJson))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", basicAuth())
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		fmt.Printf("Error: %s\n", err2)
		os.Exit(1)
	}

	defer resp.Body.Close()

	bytes, err3 := io.ReadAll(resp.Body)
	if err3 != nil {
		fmt.Printf("Error: %s\n", err3)
		os.Exit(1)
	}

	fmt.Printf("Resp: %s %s\n", resp.Status, string(bytes))

	return is2XX(resp.Status)
}


func is2XX(status string) bool {
	sNum, err := strconv.Atoi(status)
	if err == nil {
		if sNum > 199 && sNum < 300 {
			return true
		}
	}

	return false
}

