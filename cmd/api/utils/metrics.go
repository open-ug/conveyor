package utils

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func GetCPUUsage(containerName string, start string, end string) (string, error) {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	baseURL := "http://localhost:9090/api/v1/query_range?query=rate(container_cpu_user_seconds_total{name=\"" + containerName + "\"}[1m])&start=" + start + "&end=" + end + "&step=15s"

	resp, err := client.R().
		Get(baseURL)
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}
	return resp.String(), nil
}

func GetMemoryUsage(containerName string, start string, end string) (string, error) {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	baseURL := "http://localhost:9090/api/v1/query_range?query=container_memory_usage_bytes{name=\"" + containerName + "\"}&start=" + start + "&end=" + end + "&step=15s"

	resp, err := client.R().
		Get(baseURL)
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}
	return resp.String(), nil
}
