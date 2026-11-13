package filelib

import (
	"fmt"
	"os/exec"
	"strings"
)

func PrepareTmpFS(path string, sizeMb int) error {
	var stdout, stderr strings.Builder
	var cmd *exec.Cmd

	if ok, err := FileExists(path, true); !ok || err != nil {
		return fmt.Errorf("directory %s not exists — create it manually with \"sudo\"", path)
	}

	fmt.Printf("PREPARING TEMPORARY FS IN %s FOR %dMB\n", path, sizeMb)

	if ok, _ := FileExistsAny(path); ok {
		cmd = exec.Command("sudo", "umount", path)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(stderr.String())
		}
	}

	// sudo mount -t tmpfs -o size=1024M tmpfs /media/ramdisk

	cmd = exec.Command("sudo", "mount", "-t", "tmpfs", "-o", fmt.Sprintf("size=%dM", sizeMb), "tmpfs", path)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Printf("PREPARED TEMPORARY FS IN %s FOR %dMB\n", path, sizeMb)

	return nil
}
