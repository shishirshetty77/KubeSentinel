package notifier

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RecordEvent creates a Kubernetes Event on the target pod with the AI analysis
func RecordEvent(clientset *kubernetes.Clientset, pod *corev1.Pod, analysis string) error {
	t := metav1.Time{Time: time.Now()}

	event := &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: pod.Name + "-ai-analysis-",
			Namespace:    pod.Namespace,
		},
		InvolvedObject: corev1.ObjectReference{
			Kind:            "Pod",
			Namespace:       pod.Namespace,
			Name:            pod.Name,
			UID:             pod.UID,
			APIVersion:      "v1",
			ResourceVersion: pod.ResourceVersion,
		},
		Reason:  "AIAnalysis",
		Message: analysis,
		Type:    corev1.EventTypeWarning,
		Source: corev1.EventSource{
			Component: "KubeSentinel",
		},
		FirstTimestamp: t,
		LastTimestamp:  t,
		Count:          1,
	}

	_, err := clientset.CoreV1().Events(pod.Namespace).Create(context.Background(), event, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("❌ Failed to create K8s event: %v\n", err)
		return err
	}
	fmt.Printf("✅ Created K8s Event for pod %s\n", pod.Name)
	return nil
}
