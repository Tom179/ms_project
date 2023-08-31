package common

//启动服务器
import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(r *gin.Engine, srvName string, addr string, stop func()) { //传入Engine、服务名字、端口
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("%s服务器,启动！%s\n", srvName, srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("启动失败:", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //监听中断命令

	<-quit //阻塞直到接收值，才继续执行后续代码:关闭服务器
	log.Printf("管道接收到停止信号，关闭%s服务器...\n", srvName)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()   //取消等待，直接结束。这里是没有达到2秒时手动结束操作
	if stop != nil { //停止微服务
		fmt.Println("stop函数不为空,停止微服务")
		stop()
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln("服务器关闭错误：", err)
	}

	select { //select 语句，用于阻塞当前协程，直到其中的 case 中的某个通道操作可以被执行
	case <-ctx.Done(): //ctx结束（超时或取消）
		log.Println("等待超时...")
	}

	log.Printf("%s服务器成功关闭!\n", srvName)

}
