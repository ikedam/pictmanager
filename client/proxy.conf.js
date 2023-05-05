const proxy = {
  '/api': {
    'target': 'http://server:8080'
  }
};

module.exports = proxy;
