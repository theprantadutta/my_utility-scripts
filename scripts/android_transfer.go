package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const adbDownloadURL = "https://dl.google.com/android/repository/platform-tools-latest-"

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: android_transfer <source_file_path> <destination_folder_on_android>")
		os.Exit(1)
	}

	sourceFilePath := os.Args[1]
	destinationFolderOnAndroid := os.Args[2]

	// Check if ADB is installed
	adbPath, err := exec.LookPath("adb")
	if err != nil {
		fmt.Println("ADB not found in the system. Downloading ADB...")

		// Download ADB
		adbPath, err = downloadADB()
		if err != nil {
			fmt.Printf("Error downloading ADB: %s\n", err)
			os.Exit(1)
		}
	}

	// Execute ADB command to push the file to the Android device
	cmd := exec.Command(adbPath, "push", sourceFilePath, destinationFolderOnAndroid)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error copying file: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("File copied successfully!")
}

func downloadADB() (string, error) {
	fmt.Println("Downloading ADB...")

	platform := runtime.GOOS
	adbURL := adbDownloadURL + platform + ".zip"

	// Create the target directory
	targetDir := "C:\\adb"
	os.MkdirAll(targetDir, os.ModePerm)

	// Remove the existing platform-tools directory if it exists
	existingPlatformToolsPath := filepath.Join(targetDir, "platform-tools")
	if _, err := os.Stat(existingPlatformToolsPath); err == nil {
		// Try to terminate processes using files in the platform-tools directory
		err := terminateProcessesUsingPath(existingPlatformToolsPath)
		if err != nil {
			fmt.Printf("Error terminating processes: %s\n", err)
		}

		// Try to remove the directory, ignore errors
		_ = os.RemoveAll(existingPlatformToolsPath)
	}

	// Create a temporary file to store the downloaded ADB
	tempFile, err := os.CreateTemp("", "adb_download_*.zip")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())

	// Download ADB
	response, err := http.Get(adbURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Write the downloaded content to the temporary file
	_, err = io.Copy(tempFile, response.Body)
	if err != nil {
		return "", err
	}

	// Close the file explicitly
	tempFile.Close()

	// Unzip ADB using the archive/zip package
	zipReader, err := zip.OpenReader(tempFile.Name())
	if err != nil {
		return "", err
	}
	defer zipReader.Close()

	// Extract each file from the zip archive
	for _, file := range zipReader.File {
		filePath := filepath.Join(targetDir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
		} else {
			// Open the file in the zip archive
			fileReader, err := file.Open()
			if err != nil {
				return "", err
			}
			defer fileReader.Close()

			// Create the destination file
			destFile, err := os.Create(filePath)
			if err != nil {
				return "", err
			}
			defer destFile.Close()

			// Copy the file contents
			_, err = io.Copy(destFile, fileReader)
			if err != nil {
				return "", err
			}
		}
	}

	// Return the full path to the ADB executable
	adbExecutablePath := filepath.Join(targetDir, "platform-tools", "adb.exe")
	return adbExecutablePath, nil
}

func terminateProcessesUsingPath(path string) error {
	cmd := exec.Command("taskkill", "/F", "/IM", "adb.exe", "/T")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error terminating processes: %v", err)
	}

	return nil
}
