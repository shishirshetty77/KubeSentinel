package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/shishirshetty77/KubeSentinel/pkg/analyzer"
	"github.com/shishirshetty77/KubeSentinel/pkg/notifier"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type Watcher struct {
	clientset *kubernetes.Clientset
}

func NewWatcher(clientset *kubernetes.Clientset) *Watcher {
	return &Watcher{
		clientset: clientset,
	}
}

func (w *Watcher) Start(ctx context.Context) {
	// Create a shared informer factory that watches all namespaces
	factory := informers.NewSharedInformerFactory(w.clientset, 10*time.Minute)
	podInformer := factory.Core().V1().Pods().Informer()

	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// Optional: Handle new pods (maybe check if they fail immediately)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			newPod := newObj.(*corev1.Pod)
			w.checkPodHealth(newPod)
		},
		DeleteFunc: func(obj interface{}) {
			// Loop: Pod deleted
		},
	})

	fmt.Println("👀 Watching for crashes across all namespaces...")
	factory.Start(ctx.Done())
	factory.WaitForCacheSync(ctx.Done())
}

func (w *Watcher) checkPodHealth(pod *corev1.Pod) {
	// 1. Check if Pod is in a "Bad" state
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.State.Waiting != nil {
			reason := containerStatus.State.Waiting.Reason
			if reason == "CrashLoopBackOff" || reason == "ImagePullBackOff" || reason == "OOMKilled" {
				fmt.Printf("🚨 DETECTED FAILURE: Pod %s/%s is in %s. Initiating AI Analysis...\n", pod.Namespace, pod.Name, reason)

				// 1. Fetch Logs
				logs, err := FetchPodLogs(w.clientset, pod, containerStatus.Name)
				if err != nil {
					fmt.Printf("⚠️ Failed to fetch logs: %v\n", err)
				}

				// 2. AI Analysis
				// Note: accessing "pkg/analyzer" requires import
				analysis := "AI Investigation:\n" + analyzer.Analyze(logs, reason)

				// 3. Notify (K8s Event)
				// Note: accessing "pkg/notifier" requires import
				notifier.RecordEvent(w.clientset, pod, analysis)
			}
		}

		if containerStatus.State.Terminated != nil && containerStatus.State.Terminated.ExitCode != 0 {
			// Handle immediate terminations if needed
		}
	}
}
