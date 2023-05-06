const proxy = {
  '/api': {
    target: 'http://server:8080'
  },
  '/image': {
    target: 'http://storage:4443',
    changeOrigin: true,
    pathRewrite: {
      '^/image': ''
    }
  }
};

module.exports = proxy;
