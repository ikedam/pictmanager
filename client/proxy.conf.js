const proxy = {
  '/api': {
    target: 'http://server:8080'
  },
  '/storage': {
    target: 'http://storage:4443',
    changeOrigin: true,
    pathRewrite: {
      '^/storage': ''
    }
  }
};

module.exports = proxy;
