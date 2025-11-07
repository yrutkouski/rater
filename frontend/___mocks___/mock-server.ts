import * as http from 'http';

interface MockData {
  [key: string]: unknown;
}

const mockData: MockData = {
  '/api/health': { status: 'ok', mode: 'mock' },
  '/api/data': {
    message: 'Mock data from mock server',
    items: [
      { id: 1, name: 'Item 1', value: 100 },
      { id: 2, name: 'Item 2', value: 200 },
      { id: 3, name: 'Item 3', value: 300 },
    ],
  },
};

const server = http.createServer((req, res) => {
  const data = mockData[req.url || ''] || { error: 'Not found' };
  const statusCode = mockData[req.url || ''] ? 200 : 404;

  res.writeHead(statusCode, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify(data));
});

const PORT = 8080;
server.listen(PORT, () => {
  console.log(`Mock API server running on http://localhost:${PORT}`);
});

