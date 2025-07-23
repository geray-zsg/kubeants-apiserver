package k8s

import (
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"kubeants.io/config"
	"kubeants.io/util"
)

type ExecService struct{}

func NewExecService() *ExecService {
	return &ExecService{}
}

// 升级 WebSocket
func UpgradeToWebSocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	return upgrader.Upgrade(w, r, nil)
}

// WebSocket -> Exec 流转
func (s *ExecService) ExecToPod(wsConn *websocket.Conn, namespace, podName, container, command string) error {
	cmd := strings.Split(command, " ")
	cfg := config.KubeRestConfig
	clientset := config.KubeClientSet

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   cmd,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
	if err != nil {
		return err
	}

	handler := &util.WebSocketStreamHandler{Conn: wsConn}

	return exec.Stream(remotecommand.StreamOptions{
		Stdin:  handler,
		Stdout: handler,
		Stderr: handler,
		Tty:    true,
	})
}
