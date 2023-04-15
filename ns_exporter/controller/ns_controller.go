package controller

import (
	"fmt"
	"kube_expoter/ex_config"
	"kube_expoter/linkKube"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func PrometheusCotroller(optins *ex_config.Optins) {
	namespace := ""
	prometheus.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        "namespace",
		Namespace:   "kube_delete",
		Help:        "This is a deleted namespace",
		ConstLabels: prometheus.Labels{"Namespace": "Delete_Operation"},
	}, func() float64 {

		// 由于下面代码需要使用到 channel 所以这里我开启一个线程，同时执行
		go func() {
			// 实例化 SharedInformer
			// NewSharedInformerFactory 的参数
			// 1. 与k8s交互的客户端
			// 2. resync 的时间，如果传入0，则禁用resync功能，该功能使用List 操作
			stopCh := make(chan struct{})
			defer close(stopCh)
			sharedInformers := informers.NewSharedInformerFactory(linkKube.Link_kube(optins.KubeConfig.ConfigPath), time.Minute)
			// 生成 pod 资源的informer
			informer := sharedInformers.Core().V1().Namespaces().Informer()
			// 添加事件回调
			informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
				// 创建资源时触发
				AddFunc: nil,
				// 更新资源时触发
				UpdateFunc: nil,
				// 删除资源时触发
				DeleteFunc: func(obj interface{}) {
					DelNamespace := obj.(v1.Object)
					namespace = DelNamespace.GetName()
					//log.Println("delete namespace ", DelNamespace.GetName())
				},
			})
			// 运行 informer
			informer.Run(stopCh)
		}()
		//fmt.Println("PrometheusCotroller() ", namespace)

		logger := lumberjack.Logger{
			Filename: optins.Log.FileName,
			MaxAge:   optins.Log.Max_age,
			Compress: optins.Log.Compress,
			MaxSize:  optins.Log.Max_size,
		}

		defer func(logger *lumberjack.Logger) {
			err := logger.Close()
			if err != nil {
				_ = fmt.Errorf("error: %s", err.Error())
			}
		}(&logger)

		logLevel, _ := logrus.ParseLevel(optins.Log.Level)

		logrus.SetOutput(&logger)
		logrus.SetReportCaller(true)
		logrus.SetLevel(logLevel)

		logrus.WithFields(
			logrus.Fields{
				"DeleteNS": namespace,
			}).Info("Kube_Exporter")

		if namespace != "" {
			return 1.0
		} else {
			return 0.0
		}
	}))
	//fmt.Println(tempNS)
}
