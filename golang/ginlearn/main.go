package main

import "github.com/gin-gonic/gin"

//func main() {
//	setting.Setup()
//	models.SetUp()
//	logging.Setup()
//
//	//endless
//	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
//	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
//	endless.DefaultMaxHeaderBytes = 1 << 20
//	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
//
//	server := endless.NewServer(endPoint, routers.InitRouter())
//	server.BeforeBegin = func(add string) {
//		log.Printf("Actual pid %d", syscall.Getpid())
//	}
//
//	err := server.ListenAndServe()
//	if err != nil {
//		log.Printf("Server err: %v", err)
//	}
//
//	////shutdown
//	//router := routers.InitRouter()
//	//
//	//s := &http.Server{
//	//	Addr:              fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
//	//	Handler:           router,
//	//	TLSConfig:         nil,
//	//	ReadTimeout:       setting.ServerSetting.ReadTimeout,
//	//	WriteTimeout:      setting.ServerSetting.WriteTimeout,
//	//	MaxHeaderBytes:	   1 << 20,
//	//}
//	//
//	//go func() {
//	//	if err := s.ListenAndServe(); err != nil {
//	//		log.Printf("Listen: %s\n", err)
//	//	}
//	//}()
//	//
//	//quit := make(chan os.Signal)
//	//signal.Notify(quit, os.Interrupt)
//	////堵塞住
//	//<-quit
//	//
//	//log.Println("Shutdown server ...")
//	//ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
//	//defer cancel()
//	//if err := s.Shutdown(ctx); err != nil {
//	//	log.Fatal("Server shutdown:", err)
//	//}
//	//log.Println("Server exiting")
//}

func main() {
	r := gin.Default()
	r.GET("/ping", handler)
	r.GET("pong", handler)
	r.GET("/user/:id", handler)
	r.GET("/user/:id/add", handler)

	r.Run()
}

func handler(c *gin.Context) {
	c.JSON(200, gin.H{"message":"pong"})
}

