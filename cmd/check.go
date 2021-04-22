package cmd

import (
	"context"
	"net"
	"os"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	checkPartition = &cobra.Command{
		Use:   "check",
		Short: "check connectivity to a partition and stop kube-controller-manager if unavailable",
		RunE: func(cmd *cobra.Command, args []string) error {
			return checkNRestart(args)
		},
	}
	dialErrs   int
	deployment string
	namespace  string
	target     string
	maxTries   int
	timeout    time.Duration
	replicas   int32
)

func init() {
	viper.BindPFlags(checkPartition.Flags())
}

func checkNRestart(args []string) error {

	klog.Infoln("Starting partition-watchdog")
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// init
	ticker := time.NewTicker(time.Duration(viper.GetDuration("checkinterval")))
	dialErrs = 0
	replicas = 1
	namespace = os.Getenv("WATCHDOG_NAMESPACE")
	deployment = viper.GetString("deployment")
	target = viper.GetString("target")
	maxTries = viper.GetInt("tries")
	timeout = viper.GetDuration("timeout")

	for ; true; <-ticker.C {
		checkTarget(c)
		err = scaleDeployment(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkTarget(c *kubernetes.Clientset) {

	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", target)
	if err != nil {
		klog.Error(err)
		if dialErrs < maxTries {
			dialErrs++
		}
		klog.Infof("connection to %s failed dialErrs: %d/%d)", target, dialErrs, maxTries)
		return
	}
	if dialErrs > 0 {
		dialErrs--
		klog.Infof("connection to %s successfull (dialErrs: %d/%d)", target, dialErrs, maxTries)
	}
	conn.Close()
}

func scaleDeployment(c *kubernetes.Clientset) error {

	// dialErrs reached max, scaledown deployment
	if dialErrs == maxTries {
		replicas = 0
	}
	// dialErrs ok again, scaleup deployment
	if replicas == 0 && dialErrs == 0 {
		replicas = 1
	}
	d := c.AppsV1().Deployments(namespace)
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, err := d.Get(context.Background(), deployment, v1.GetOptions{})
		if err != nil {
			klog.Errorf("Failed to get latest version of %s deployment in namespace %s: %s", deployment, namespace, err)
			return err
		}
		if *result.Spec.Replicas != replicas {
			result.Spec.Replicas = &replicas
			_, err = d.Update(context.Background(), result, v1.UpdateOptions{})
			klog.Infof("%q in namespace %q scaled to %d", deployment, namespace, replicas)
			return err
		}
		return nil
	})
	return err
}
