service_name = "test-servername"
log_addr = "logstash"
log_level = "error"

server = {
  http = {
    listen = "test:port"
  }
  grpc {
    listen = "test:port"
  }
  profile {
    mode = "net"
  }
}