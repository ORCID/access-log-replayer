from http.server import SimpleHTTPRequestHandler, HTTPServer

class LoggingHandler(SimpleHTTPRequestHandler):
    def log_message(self, format, *args):
        with open("server.log", "a") as f:
            f.write("%s - - [%s] %s\n" %
                    (self.client_address[0],
                     self.log_date_time_string(),
                     format % args))

if __name__ == "__main__":
    server = HTTPServer(("localhost", 8983), LoggingHandler)
    print("Serving on port 8983...")
    server.serve_forever()