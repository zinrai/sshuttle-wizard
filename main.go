package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var privateNetworks []*net.IPNet

func init() {
	privateRanges := []string{
		"10.0.0.0/8",     // Class A
		"172.16.0.0/12",  // Class B
		"192.168.0.0/16", // Class C
	}

	for _, cidr := range privateRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Sprintf("Invalid CIDR notation: %s", cidr))
		}
		privateNetworks = append(privateNetworks, network)
	}
}

func main() {
	fmt.Println("sshuttle-wizard")
	fmt.Println("---------------")
	fmt.Println("Welcome to the sshuttle command builder and executor wizard!")

	if !isSshuttleInstalled() {
		fmt.Println("Error: sshuttle is not installed or not in your PATH.")
		fmt.Println("Please install sshuttle and try again.")
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	remoteHost := getInput(reader, "Enter remote host (e.g. user@example.com)")

	subnets, err := getRemoteSubnets(remoteHost)
	if err != nil {
		fmt.Printf("Error getting remote subnets: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nDetected private subnets:")
	for i, subnet := range subnets {
		fmt.Printf("%d. %s\n", i+1, subnet)
	}

	chosenSubnets := getChosenSubnets(reader, subnets)

	options := getInput(reader, "Enter additional options (e.g. -v for verbose)")

	command := buildCommand(remoteHost, chosenSubnets, options)

	fmt.Println("\nPrepared sshuttle command:")
	fmt.Println(command)

	if getInput(reader, "Do you want to execute this command? (y/n)") == "y" {
		executeCommand(command)
	} else {
		fmt.Println("Command not executed. You can run it manually if needed.")
	}
}

func getInput(reader *bufio.Reader, prompt string) string {
	fmt.Printf("%s: ", prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func buildCommand(remoteHost string, subnets []string, options string) string {
	command := "sshuttle"

	if options != "" {
		command += " " + options
	}

	for _, subnet := range subnets {
		command += " " + subnet
	}

	command += " -r " + remoteHost

	return command
}

func isSshuttleInstalled() bool {
	_, err := exec.LookPath("sshuttle")
	return err == nil
}

func getRemoteSubnets(remoteHost string) ([]string, error) {
	cmd := exec.Command("ssh", remoteHost, "ip route")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("unable to run command: %v", err)
	}

	return parsePrivateSubnets(string(output)), nil
}

func parsePrivateSubnets(output string) []string {
	var subnets []string
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 {
			subnet := fields[0]
			if isPrivateSubnet(subnet) {
				subnets = append(subnets, subnet)
			}
		}
	}
	return subnets
}

func isPrivateSubnet(subnet string) bool {
	ip, _, err := net.ParseCIDR(subnet)
	if err != nil {
		return false
	}

	for _, network := range privateNetworks {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

func getChosenSubnets(reader *bufio.Reader, subnets []string) []string {
	var chosenSubnets []string
	for {
		input := getInput(reader, "Enter the numbers of the subnets to route (comma-separated, or press Enter to finish)")
		if input == "" {
			break
		}

		choices := strings.Split(input, ",")
		for _, choice := range choices {
			index := parseInt(strings.TrimSpace(choice)) - 1
			if index >= 0 && index < len(subnets) {
				chosenSubnets = append(chosenSubnets, subnets[index])
			} else {
				fmt.Printf("Invalid choice: %s. Skipping.\n", choice)
			}
		}

		if len(chosenSubnets) > 0 {
			break
		} else {
			fmt.Println("No valid subnets selected. Please try again.")
		}
	}
	return chosenSubnets
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func executeCommand(command string) {
	fmt.Println("Executing sshuttle command...")

	args := strings.Fields(command)

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing sshuttle: %v\n", err)
	}
}
