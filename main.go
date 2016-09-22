package main

import "fmt"
import "os"
import "errors"
import "bufio"
import "strings"
import "bytes"

func main() {
	path, err := ansibleConfigPath()
	if err != nil {
		fmt.Println("Cannot find Ansible configuration file")
	}

	rolesPath, err := readRolesPath(path)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Roles path is:", rolesPath)

	createPlaybookStructure(rolesPath, "my-playbook")
}

func ansibleConfigPath() (path string, err error) {
	envPath := os.Getenv("ANSIBLE_CONFIG")

	if envPath != "" {
		return envPath, nil
	}

	if _, err := os.Stat("./ansible.cfg"); err == nil {
		return "./ansible.cfg", nil
	}

	if _, err := os.Stat("./.ansible.cfg"); err == nil {
		return "./.ansible.cfg", nil
	}

	if _, err := os.Stat("/etc/ansible/ansible.cfg"); err == nil {
		return "/etc/ansible/ansible.cfg", nil
	}

	return "", errors.New("Ansible config not found")
}

func readRolesPath(path string) (rolesPath string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", errors.New("Cannot open Ansible configuration file")
	}
	defer file.Close()

	parts := strings.SplitAfter(path, "/")
	prefix := strings.Join(parts[:len(parts)-1], "")

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "roles_path") {
			option := scanner.Text()
			rolesPath = strings.TrimSpace(strings.Split(option, "=")[1])
			return concat(prefix, rolesPath), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", errors.New("Cannot read data from Ansible configuration file")
	}

	return concat(prefix, "roles"), nil
}

func concat(prefix string, suffix string) string {
	var buffer bytes.Buffer
	buffer.WriteString(prefix)
	buffer.WriteString(suffix)
	return buffer.String()
}

func createPlaybookStructure(rolesPath string, name string) {
	folders := [...]string{"tasks", "handlers", "templates", "files", "vars", "defaults", "meta"}

	if string(rolesPath[len(rolesPath)-1]) != "/" {
		rolesPath = concat(rolesPath, "/")
	}

	playbookPath := concat(rolesPath, name)
	fmt.Println(playbookPath)
	if string(playbookPath[len(playbookPath)-1]) != "/" {
		playbookPath = concat(playbookPath, "/")
	}
	fmt.Println(playbookPath)

	for _, folder := range folders {
		folderPath := concat(playbookPath, folder)
		fmt.Println(folderPath)
		os.MkdirAll(folderPath, 0755)

		if folder != "files" && folder != "templates" {
			filePath := concat(folderPath, "/main.yml")
			fmt.Println(filePath)
			os.Create(filePath)
		}
	}
}
