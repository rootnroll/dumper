package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func dumpContainer(cli *client.Client, container types.Container, baseLogsDir string) {
	demoName := container.Labels["rootnroll.demo.name"]
	createdAt := time.Unix(container.Created, 0).UTC()
	logsFilename := fmt.Sprintf("%s_%s.log", createdAt.Format("2006-01-02T15-04-05"), container.ID[:12])
	demoLogsDir := filepath.Join(baseLogsDir, demoName)
	logsFilepath := filepath.Join(demoLogsDir, logsFilename)
	if _, err := os.Stat(logsFilepath); !os.IsNotExist(err) {
		// Do not dump the already dumped container
		return
	}
	if err := os.MkdirAll(demoLogsDir, os.ModePerm); err != nil {
		fmt.Println("Failed to create a directory for demo logs:", err)
		return
	}

	logsOptions := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true}
	output, err := cli.ContainerLogs(context.Background(), container.ID, logsOptions)
	if err != nil {
		fmt.Println("Failed to fetch container logs:", err)
		return
	}
	// Write logs to a file
	logsFile, err := os.Create(logsFilepath)
	if err != nil {
		fmt.Println("Failed to create a file:", err)
		return
	}
	defer logsFile.Close()
	logsWriter := bufio.NewWriter(logsFile)
	if _, err := io.Copy(logsWriter, output); err != nil {
		fmt.Println("Failed to write logs to the file:", err)
		return
	}
	if err := logsWriter.Flush(); err != nil {
		fmt.Println("Failed to write logs to the file:", err)
		return
	}

	fmt.Printf("Dumped container: %s (%s)\n", container.ID[:12], demoName)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Dump logs of rootnroll demo containers created more than LIFETIME minutes ago.\n\n")
		fmt.Printf("Usage: %s BASE_LOGS_DIR LIFETIME\n", os.Args[0])
		os.Exit(1)
	}
	baseLogsDir := os.Args[1]
	lifetime, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	filters := filters.NewArgs()
	filters.Add("label", "rootnroll.demo.name")
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters,
	})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		createdAt := time.Unix(container.Created, 0)
		if int(time.Now().Sub(createdAt).Minutes()) >= lifetime {
			// Dump a container created more than 20 minutes ago
			dumpContainer(cli, container, baseLogsDir)
		}
	}
}
