package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/shishirshetty77/KubeSentinel/pkg/k8s"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	// 1. Setup Configuration
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// 2. Initialize K8s Client
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("Using local kubeconfig: %v\n", err)
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Fatalf("Error loading kubeconfig: %s", err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating clientset: %s", err.Error())
	}

	log.Println("Starting KubeSentinel Controller...")
	log.Println("Successfully connected to Kubernetes Cluster")

	// 3. Start Controller
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Watcher
	watcher := k8s.NewWatcher(clientset)
	go watcher.Start(ctx)

	// 4. Handle Shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down KubeSentinel...")
}
