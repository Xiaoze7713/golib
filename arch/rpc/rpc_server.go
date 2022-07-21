/**
 * @Author: wenliangzhang
 * @Description:
 * @File: rpc_server
 * @Version: 1.0.0
 * @Date: 2022/4/12 3:35 PM
 */
package rpc

// WebServer 基于rpc协议的服务
type RPCServer struct {
}

// NewRPCServer 创建RPCServer
func NewRPCServer() *RPCServer {
	return &RPCServer{}
}
