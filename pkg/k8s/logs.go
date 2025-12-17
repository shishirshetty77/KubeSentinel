package k8s

import (
	"bytes"
	"context"
	"fmt"
	"io"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// FetchPodLogs retrieves the last N lines of logs from a specific container
func FetchPodLogs(clientset *kubernetes.Clientset, pod *corev1.Pod, containerName string) (string, error) {
	podLogOpts := corev1.PodLogOptions{
		Container: containerName,
		TailLines: func() *int64 { i := int64(50); return &i }(), // Last 50 lines
		Previous:  true,                                          // Important: Get logs from the CRASHED instance, not the current restarting one
	}

	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		// Try without "Previous" if it fails (maybe it's the first crash)
		podLogOpts.Previous = false
		req = clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
		podLogs, err = req.Stream(context.Background())
		if err != nil {
			return "", fmt.Errorf("error in opening stream: %v", err)
		}
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", fmt.Errorf("error in copy information from podLogs to buf: %v", err)
	}

	return buf.String(), nil
}
