import http = require("http");

const requestListener = (req: http.IncomingMessage, res: http.ServerResponse) => {
    res.writeHead(200);
    res.end('Hello, World!');
}
  
const server = http.createServer(requestListener);
server.listen(8080);
